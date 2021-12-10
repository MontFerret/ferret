package http

import (
	stdhttp "net/http"
	"time"

	"github.com/gobwas/glob"
	"github.com/sethgrid/pester"

	"github.com/MontFerret/ferret/pkg/drivers"
)

var (
	DefaultConcurrency = 1
	DefaultMaxRetries  = 5
	DefaultTimeout     = time.Second * 30
)

type (
	Option func(opts *Options)

	compiledStatusCodeFilter struct {
		URL  glob.Glob
		Code int
	}

	Options struct {
		*drivers.Options
		Backoff         pester.BackoffStrategy
		MaxRetries      int
		Concurrency     int
		HTTPCodesFilter []compiledStatusCodeFilter
		HTTPTransport   *stdhttp.Transport
		Timeout         time.Duration
	}
)

func NewOptions(setters []Option) *Options {
	opts := new(Options)
	opts.Options = new(drivers.Options)
	opts.Name = DriverName
	opts.Backoff = pester.ExponentialBackoff
	opts.Concurrency = DefaultConcurrency
	opts.MaxRetries = DefaultMaxRetries
	opts.Timeout = DefaultTimeout
	opts.HTTPCodesFilter = make([]compiledStatusCodeFilter, 0, 5)

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
		drivers.WithProxy(address)(opts.Options)
	}
}

func WithUserAgent(value string) Option {
	return func(opts *Options) {
		drivers.WithUserAgent(value)(opts.Options)
	}
}

func WithCustomName(name string) Option {
	return func(opts *Options) {
		drivers.WithCustomName(name)(opts.Options)
	}
}

func WithHeader(name string, value []string) Option {
	return func(opts *Options) {
		drivers.WithHeader(name, value)(opts.Options)
	}
}

func WithHeaders(headers *drivers.HTTPHeaders) Option {
	return func(opts *Options) {
		drivers.WithHeaders(headers)(opts.Options)
	}
}

func WithCookie(cookie drivers.HTTPCookie) Option {
	return func(opts *Options) {
		drivers.WithCookie(cookie)(opts.Options)
	}
}

func WithCookies(cookies []drivers.HTTPCookie) Option {
	return func(opts *Options) {
		drivers.WithCookies(cookies)(opts.Options)
	}
}

func WithAllowedHTTPCode(httpCode int) Option {
	return func(opts *Options) {
		opts.HTTPCodesFilter = append(opts.HTTPCodesFilter, compiledStatusCodeFilter{
			Code: httpCode,
		})
	}
}

func WithAllowedHTTPCodes(httpCodes []int) Option {
	return func(opts *Options) {
		for _, code := range httpCodes {
			opts.HTTPCodesFilter = append(opts.HTTPCodesFilter, compiledStatusCodeFilter{
				Code: code,
			})
		}
	}
}

func WithCustomTransport(transport *stdhttp.Transport) Option {
	return func(opts *Options) {
		opts.HTTPTransport = transport
	}
}

func WithTimeout(duration time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = duration
	}
}
