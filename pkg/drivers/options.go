package drivers

type (
	globalOptions struct {
		defaultDriver string
	}

	GlobalOption func(drv Driver, opts *globalOptions)

	Options struct {
		Name      string
		Proxy     string
		UserAgent string
		Headers   *HTTPHeaders
		Cookies   *HTTPCookies
	}

	Option func(opts *Options)
)

func AsDefault() GlobalOption {
	return func(drv Driver, opts *globalOptions) {
		opts.defaultDriver = drv.Name()
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

func WithCustomName(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func WithHeader(name string, value []string) Option {
	return func(opts *Options) {
		if opts.Headers == nil {
			opts.Headers = NewHTTPHeaders()
		}

		opts.Headers.SetArr(name, value)
	}
}

func WithHeaders(headers *HTTPHeaders) Option {
	return func(opts *Options) {
		if opts.Headers == nil {
			opts.Headers = NewHTTPHeaders()
		}

		headers.ForEach(func(value []string, key string) bool {
			opts.Headers.SetArr(key, value)

			return true
		})
	}
}

func WithCookie(cookie HTTPCookie) Option {
	return func(opts *Options) {
		if opts.Cookies == nil {
			opts.Cookies = NewHTTPCookies()
		}

		opts.Cookies.Set(cookie)
	}
}

func WithCookies(cookies []HTTPCookie) Option {
	return func(opts *Options) {
		if opts.Cookies == nil {
			opts.Cookies = NewHTTPCookies()
		}

		for _, c := range cookies {
			opts.Cookies.Set(c)
		}
	}
}
