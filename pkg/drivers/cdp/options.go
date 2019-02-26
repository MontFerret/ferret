package cdp

type (
	Options struct {
		name      string
		proxy     string
		userAgent string
		address   string
		cookies   bool
	}

	Option func(opts *Options)
)

const DefaultAddress = "http://127.0.0.1:9222"

func newOptions(setters []Option) *Options {
	opts := new(Options)
	opts.name = DriverName
	opts.address = DefaultAddress

	for _, setter := range setters {
		setter(opts)
	}

	return opts
}

func WithAddress(address string) Option {
	return func(opts *Options) {
		opts.address = address
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

func WithCookies() Option {
	return func(opts *Options) {
		opts.cookies = true
	}
}

func WithCustomName(name string) Option {
	return func(opts *Options) {
		opts.name = name
	}
}
