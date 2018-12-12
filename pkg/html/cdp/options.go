package cdp

type (
	Options struct {
		proxy     string
		userAgent string
		cdp       string
	}

	Option func(opts *Options)
)

func newOptions(setters []Option) *Options {
	opts := new(Options)
	opts.cdp = "http://127.0.0.1:9222"

	for _, setter := range setters {
		setter(opts)
	}

	return opts
}

func WithCDP(address string) Option {
	return func(opts *Options) {
		opts.cdp = address
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
