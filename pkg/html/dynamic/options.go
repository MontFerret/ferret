package dynamic

type (
	Options struct {
		proxy     string
		userAgent string
	}

	Option func(opts *Options)
)

func WithProxy(address string) Option {
	return func(opts *Options) {
		opts.proxy = address
	}
}

func WithUserAgent(value string) Option {
	return func(opts *Options) {
		opts.userAgent = value
	}
}
