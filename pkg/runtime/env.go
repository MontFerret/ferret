package runtime

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	EnvironmentOption func(env *Environment)

	Environment struct {
		functions map[string]core.Function
		params    []core.Value
	}
)

var noopEnv = &Environment{
	functions: make(map[string]core.Function),
	params:    make([]core.Value, 0),
}

func newEnvironment(opts []EnvironmentOption) *Environment {
	if len(opts) == 0 {
		return noopEnv
	}

	env := &Environment{
		functions: make(map[string]core.Function),
		params:    make([]core.Value, 0),
	}

	for _, opt := range opts {
		opt(env)
	}

	return env
}

func (env *Environment) GetFunction(name string) core.Function {
	return env.functions[name]
}
