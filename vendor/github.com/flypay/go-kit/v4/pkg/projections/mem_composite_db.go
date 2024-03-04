package projections

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var _ CompositeReadWriter = (*MemCompositeDB)(nil)

const attrLenBetween = 2

// MemCompositeDB is an implementation of the projection composite reader and writer interfaces backed my an in memory slice.
// It is safe to use concurrently and initilizes itself on first use.
type MemCompositeDB struct {
	storage []*dynamodb.AttributeValue
	mu      sync.Mutex
}

func validOutput(projections interface{}) error {
	pOut := reflect.ValueOf(projections)
	if reflect.Ptr != pOut.Kind() {
		return fmt.Errorf("MemCompositeDB.ReadComposite called with non-pointer projections")
	}

	out := pOut.Elem()
	if out.Kind() != reflect.Slice {
		return ErrInvalidProjectionsDestination
	}

	return nil
}

// ReadComposite reads composite data from memory.
func (db *MemCompositeDB) ReadComposite(ctx context.Context, id string, rangeKey *dynamodb.Condition, projections interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := validOutput(projections); err != nil {
		return err
	}

	match, err := buildRangeMatcher(rangeKey)
	if err != nil {
		return err
	}

	var result []*dynamodb.AttributeValue
	for _, record := range db.storage {
		if identifier(record) == id && match(rangeValue(record)) {
			result = append(result, record)
		}
	}

	if len(result) == 0 {
		return ErrNotFound
	}

	return dynamodbattribute.UnmarshalList(result, projections)
}

func buildRangeMatcher(rangeCond *dynamodb.Condition) (func(rangeVal string) bool, error) {
	match := func(rangeVal string) bool { return true }

	if rangeCond == nil {
		return match, nil
	}

	if rangeCond.ComparisonOperator == nil {
		return match, fmt.Errorf("no ComparisonOperator given")
	}

	attrValues := rangeCond.AttributeValueList

	if len(attrValues) < 1 {
		return match, fmt.Errorf("empty AttributeValueList given")
	}

	if attrValues[0].S == nil {
		return match, fmt.Errorf("no string attribute given, no other method implemented")
	}

	rangeQueryValue := *attrValues[0].S

	switch comOp := *rangeCond.ComparisonOperator; comOp {
	case dynamodb.ComparisonOperatorEq:
		match = func(rangeVal string) bool { return rangeVal == rangeQueryValue }
	case dynamodb.ComparisonOperatorLe:
		match = func(rangeVal string) bool { return rangeVal <= rangeQueryValue }
	case dynamodb.ComparisonOperatorLt:
		match = func(rangeVal string) bool { return rangeVal < rangeQueryValue }
	case dynamodb.ComparisonOperatorGe:
		match = func(rangeVal string) bool { return rangeVal >= rangeQueryValue }
	case dynamodb.ComparisonOperatorGt:
		match = func(rangeVal string) bool { return rangeVal > rangeQueryValue }
	case dynamodb.ComparisonOperatorBeginsWith:
		match = func(rangeVal string) bool { return strings.HasPrefix(rangeVal, rangeQueryValue) }
	case dynamodb.ComparisonOperatorBetween:
		if len(attrValues) != attrLenBetween {
			return match, fmt.Errorf("AttributeValueList must have two values for BETWEEN")
		}

		if attrValues[0].S == nil || attrValues[1].S == nil {
			return match, fmt.Errorf("no string attribute given, no other method implemented")
		}

		lower := *attrValues[0].S
		upper := *attrValues[1].S

		match = func(rangeVal string) bool { return rangeVal >= lower && rangeVal <= upper }
	default:
		return match, fmt.Errorf("unknown comparision operator %q", comOp)
	}

	return match, nil
}

// Write writes composite data to memory. Expiry in the future is not implemented.
// To delete data write with an expiry before time.Now()
// TODO with an expiry in the future the data should be removed at that time, consider launching a goroutine to remove the data.
func (db *MemCompositeDB) Write(ctx context.Context, data interface{}, expires *time.Time) error {
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

	newRange := rangeValue(tostore)
	idx := db.index(id, newRange)
	isNewer := sequence(data)

	if expires != nil && expires.Before(time.Now()) {
		if idx < 0 {
			// Not in db and want to delete
			return nil
		}

		existing := db.storage[idx]
		if isNewer(existing) {
			return ErrDataIntegrityMismatch
		}
		db.delete(idx)
		return nil
	}

	if idx < 0 {
		db.storage = append(db.storage, tostore)
		return nil
	}

	existing := db.storage[idx]
	if isNewer(existing) {
		return ErrDataIntegrityMismatch
	}
	db.storage[idx] = tostore
	return nil
}

func (db *MemCompositeDB) delete(i int) {
	db.storage[i] = db.storage[len(db.storage)-1]
	db.storage = db.storage[:len(db.storage)-1]
}

// index returns the position of the data in the storage or -1 if not found
func (db *MemCompositeDB) index(id, newRange string) int {
	for i, existing := range db.storage {
		if identifier(existing) == id && newRange == rangeValue(existing) {
			return i
		}
	}
	return -1
}

// ErrNotImplemented is returned when a feature is not implemented yet
var ErrNotImplemented = errors.New("not implemented yet")

// PatchComposite patches in memory composite data.
func (db *MemCompositeDB) PatchComposite(ctx context.Context, id, rangeVal, path string, value interface{}, opts ...Option) error {
	return db.PatchFieldsComposite(ctx, id, rangeVal, map[string]interface{}{path: value}, opts...)
}

func (db *MemCompositeDB) PatchFieldsComposite(ctx context.Context, id, rangeVal string, keyValues map[string]interface{}, opts ...Option) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	var cfg options
	cfg.apply(opts...)

	var existing *dynamodb.AttributeValue
	var isNew bool
	idx := db.index(id, rangeVal)
	if idx < 0 {
		if !cfg.createOnNotExist {
			return ErrNotFound
		}

		existing = &dynamodb.AttributeValue{M: make(map[string]*dynamodb.AttributeValue)}

		idv, _ := dynamodbattribute.Marshal(id)
		existing.M[identifierField] = idv

		rangev, _ := dynamodbattribute.Marshal(rangeVal)
		existing.M[rangeField] = rangev
		isNew = true
	} else {
		existing = db.storage[idx]
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

	if idx < 0 {
		db.storage = append(db.storage, existing)
		return nil
	}

	db.storage[idx] = existing

	return nil
}

func (db *MemCompositeDB) DeleteComposite(ctx context.Context, id, rangeVal string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.storage == nil {
		return nil
	}

	for i, record := range db.storage {
		if identifier(record) == id && rangeValue(record) == rangeVal {
			db.delete(i)
			return nil
		}
	}

	return nil
}

func rangeValue(v *dynamodb.AttributeValue) string {
	// return v.FieldByName(rangeField).Interface().(string)
	return *v.M[rangeField].S
}

func identifier(v *dynamodb.AttributeValue) string {
	if v.M[identifierField] == nil {
		return ""
	}

	if v.M[identifierField].S == nil {
		return ""
	}

	return *v.M[identifierField].S
}

type hasChanged func(*dynamodb.AttributeValue) bool

// sequence returns a function that returns true if the attribute value is newer than the data trying to be written
func sequence(data interface{}) hasChanged {
	field, v := sequenceField(data)
	if field == nil {
		// Allow data to be written since there was no sequence field
		return func(v *dynamodb.AttributeValue) bool {
			return false
		}
	}

	// Only write data if previous is less than new value
	return func(existing *dynamodb.AttributeValue) bool {
		prev, err := strconv.Atoi(*existing.M[field.Name].N)
		if err != nil {
			return true
		}
		return int64(prev) > v.FieldByName(field.Name).Int()
	}
}
