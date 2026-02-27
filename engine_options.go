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
		compiler []compiler.Option
		noStdlib bool
		lib      runtime.RootNamespace
		params   map[string]runtime.Value
		logging  runtime.LogSettings
		encodig  *encoding.Registry
	}

	Option func(env *options) error
)

func newOptions(setters []Option) (*options, error) {
	ns := runtime.NewRootNamespace()
	opts := &options{
		lib:     ns,
		params:  make(map[string]runtime.Value),
		encodig: encoding.NewRegistry(),
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
		stdlib.RegisterLib(ns)
	}

	return opts, nil
}

// WithNamespace creates an Option that sets the lib from the provided runtime.Namespace to the options if not nil.
func WithNamespace(ns runtime.Namespace) Option {
	return func(opts *options) error {
		if ns == nil {
			return nil
		}

		opts.lib.Function().From(ns.Function())

		return nil
	}
}

// WithFunctionsBuilder creates an Option that sets the lib from the provided runtime.FunctionDefs to the options if not nil.
func WithFunctionsRegistrar(setter func(fns runtime.FunctionDefs)) Option {
	return func(env *options) error {
		if setter != nil {
			setter(env.lib.Function())
		}

		return nil
	}
}

// WithFunctions creates an Option that sets the provided *runtime.Functions to the options if not nil.
func WithFunctions(funcs *runtime.Functions) Option {
	return func(opts *options) error {
		if funcs != nil {
			opts.lib.Function().From(runtime.NewFunctionsBuilderFrom(funcs))
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

		opts.lib.Function().Var().Add(name, function)

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

		opts.encodig = registry
		return nil
	}
}

// WithEncodingCodec registers or overrides a codec for the given content type.
func WithEncodingCodec(contentType string, codec encoding.Codec) Option {
	return func(opts *options) error {
		if opts.encodig == nil {
			opts.encodig = encoding.NewRegistry()
		}

		return opts.encodig.Register(contentType, codec)
	}
}
