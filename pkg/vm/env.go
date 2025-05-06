package vm

import (
	"os"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/logging"
)

type (
	EnvironmentOption func(env *Environment)

	Environment struct {
		functions map[string]runtime.Function
		params    map[string]runtime.Value
		logging   logging.Options
	}
)

var noopEnv = &Environment{
	functions: make(map[string]runtime.Function),
	params:    make(map[string]runtime.Value),
}

func newEnvironment(opts []EnvironmentOption) *Environment {
	if len(opts) == 0 {
		return noopEnv
	}

	env := &Environment{
		functions: make(map[string]runtime.Function),
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

func (env *Environment) GetFunction(name string) runtime.Function {
	return env.functions[name]
}

func (env *Environment) HasFunction(name string) bool {
	_, exists := env.functions[name]

	return exists
}

func (env *Environment) GetParam(name string) runtime.Value {
	return env.params[name]
}

func (env *Environment) HasParam(name string) bool {
	_, exists := env.params[name]

	return exists
}
