package engine

import (
	"fmt"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

// WithParams applies custom parameters to the options by merging them with existing ones, initializing if necessary.
// If a parameter already exists, it will be overwritten.
// All host values will be converted to a runtime.Value.
func WithParams(params map[string]any) Option {
	return func(opts *config) error {
		if len(params) == 0 {
			return nil
		}

		if opts.params == nil {
			opts.params = runtime.NewParams()
		}

		merged, err := opts.params.Merge(params)
		if err != nil {
			return err
		}

		opts.params = merged

		return nil
	}
}

// WithRuntimeParams configures runtime parameters by merging the provided runtime.Params with existing ones in options.
// If a parameter already exists, it will be overwritten.
func WithRuntimeParams(params runtime.Params) Option {
	return func(opts *config) error {
		if len(params) == 0 {
			return nil
		}

		if opts.params == nil {
			opts.params = runtime.NewParams()
		}

		opts.params = opts.params.MergeParams(params)

		return nil
	}
}

// WithParam returns an Option that sets a parameter with the specified name and value in the options configuration.
// The name cannot be empty, and the value cannot be nil. It ensures the parameter value is correctly parsed and stored.
func WithParam(name string, value any) Option {
	return func(opts *config) error {
		if name == "" {
			return fmt.Errorf("param name cannot be empty")
		}

		if value == nil {
			return fmt.Errorf("param value cannot be nil")
		}

		if opts.params == nil {
			opts.params = runtime.NewParams()
		}

		parsed, err := runtime.ValueOf(value)
		if err != nil {
			return fmt.Errorf("invalid param value: %w", err)
		}

		opts.params.SetValue(name, parsed)

		return nil
	}
}

// WithRuntimeParam returns an Option that sets a runtime parameter with the specified name and pre-converted runtime.Value.
// The name cannot be empty, and the value cannot be nil.
func WithRuntimeParam(name string, value runtime.Value) Option {
	return func(opts *config) error {
		if name == "" {
			return fmt.Errorf("param name cannot be empty")
		}

		if value == nil {
			return fmt.Errorf("param value cannot be nil")
		}

		if opts.params == nil {
			opts.params = runtime.NewParams()
		}

		opts.params.SetValue(name, value)

		return nil
	}
}

// WithNamespace merges the functions from the provided runtime.Namespace into the engine's function library.
func WithNamespace(ns runtime.Namespace) Option {
	return func(opts *config) error {
		if ns == nil {
			return fmt.Errorf("namespace cannot be nil")
		}

		opts.library.Function().From(ns.Function())

		return nil
	}
}

// WithFunctionsRegistrar creates an Option that invokes the provided registrar with the engine's runtime.Namespace if the registrar is not nil.
// Registered host-function names and namespace segments are canonicalized to lowercase and resolve case-insensitively in FQL.
func WithFunctionsRegistrar(setter func(ns runtime.Namespace)) Option {
	return func(env *config) error {
		if setter == nil {
			return fmt.Errorf("functions registrar cannot be nil")
		}

		setter(env.library)

		return nil
	}
}

// WithFunctions merges the provided *runtime.Functions into the engine's function library.
func WithFunctions(funcs *runtime.Functions) Option {
	return func(opts *config) error {
		if funcs == nil {
			return fmt.Errorf("functions cannot be nil")
		}

		opts.library.Function().From(runtime.NewFunctionsBuilderFrom(funcs))

		return nil
	}
}

// WithoutStdlib disables the standard library, so no built-in functions are registered by default.
func WithoutStdlib() Option {
	return WithStdlib(stdlib.Empty())
}

// WithStdlib configures which standard library groups are registered by default.
func WithStdlib(set stdlib.Set) Option {
	return gooptions.New(func(opts *config, set stdlib.Set) {
		opts.stdlib = set
	}).
		Value(set).
		Named("stdlib").
		Build()
}

// WithModules creates an Option that appends the provided module.Module values if not empty.
func WithModules(mods ...module.Module) Option {
	return func(env *config) error {
		if len(mods) == 0 {
			return nil
		}

		if env.modules == nil {
			env.modules = make([]module.Module, 0, len(mods))
		}

		for _, m := range mods {
			if m == nil {
				return fmt.Errorf("module cannot be nil")
			}

			env.modules = append(env.modules, m)
		}

		return nil
	}
}
