package cdp

import "github.com/MontFerret/ferret/pkg/drivers"

type (
	Options struct {
		Name        string
		Proxy       string
		UserAgent   string
		Address     string
		KeepCookies bool
		Headers     drivers.HTTPHeaders
		Cookies     drivers.HTTPCookies
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

func WithHeader(name string, value []string) Option {
	return func(opts *Options) {
		if opts.Headers == nil {
			opts.Headers = make(drivers.HTTPHeaders)
		}

		opts.Headers[name] = value
	}
}

func WithHeaders(headers drivers.HTTPHeaders) Option {
	return func(opts *Options) {
		if opts.Headers == nil {
			opts.Headers = make(drivers.HTTPHeaders)
		}

		for k, v := range headers {
			opts.Headers[k] = v
		}
	}
}

func WithCookie(cookie drivers.HTTPCookie) Option {
	return func(opts *Options) {
		if opts.Cookies == nil {
			opts.Cookies = make(drivers.HTTPCookies)
		}

		opts.Cookies[cookie.Name] = cookie
	}
}

func WithCookies(cookies []drivers.HTTPCookie) Option {
	return func(opts *Options) {
		if opts.Cookies == nil {
			opts.Cookies = make(drivers.HTTPCookies)
		}

		for _, c := range cookies {
			opts.Cookies[c.Name] = c
		}
	}
}
