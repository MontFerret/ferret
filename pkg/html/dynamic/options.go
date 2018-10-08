package dynamic

type (
	Options struct {
		proxy string
	}

	Option func(opts *Options)
)

func WithProxy(address string) Option {
	return func(opts *Options) {
		opts.proxy = address
	}
}
