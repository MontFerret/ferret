package http

import (
	stdhttp "net/http"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/sethgrid/pester"
)

type (
	Option func(opts *Options)

	Options struct {
		Name             string
		Backoff          pester.BackoffStrategy
		MaxRetries       int
		Concurrency      int
		Proxy            string
		UserAgent        string
		Headers          drivers.HTTPHeaders
		Cookies          drivers.HTTPCookies
		AllowedHTTPCodes []int
	}
)

func newOptions(setters []Option) *Options {
	opts := new(Options)
	opts.Name = DriverName
	opts.Backoff = pester.ExponentialBackoff
	opts.Concurrency = 3
	opts.MaxRetries = 5
	opts.AllowedHTTPCodes = []int{stdhttp.StatusOK}

	for _, setter := range setters {
		setter(opts)
	}

	return opts
}

func WithDefaultBackoff() Option {
	return func(opts *Options) {
		opts.Backoff = pester.DefaultBackoff
	}
}

func WithLinearBackoff() Option {
	return func(opts *Options) {
		opts.Backoff = pester.LinearBackoff
	}
}

func WithExponentialBackoff() Option {
	return func(opts *Options) {
		opts.Backoff = pester.ExponentialBackoff
	}
}

func WithMaxRetries(value int) Option {
	return func(opts *Options) {
		opts.MaxRetries = value
	}
}

func WithConcurrency(value int) Option {
	return func(opts *Options) {
		opts.Concurrency = value
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

func WithAllowedHTTPCode(httpCode int) Option {
	return func(opts *Options) {
		if opts.AllowedHTTPCodes == nil {
			opts.AllowedHTTPCodes = make([]int, 0, 1)
		}
		opts.AllowedHTTPCodes = append(opts.AllowedHTTPCodes, httpCode)
	}
}

func WithAllowedHTTPCodes(httpCodes []int) Option {
	return func(opts *Options) {
		if opts.AllowedHTTPCodes == nil {
			opts.AllowedHTTPCodes = httpCodes
		}
		opts.AllowedHTTPCodes = append(opts.AllowedHTTPCodes, httpCodes...)
	}
}
