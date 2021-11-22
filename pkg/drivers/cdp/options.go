package cdp

import (
	"github.com/MontFerret/ferret/pkg/drivers"
)

type (
	Options struct {
		*drivers.Options
		Address     string
		KeepCookies bool
		Connection  *ConnectionOptions
	}

	ConnectionOptions struct {
		BufferSize  int
		Compression bool
	}

	Option func(opts *Options)
)

const (
	DefaultAddress    = "http://127.0.0.1:9222"
	DefaultBufferSize = 1048562
)

func NewOptions(setters []Option) *Options {
	opts := new(Options)
	opts.Options = new(drivers.Options)
	opts.Name = DriverName
	opts.Address = DefaultAddress
	opts.Connection = &ConnectionOptions{
		BufferSize:  DefaultBufferSize,
		Compression: true,
	}

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
		drivers.WithProxy(address)(opts.Options)
	}
}

func WithUserAgent(value string) Option {
	return func(opts *Options) {
		drivers.WithUserAgent(value)(opts.Options)
	}
}

func WithKeepCookies() Option {
	return func(opts *Options) {
		opts.KeepCookies = true
	}
}

func WithCustomName(name string) Option {
	return func(opts *Options) {
		drivers.WithCustomName(name)(opts.Options)
	}
}

func WithHeader(name string, header []string) Option {
	return func(opts *Options) {
		drivers.WithHeader(name, header)(opts.Options)
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

func WithBufferSize(size int) Option {
	return func(opts *Options) {
		opts.Connection.BufferSize = size
	}
}

func WithCompression() Option {
	return func(opts *Options) {
		opts.Connection.Compression = true
	}
}

func WithNoCompression() Option {
	return func(opts *Options) {
		opts.Connection.Compression = false
	}
}
