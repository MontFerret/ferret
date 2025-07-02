package vm

import (
	"os"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	EnvironmentOption func(env *Environment)

	Environment struct {
		functions runtime.Functions
		params    map[string]runtime.Value
		logging   runtime.Options
	}
)

var noopEnv = &Environment{
	functions: runtime.NewFunctions(),
	params:    make(map[string]runtime.Value),
}

func newEnvironment(opts []EnvironmentOption) *Environment {
	if len(opts) == 0 {
		return noopEnv
	}

	env := &Environment{
		functions: runtime.NewFunctions(),
		params:    make(map[string]runtime.Value),
		logging: runtime.Options{
			Writer: os.Stdout,
			Level:  runtime.ErrorLevel,
		},
	}

	for _, opt := range opts {
		opt(env)
	}

	return env
}
