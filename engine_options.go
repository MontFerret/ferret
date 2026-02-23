package ferret

import (
	"io"
	"os"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
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
		registry  *encoding.Registry
	}

	Option func(env *options) error
)

func newOptions(setters []Option) (*options, error) {
	ns := runtime.NewRootNamespace()
	opts := &options{
		functions: ns.Functions(),
		params:    make(map[string]runtime.Value),
		registry:  encoding.NewRegistry(),
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

// WithNamespace creates an Option that sets the functions from the provided runtime.Namespace to the options if not nil.
func WithNamespace(ns runtime.Namespace) Option {
	return func(opts *options) error {
		if ns == nil {
			return nil
		}

		if opts.functions == nil {
			opts.functions = ns.Functions()
		} else {
			opts.functions.SetFrom(ns.Functions().Build())
		}

		return nil
	}
}

// WithFunctions creates an Option that sets the provided runtime.Functions to the options if not nil.
func WithFunctions(funcs runtime.Functions) Option {
	return func(opts *options) error {
		if funcs != nil {
			opts.functions.SetFrom(funcs)
		}

		return nil
	}
}

// WithFunction returns an Option that registers a runtime.Function with a given name in the options' function builder.
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

// WithLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithLog(writer io.Writer) Option {
	return func(opts *options) error {
		opts.logging.Writer = writer
		return nil
	}
}

// WithLogLevel sets the logging level for the engine.
// The logging level determines the severity of log messages that will be recorded.
func WithLogLevel(lvl runtime.LogLevel) Option {
	return func(opts *options) error {
		opts.logging.Level = lvl
		return nil
	}
}

// WithLogFields sets the fields to be included in log entries.
// These fields can provide additional context for debugging and monitoring purposes.
func WithLogFields(fields map[string]any) Option {
	return func(opts *options) error {
		opts.logging.Fields = fields
		return nil
	}
}

// WithEncodingRegistry sets a custom encoding registry for query execution.
func WithEncodingRegistry(registry *encoding.Registry) Option {
	return func(opts *options) error {
		if registry == nil {
			return runtime.Error(runtime.ErrMissedArgument, "registry")
		}

		opts.registry = registry
		return nil
	}
}

// WithEncodingCodec registers or overrides a codec for the given content type.
func WithEncodingCodec(contentType string, codec encoding.Codec) Option {
	return func(opts *options) error {
		if opts.registry == nil {
			opts.registry = encoding.NewRegistry()
		}

		return opts.registry.Register(contentType, codec)
	}
}
