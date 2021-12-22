package runtime

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Options struct {
		mxParam sync.RWMutex
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
		options.mxParam.Lock()
		defer options.mxParam.Unlock()

		options.params[name] = values.Parse(value)
	}
}

func WithParams(params map[string]interface{}) Option {
	return func(options *Options) {
		for name, value := range params {
			options.mxParam.Lock()
			options.params[name] = values.Parse(value)
			options.mxParam.Unlock()
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
