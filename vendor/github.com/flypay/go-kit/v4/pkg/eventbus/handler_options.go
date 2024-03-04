package eventbus

import (
	"time"

	"github.com/flypay/go-kit/v4/pkg/log"
)

type HandlerOption interface {
	apply(*HandlerOptions)
}

// WithTimeout can be passed as a variadic parameter to the Consumer.On method and it will override the default event
// handling timeout for that event
func WithTimeout(to time.Duration) HandlerOption {
	return newFuncHandlerOption(func(options *HandlerOptions) {
		if to <= 0 {
			log.Warnf("Invalid handler timeout: %s, default will be used", to)
			return
		}

		options.Timeout = to
	})
}

// MaxAttempts overrides the default setting of HandlerOptions.MaxAttempts
func MaxAttempts(n int) HandlerOption {
	return newFuncHandlerOption(func(opts *HandlerOptions) {
		if n <= 0 {
			log.Warnf("Invalid MaxAttempts %d, default will be used", n)
			opts.MaxAttempts = 1
			return
		}

		opts.MaxAttempts = n
	})
}

// HandlerOptions are settings that can be applied to consumer handlers
type HandlerOptions struct {
	// Timeout is the amount of time after which the consumer handler function will be canceled.
	Timeout time.Duration

	// MaxAttempts is the number of attempts a consumer will try and handle an event that results in a non-nil error
	// before it becomes a dead event (sent to deadletter queue).
	// Setting this to a value > 1 will enable retrying when the handler returns an error.
	// default: 1  (no retries)
	MaxAttempts int
}

func (o *HandlerOptions) Apply(options ...HandlerOption) {
	for _, option := range options {
		option.apply(o)
	}
}

func NewDefaultHandlerOptions(timeout time.Duration) HandlerOptions {
	return HandlerOptions{
		Timeout: timeout,
	}
}

type funcHandlerOption struct {
	f func(options *HandlerOptions)
}

func (fho *funcHandlerOption) apply(options *HandlerOptions) {
	fho.f(options)
}

func newFuncHandlerOption(f func(options *HandlerOptions)) *funcHandlerOption {
	return &funcHandlerOption{
		f: f,
	}
}
