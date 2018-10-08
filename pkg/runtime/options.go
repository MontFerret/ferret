package runtime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/env"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"io"
	"os"
)

type (
	Options struct {
		proxy     string
		cdp       string
		params    map[string]core.Value
		logging   *logging.Options
		userAgent string
	}

	Option func(*Options)
)

func newOptions() *Options {
	return &Options{
		cdp:    "http://0.0.0.0:9222",
		params: make(map[string]core.Value),
		logging: &logging.Options{
			Writer: os.Stdout,
			Level:  logging.ErrorLevel,
		},
	}
}

func WithParam(name string, value interface{}) Option {
	return func(options *Options) {
		options.params[name] = values.Parse(value)
	}
}

func WithParams(params map[string]interface{}) Option {
	return func(options *Options) {
		for name, value := range params {
			options.params[name] = values.Parse(value)
		}
	}
}

func WithBrowser(address string) Option {
	return func(options *Options) {
		options.cdp = address
	}
}

func WithProxy(address string) Option {
	return func(options *Options) {
		options.proxy = address
	}
}

func WithUserAgent(value string) Option {
	return func(options *Options) {
		options.userAgent = value
	}
}

func WithRandomUserAgent() Option {
	return func(options *Options) {
		options.userAgent = env.RandomUserAgent
	}
}

func WithLog(writer io.Writer) Option {
	return func(options *Options) {
		options.logging.Writer = writer
	}
}

func WithLogLevel(lvl logging.Level) Option {
	return func(options *Options) {
		options.logging.Level = lvl
	}
}

func (opts *Options) withContext(parent context.Context) context.Context {
	ctx := core.ParamsWith(parent, opts.params)
	ctx = logging.WithContext(ctx, opts.logging)
	ctx = env.WithContext(ctx, env.Environment{
		CDPAddress:   opts.cdp,
		ProxyAddress: opts.proxy,
		UserAgent:    opts.userAgent,
	})

	return ctx
}
