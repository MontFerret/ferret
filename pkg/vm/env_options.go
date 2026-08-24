package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func WithParams(params runtime.Params) EnvironmentOption {
	return func(env *environmentConfig) error {
		if params == nil {
			return nil
		}

		if env.params == nil {
			env.params = runtime.NewParams()
		}

		env.params.MergeParams(params)

		return nil
	}
}

func WithParam(name string, value runtime.Value) EnvironmentOption {
	return func(env *environmentConfig) error {
		if env.params == nil {
			env.params = runtime.NewParams()
		}

		env.params[name] = value

		return nil
	}
}

func WithFunctions(funcs *runtime.Functions) EnvironmentOption {
	return func(env *environmentConfig) error {
		if funcs != nil {
			env.functions.From(runtime.NewFunctionsBuilderFrom(funcs))
		}

		return nil
	}
}

func WithFunction(name string, function runtime.Function) EnvironmentOption {
	return func(env *environmentConfig) error {
		if name != "" && function != nil {
			env.functions.Var().Add(name, function)
		}

		return nil
	}
}

func WithNamespace(ns runtime.Namespace) EnvironmentOption {
	return func(env *environmentConfig) error {
		if ns != nil {
			env.functions.From(ns.Function())
		}

		return nil
	}
}

func WithFunctionsBuilder(builder runtime.FunctionDefs) EnvironmentOption {
	return func(env *environmentConfig) error {
		if builder != nil {
			env.functions.From(builder)
		}

		return nil
	}
}

func WithFunctionsRegistrar(setter func(fns runtime.FunctionDefs)) EnvironmentOption {
	return func(env *environmentConfig) error {
		if setter != nil {
			setter(env.functions)
		}

		return nil
	}
}
