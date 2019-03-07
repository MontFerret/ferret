package cdp

type (
	Options struct {
		proxy     string
		userAgent string
		address   string
	}

	Option func(opts *Options)
)

const DefaultAddress = "http://127.0.0.1:9222"

func newOptions(setters []Option) *Options {
	opts := new(Options)
	opts.address = DefaultAddress

	for _, setter := range setters {
		setter(opts)
	}

	return opts
}

func WithAddress(address string) Option {
	return func(opts *Options) {
		if address != "" {
			opts.address = address
		}
	}
}

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
