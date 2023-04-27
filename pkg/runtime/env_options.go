package runtime

import "github.com/MontFerret/ferret/pkg/runtime/core"

func WithParams(params []core.Value) EnvironmentOption {
	return func(env *Environment) {
		env.params = params
	}
}

func WithFunctions(functions map[string]core.Function) EnvironmentOption {
	return func(env *Environment) {
		env.functions = functions
	}
}
