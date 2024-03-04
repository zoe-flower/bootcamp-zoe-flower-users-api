package projections

type Option interface {
	apply(*options)
}

// CreateOnNotExist allows projection methods to create a record even when the underlying record does not exist.
// Use this to patch without having to check for existance first.
func CreateOnNotExist() Option {
	return opt(func(o *options) {
		o.createOnNotExist = true
	})
}

type options struct {
	createOnNotExist bool
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

type opt func(*options)

func (f opt) apply(o *options) {
	f(o)
}
