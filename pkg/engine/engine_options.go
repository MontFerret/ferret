package engine

import (
	"io"
	"os"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

type (
	options struct {
		compiler  []compiler.Option
		noStdlib  bool
		functions runtime.FunctionsBuilder
		params    map[string]runtime.Value
		logging   runtime.LogSettings
	}

	Option func(env *options) error
)

func newOptions(setters []Option) (*options, error) {
	ns := runtime.NewRootNamespace()
	opts := &options{
		functions: ns.Functions(),
		params:    make(map[string]runtime.Value),
		logging: runtime.LogSettings{
			Writer: os.Stdout,
			Level:  runtime.ErrorLevel,
		},
	}

	for _, opt := range setters {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}

	if !opts.noStdlib {
		if err := stdlib.RegisterLib(ns); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func WithFunctions(funcs runtime.Functions) Option {
	return func(opts *options) error {
		if funcs != nil {
			opts.functions.SetFrom(funcs)
		}

		return nil
	}
}

func WithFunction(name string, function runtime.Function) Option {
	return func(opts *options) error {
		if name == "" {
			return runtime.Error(runtime.ErrMissedArgument, "function name")
		}

		if name == "" || function == nil {
			return runtime.Error(runtime.ErrMissedArgument, "function")
		}

		opts.functions.Set(name, function)

		return nil
	}
}

func WithParam(name string, value any) Option {
	return func(opts *options) error {
		opts.params[name] = runtime.Parse(value)

		return nil
	}
}

func WithParams(params map[string]interface{}) Option {
	return func(options *options) error {
		for name, value := range params {
			options.params[name] = runtime.Parse(value)
		}

		return nil
	}
}

func WithLog(writer io.Writer) Option {
	return func(opts *options) error {
		opts.logging.Writer = writer
		return nil
	}
}

func WithLogLevel(lvl runtime.LogLevel) Option {
	return func(opts *options) error {
		opts.logging.Level = lvl
		return nil
	}
}

func WithLogFields(fields map[string]interface{}) Option {
	return func(opts *options) error {
		opts.logging.Fields = fields
		return nil
	}
}
