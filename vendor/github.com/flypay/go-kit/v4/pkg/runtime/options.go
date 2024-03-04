package runtime

import (
	"fmt"
	"path/filepath"

	"github.com/flypay/go-kit/v4/internal/service"
	"github.com/spf13/afero"
)

type Option interface {
	apply(*options)
}

// WithFromENV will set the runtime strategy to use values
// from the ENV when generating the runtime configuration
func WithFromENV() Option {
	return funcOption(func(o *options) {
		o.FromENV = true
	})
}

// WithFromFs will set the runtime strategy to use values from the service.json
// in the specified filesystem when generating the runtime configuration.
func WithFromFs(root afero.Fs) Option {
	return funcOption(func(o *options) {
		o.FromServiceJSON = true
		o.Fs = root
		o.ServiceJSONName = service.JSONFileName
	})
}

// WithFromServiceJSON will set the runtime strategy to use values
// from the service.json when generating the runtime configuration.
func WithFromServiceJSON(path string) Option {
	return funcOption(func(o *options) {
		fmt.Println("NOTE: WithFromServiceJSON is deprecated, please use WithFromFs.")

		o.FromServiceJSON = true
		o.ServiceJSONName = filepath.Base(path)
		o.Fs = afero.NewBasePathFs(afero.NewOsFs(), filepath.Dir(path))
	})
}

func WithShouldRunProducer(enabled bool) Option {
	return funcOption(func(o *options) { o.ShouldRunProducer = enabled })
}

func WithShouldRunObjstore(enabled bool) Option {
	return funcOption(func(o *options) { o.ShouldRunObjstore = enabled })
}

func WithShouldRunHTTPClient(enabled bool) Option {
	return funcOption(func(o *options) { o.ShouldRunHTTPClient = enabled })
}

type options struct {
	FromENV bool

	FromServiceJSON bool
	Fs              afero.Fs

	// DEPRECATED
	ServiceJSONName string

	ShouldRunProducer   bool
	ShouldRunObjstore   bool
	ShouldRunHTTPClient bool
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

type funcOption func(options *options)

func (f funcOption) apply(options *options) {
	f(options)
}
