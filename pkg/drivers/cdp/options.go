package cdp

type (
	Options struct {
		Name        string
		Proxy       string
		UserAgent   string
		Address     string
		KeepCookies bool
	}

	Option func(opts *Options)
)

const DefaultAddress = "http://127.0.0.1:9222"

func newOptions(setters []Option) *Options {
	opts := new(Options)
	opts.Name = DriverName
	opts.Address = DefaultAddress

	for _, setter := range setters {
		setter(opts)
	}

	return opts
}

func WithAddress(address string) Option {
	return func(opts *Options) {
		if address != "" {
			opts.Address = address
		}
	}
}

func WithProxy(address string) Option {
	return func(opts *Options) {
		opts.Proxy = address
	}
}

func WithUserAgent(value string) Option {
	return func(opts *Options) {
		opts.UserAgent = value
	}
}

func WithKeepCookies() Option {
	return func(opts *Options) {
		opts.KeepCookies = true
	}
}

func WithCustomName(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}
