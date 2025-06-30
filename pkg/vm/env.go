package vm

import (
	"os"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/logging"
)

type (
	EnvironmentOption func(env *Environment)

	Environment struct {
		functions runtime.Functions
		params    map[string]runtime.Value
		logging   logging.Options
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
		logging: logging.Options{
			Writer: os.Stdout,
			Level:  logging.ErrorLevel,
		},
	}

	for _, opt := range opts {
		opt(env)
	}

	return env
}
