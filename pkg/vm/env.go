package vm

import (
	"os"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	EnvironmentOption func(env *Environment)

	Environment struct {
		Functions runtime.Functions
		Params    map[string]runtime.Value
		Logging   runtime.LogOptions
	}
)

var noopEnv = &Environment{
	Functions: runtime.NewFunctions(),
	Params:    make(map[string]runtime.Value),
}

func NewDefaultEnvironment() *Environment {
	return &Environment{
		Functions: runtime.NewFunctions(),
		Params:    make(map[string]runtime.Value),
		Logging: runtime.LogOptions{
			Writer: os.Stdout,
			Level:  runtime.ErrorLevel,
		},
	}
}

func NewEnvironment(opts []EnvironmentOption) *Environment {
	env := NewDefaultEnvironment()

	for _, opt := range opts {
		opt(env)
	}

	return env
}

func MergeEnvironments(envs ...*Environment) *Environment {
	if len(envs) == 0 {
		return NewDefaultEnvironment()
	}

	if len(envs) == 1 {
		return envs[0]
	}

	merged := NewDefaultEnvironment()

	for _, env := range envs {
		if env == nil {
			continue
		}

		// merge functions
		merged.Functions.SetAll(env.Functions)

		// merge params
		for name, val := range env.Params {
			merged.Params[name] = val
		}

		// merge logging options
		if env.Logging.Writer != nil {
			merged.Logging.Writer = env.Logging.Writer
		}

		merged.Logging.Level = env.Logging.Level
	}

	return merged
}
