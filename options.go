package ferret

import (
	"fmt"
	"io"
	"os"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	options struct {
		noStdlib bool
		compiler []compiler.Option
		library  runtime.Library
		params   runtime.Params
		logging  runtime.LogSettings
		encoding *encoding.Registry
		modules  []Module
		hooks    *hookRegistry
	}

	Option func(env *options) error

	SessionOption = vm.EnvironmentOption
)

var (
	WithSessionParams = vm.WithParams
	WithSessionParam  = vm.WithParam
)

func newOptions(setters []Option) (*options, error) {
	opts := &options{
		compiler: []compiler.Option{},
		library:  runtime.NewLibrary(),
		params:   make(map[string]runtime.Value),
		encoding: encoding.NewRegistry(),
		hooks:    newHookRegistry(),
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
		stdlib.RegisterLib(opts.library)
	}

	return opts, nil
}

// WithNamespace creates an Option that sets the library from the provided runtime.Namespace to the options if not nil.
func WithNamespace(ns runtime.Namespace) Option {
	return func(opts *options) error {
		if ns == nil {
			return fmt.Errorf("namespace cannot be nil")
		}

		opts.library.Function().From(ns.Function())

		return nil
	}
}

// WithFunctionsRegistrar creates an Option that invokes the provided registrar with the engine's runtime.FunctionDefs if the registrar is not nil.
func WithFunctionsRegistrar(setter func(fns runtime.FunctionDefs)) Option {
	return func(env *options) error {
		if setter == nil {
			return fmt.Errorf("functions registrar cannot be nil")
		}

		setter(env.library.Function())

		return nil
	}
}

// WithFunctions creates an Option that sets the provided *runtime.Functions to the options if not nil.
func WithFunctions(funcs *runtime.Functions) Option {
	return func(opts *options) error {
		if funcs == nil {
			return fmt.Errorf("functions cannot be nil")
		}

		opts.library.Function().From(runtime.NewFunctionsBuilderFrom(funcs))

		return nil
	}
}

// WithLog sets the writer for logging output.
// The writer can be any io.Writer, such as os.Stdout or a file.
func WithLog(writer io.Writer) Option {
	return func(opts *options) error {
		if writer == nil {
			return fmt.Errorf("log writer cannot be nil")
		}

		opts.logging.Writer = writer
		return nil
	}
}

// WithLogLevel sets the logging level for the engine.
// The logging level determines the severity of log messages that will be recorded.
func WithLogLevel(lvl runtime.LogLevel) Option {
	return func(opts *options) error {
		if lvl < runtime.TraceLevel || lvl > runtime.Disabled {
			return fmt.Errorf("invalid log level: %v", lvl)
		}

		opts.logging.Level = lvl
		return nil
	}
}

// WithLogFields sets the fields to be included in log entries.
// These fields can provide additional context for debugging and monitoring purposes.
func WithLogFields(fields map[string]any) Option {
	return func(opts *options) error {
		if fields == nil {
			return fmt.Errorf("log fields cannot be nil")
		}

		if opts.logging.Fields == nil {
			opts.logging.Fields = make(map[string]any)
		}

		for k, v := range fields {
			opts.logging.Fields[k] = v
		}

		return nil
	}
}

// WithEncodingRegistry sets a custom encoding registry for query execution.
func WithEncodingRegistry(registry *encoding.Registry) Option {
	return func(opts *options) error {
		if registry == nil {
			return fmt.Errorf("encoding registry is nil")
		}

		opts.encoding = registry
		return nil
	}
}

// WithModules creates an Option that appends the provided modules to the options if not empty.
func WithModules(module ...Module) Option {
	return func(env *options) error {
		if len(module) == 0 {
			return fmt.Errorf("modules cannot be empty")
		}

		if env.modules == nil {
			env.modules = make([]Module, 0, len(module))
		}

		env.modules = append(env.modules, module...)

		return nil
	}
}

// WithEncodingCodec registers or overrides a codec for the given content type.
func WithEncodingCodec(contentType string, codec encoding.Codec) Option {
	return func(opts *options) error {
		if opts.encoding == nil {
			opts.encoding = encoding.NewRegistry()
		}

		return opts.encoding.Register(contentType, codec)
	}
}

// WithCompilerOptions creates an Option that appends the provided compiler options to the options if not empty.
func WithCompilerOptions(opts ...compiler.Option) Option {
	return func(o *options) error {
		if len(opts) == 0 {
			return nil
		}

		if o.compiler == nil {
			o.compiler = opts
			return nil
		}

		o.compiler = append(o.compiler, opts...)

		return nil
	}
}

// WithEngineInitHook returns an Option that registers a hook to be executed during engine initialization.
func WithEngineInitHook(hook EngineInitHook) Option {
	return func(opts *options) error {
		if hook == nil {
			return fmt.Errorf("engine init hook is nil")
		}

		opts.hooks.engine.OnInit(hook)
		return nil
	}
}

// WithEngineCloseHook returns an Option that registers a hook to be executed when the engine is closed.
func WithEngineCloseHook(hook EngineCloseHook) Option {
	return func(opts *options) error {
		if hook == nil {
			return fmt.Errorf("engine close hook is nil")
		}

		opts.hooks.engine.OnClose(hook)
		return nil
	}
}

// WithBeforeCompileHook returns an Option that registers a hook to be executed before compilation.
func WithBeforeCompileHook(hook BeforeCompileHook) Option {
	return func(opts *options) error {
		if hook == nil {
			return fmt.Errorf("before compile hook is nil")
		}

		opts.hooks.plan.BeforeCompile(hook)
		return nil
	}
}

func WithAfterCompileHook(hook AfterCompileHook) Option {
	return func(opts *options) error {
		if hook == nil {
			return fmt.Errorf("after compile hook is nil")
		}

		opts.hooks.plan.AfterCompile(hook)
		return nil
	}
}

// WithPlanCloseHook returns an Option that registers a hook to be executed when a plan is closed.
func WithPlanCloseHook(hook PlanCloseHook) Option {
	return func(opts *options) error {
		if hook == nil {
			return fmt.Errorf("plan close hook is nil")
		}

		opts.hooks.plan.OnClose(hook)
		return nil
	}
}

// WithBeforeRunHook returns an Option that registers a hook to be executed during session initialization.
func WithBeforeRunHook(hook BeforeRunHook) Option {
	return func(opts *options) error {
		if hook == nil {
			return fmt.Errorf("before run hook is nil")
		}

		opts.hooks.session.BeforeRun(hook)
		return nil
	}
}

// WithAfterRunHook returns an Option that registers a hook to be executed during session initialization.
func WithAfterRunHook(hook AfterRunHook) Option {
	return func(opts *options) error {
		if hook == nil {
			return fmt.Errorf("after run hook is nil")
		}

		opts.hooks.session.AfterRun(hook)
		return nil
	}
}

// WithSessionCloseHook returns an Option that registers a hook to be executed when a session is closed.
func WithSessionCloseHook(hook SessionCloseHook) Option {
	return func(opts *options) error {
		if hook == nil {
			return fmt.Errorf("session close hook is nil")
		}

		opts.hooks.session.OnClose(hook)
		return nil
	}
}
