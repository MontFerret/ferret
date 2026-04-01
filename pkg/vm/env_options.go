package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func WithParams(params runtime.Params) EnvironmentOption {
	return func(env *environmentBuilder) {
		if params == nil {
			return
		}

		if env.params == nil {
			env.params = runtime.NewParams()
		}

		env.params.MergeParams(params)
	}
}

func WithParam(name string, value runtime.Value) EnvironmentOption {
	return func(env *environmentBuilder) {
		if env.params == nil {
			env.params = runtime.NewParams()
		}

		env.params[name] = value
	}
}

func WithFunctions(funcs *runtime.Functions) EnvironmentOption {
	return func(env *environmentBuilder) {
		if funcs != nil {
			env.functions.From(runtime.NewFunctionsBuilderFrom(funcs))
		}
	}
}

func WithFunction(name string, function runtime.Function) EnvironmentOption {
	return func(env *environmentBuilder) {
		if name != "" && function != nil {
			env.functions.Var().Add(name, function)
		}
	}
}

func WithNamespace(ns runtime.Namespace) EnvironmentOption {
	return func(env *environmentBuilder) {
		if ns != nil {
			env.functions.From(ns.Function())
		}
	}
}

func WithFunctionsBuilder(builder runtime.FunctionDefs) EnvironmentOption {
	return func(env *environmentBuilder) {
		if builder != nil {
			env.functions.From(builder)
		}
	}
}

func WithFunctionsRegistrar(setter func(fns runtime.FunctionDefs)) EnvironmentOption {
	return func(env *environmentBuilder) {
		if setter != nil {
			setter(env.functions)
		}
	}
}
