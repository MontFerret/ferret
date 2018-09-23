package runtime

import "context"

type (
	Options struct {
		proxy     string
		cdp       string
		variables map[string]interface{}
	}

	Option func(*Options)
)

func newOptions() *Options {
	return &Options{
		cdp:       "http://127.0.0.1:9222",
		variables: make(map[string]interface{}),
	}
}

func WithParam(name string, value interface{}) Option {
	return func(options *Options) {
		options.variables[name] = value
	}
}

func WithBrowser(address string) Option {
	return func(options *Options) {
		options.cdp = address
	}
}

func WithProxy(address string) Option {
	return func(options *Options) {
		// TODO: add implementation
		options.proxy = address
	}
}

func (opts *Options) withContext(parent context.Context) context.Context {
	return context.WithValue(
		parent,
		"variables",
		opts.variables,
	)
}
