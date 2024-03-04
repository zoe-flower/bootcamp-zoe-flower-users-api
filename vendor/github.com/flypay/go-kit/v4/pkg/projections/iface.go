package projections

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . ReadWriter
type ReadWriter interface {
	Reader
	Writer
	Patcher
	Deleter
}

//counterfeiter:generate . CompositeReadWriter
type CompositeReadWriter interface {
	CompositeReader
	Writer
	CompositePatcher
	CompositeDeleter
}

var (
	// ErrNotFound is returned when the identifier was not found in the database.
	ErrNotFound = errors.New("Item not found")

	// ErrDataIntegrityMismatch is returned when the data being written was changed by another process.
	ErrDataIntegrityMismatch = errors.New("Data integrity mismatch")

	// ErrInvalidProjectionsDestination indicates that the destination is not the correct type
	ErrInvalidProjectionsDestination = errors.New("projections must be a slice")

	// ErrIdentifierFieldMissing is returned when the passed data to store does not contain an Identifier field
	ErrIdentifierFieldMissing = errors.New("does not have 'Identifier' field")
)

//counterfeiter:generate . Reader
type Reader interface {
	// Read fetches the data stored with the given identifier and writes the data into the provided projection parameter.
	// Therefore you must provide a pointer to a type which should match the structure of the data written with `Writer.Write`.
	// If there is no projected data for the ID you will receive the `projections.ErrNotFound` error value.
	Read(ctx context.Context, id string, projection interface{}) error
}

//counterfeiter:generate . CompositeReader
type CompositeReader interface {
	// ReadComposite fetches the data stored with the given identifier and range key and writes the data into the provided projections parameter.
	// If there is no projected data for the ID you will receive the `projections.ErrNotFound` error value.
	ReadComposite(ctx context.Context, id string, rangeKey *dynamodb.Condition, projections interface{}) error
}

//counterfeiter:generate . Writer
type Writer interface {
	// Write stores the data as a projection.
	// It will replace any existing data with the same `Identifier` value.
	// The `expires` param allows you to define a time after which the projection will be removed.
	// If nil `expires` is provided the data will live forever.
	Write(ctx context.Context, data interface{}, expires *time.Time) error
}

//counterfeiter:generate . Patcher
type Patcher interface {
	// Patch updates a property of an object without affecting others.
	// Supports nested values with dot notion, e.g. "Name.FirstName".
	// By default Patch only updates existing records, otherwise it returns ErrNotFound
	// To Patch a record that doesn't already exist supply the option `CreateOnNotExist()`
	//
	// Example:
	//   // for an existing record { Identifier: "1", FirstName: "Will", LastName: "Smith", City: "LA" }
	//	 err := db.Patch(ctx, "1", "City", "Philadelphia")
	Patch(ctx context.Context, id, path string, value interface{}, opts ...Option) error

	// PatchFields does the same as Patch, however it does this for multiple fields (paths) at once.
	PatchFields(ctx context.Context, id string, keyValues map[string]interface{}, opts ...Option) error
}

//counterfeiter:generate . CompositePatcher
type CompositePatcher interface {
	// PatchComposite does the same as Patch but for an object in a composite table
	PatchComposite(ctx context.Context, id, rangeVal, path string, value interface{}, opts ...Option) error

	// PatchFieldsComposite does the same as PatchComposite, however it does this for multiple fields (paths) at once.
	PatchFieldsComposite(ctx context.Context, id, rangeVal string, keyValues map[string]interface{}, opts ...Option) error
}

//counterfeiter:generate . Deleter

type Deleter interface {
	Delete(ctx context.Context, id string) error
}

//counterfeiter:generate . CompositeDeleter
type CompositeDeleter interface {
	DeleteComposite(ctx context.Context, id, rangeVal string) error
}
