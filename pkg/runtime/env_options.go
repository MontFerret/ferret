package runtime

import "github.com/MontFerret/ferret/pkg/runtime/core"

func WithParams(params []core.Value) EnvironmentOption {
	return func(env *Environment) {
		env.params = params
	}
}

func WithFunctions(functions map[string]core.Function) EnvironmentOption {
	return func(env *Environment) {
		if env.functions == nil {
			env.functions = make(map[string]core.Function)
		}

		for name, function := range functions {
			env.functions[name] = function
		}
	}
}

func WithFunction(name string, function core.Function) EnvironmentOption {
	return func(env *Environment) {
		if env.functions == nil {
			env.functions = make(map[string]core.Function)
		}

		env.functions[name] = function
	}
}
