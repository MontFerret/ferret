package runtime

import (
	"context"
	"io"
	"os"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Options struct {
		params  map[string]core.Value
		logging logging.Options
	}

	Option func(*Options)
)

func NewOptions(setters []Option) *Options {
	opts := &Options{
		params: make(map[string]core.Value),
		logging: logging.Options{
			Writer: os.Stdout,
			Level:  logging.ErrorLevel,
		},
	}

	for _, setter := range setters {
		setter(opts)
	}

	return opts
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

func WithLogFields(fields map[string]interface{}) Option {
	return func(options *Options) {
		options.logging.Fields = fields
	}
}

func (opts *Options) WithContext(parent context.Context) context.Context {
	ctx := core.ParamsWith(parent, opts.params)
	ctx = logging.WithContext(ctx, opts.logging)

	return ctx
}
