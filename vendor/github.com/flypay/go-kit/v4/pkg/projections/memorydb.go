package projections

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var errInvalidDocumentPath = errors.New("ValidationException: The document path provided in the update expression is invalid for update")

// MemDB is an implementation of the projection reader and writer interfaces
// backed my an in memory map. It is safe to use concurrently.
type MemDB struct {
	storage map[string]*dynamodb.AttributeValue
	mu      sync.RWMutex
}

func (db *MemDB) Write(ctx context.Context, data interface{}, expires *time.Time) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tostore, err := dynamodbattribute.Marshal(data)
	if err != nil {
		return err
	}

	id := identifier(tostore)
	if id == "" {
		return ErrIdentifierFieldMissing
	}

	isNewer := sequence(data)

	if db.storage == nil {
		db.storage = make(map[string]*dynamodb.AttributeValue)
	}

	if expires != nil && expires.Before(time.Now()) {
		if _, ok := db.storage[id]; !ok {
			// Not in db and want to delete
			return nil
		}

		existing := db.storage[id]
		if isNewer(existing) {
			return ErrDataIntegrityMismatch
		}
		delete(db.storage, id)
		return nil
	}

	if expires != nil {
		// This timer will delete the key if it is written again with a new expiry
		// we will need to store a map of the timers and reset them.
		// See https://github.com/flypay/go-kit/issues/344
		time.AfterFunc(time.Until(*expires), func() {
			db.mu.Lock()
			defer db.mu.Unlock()

			if isNewer(db.storage[id]) {
				return
			}

			delete(db.storage, id)
		})
	}

	if _, ok := db.storage[id]; !ok {
		db.storage[id] = tostore
		return nil
	}

	existing := db.storage[id]
	if isNewer(existing) {
		return ErrDataIntegrityMismatch
	}

	db.storage[id] = tostore

	return nil
}

func (db *MemDB) Read(ctx context.Context, id string, dest interface{}) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	existing, found := db.storage[id]
	if !found {
		return ErrNotFound
	}

	return dynamodbattribute.Unmarshal(existing, &dest)
}

func (db *MemDB) Patch(ctx context.Context, id, path string, value interface{}, opts ...Option) error {
	return db.PatchFields(ctx, id, map[string]interface{}{path: value}, opts...)
}

func (db *MemDB) PatchFields(ctx context.Context, id string, keyValues map[string]interface{}, opts ...Option) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.storage == nil {
		db.storage = make(map[string]*dynamodb.AttributeValue)
	}

	var cfg options
	cfg.apply(opts...)

	existing := db.storage[id]
	var isNew bool
	if existing == nil {
		if !cfg.createOnNotExist {
			return ErrNotFound
		}

		existing = &dynamodb.AttributeValue{M: make(map[string]*dynamodb.AttributeValue)}
		idv, _ := dynamodbattribute.Marshal(id)
		existing.M[identifierField] = idv
		isNew = true
	}

	for key, value := range keyValues {
		fieldNames := strings.Split(key, ".")
		if len(fieldNames) > 1 && isNew {
			return errInvalidDocumentPath
		}

		next := existing

		var i int
		for i = 0; i < len(fieldNames)-1; i++ {
			if _, ok := next.M[fieldNames[i]]; !ok {
				newValue := &dynamodb.AttributeValue{M: make(map[string]*dynamodb.AttributeValue)}
				next.M[fieldNames[i]] = newValue
			}

			next = next.M[fieldNames[i]]
		}

		if next.M == nil {
			next.M = make(map[string]*dynamodb.AttributeValue)
		}

		newValue, err := dynamodbattribute.Marshal(value)
		if err != nil {
			return err
		}

		next.M[fieldNames[i]] = newValue
	}

	db.storage[id] = existing

	return nil
}

func (db *MemDB) Delete(ctx context.Context, id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.storage == nil {
		return nil
	}

	delete(db.storage, id)

	return nil
}
