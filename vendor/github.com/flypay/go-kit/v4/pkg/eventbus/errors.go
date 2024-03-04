package eventbus

import (
	"github.com/pkg/errors"
)

var (
	// ErrTooManyHandlersForEvent is returned when attempting to register an event handler
	// on a consumer that has more than 'maxHandlersPerEvent' handlers already
	// registered
	ErrTooManyHandlersForEvent = errors.New("Too many handlers registered for event")

	// ErrTopicNameNotDefined indicates that an event proto does not return a value for the topic_name option
	ErrTopicNameNotDefined = errors.New("topic name not defined for event")
)

// ErrorCode contains the error string
type ErrorCode string

const (
	StatusOK            ErrorCode = "OK"
	StatusInternalError ErrorCode = "InternalError"
	StatusExternalError ErrorCode = "ExternalError"
)

func (ec ErrorCode) String() string {
	return string(ec)
}

// HandlerError is returned from service handlers
// so we can check the error code and Inc() the metrics accordingly
type HandlerError struct {
	ErrorMessage error
	ErrorCode    ErrorCode
	Terminal     bool
}

func (he *HandlerError) Error() string {
	return he.ErrorMessage.Error()
}

func (he *HandlerError) Unwrap() error {
	return he.ErrorMessage
}

func Wrap(err error, code ErrorCode, message string) *HandlerError {
	if err == nil {
		return nil
	}
	return &HandlerError{
		ErrorMessage: errors.Wrap(err, message),
		ErrorCode:    code,
	}
}

func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *HandlerError {
	if err == nil {
		return nil
	}
	return &HandlerError{
		ErrorMessage: errors.Wrapf(err, format, args...),
		ErrorCode:    code,
	}
}

func TerminalError(err error, code ErrorCode, format string, args ...interface{}) *HandlerError {
	if err == nil {
		return nil
	}
	herr := Wrapf(err, code, format, args...)
	herr.Terminal = true
	return herr
}

func NewHandlerError(err error, code ErrorCode) *HandlerError {
	if err == nil {
		return nil
	}
	return &HandlerError{
		ErrorMessage: err,
		ErrorCode:    code,
	}
}
