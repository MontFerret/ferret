package vm

import (
	"os"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	environmentBuilder struct {
		functions *runtime.FunctionsBuilder
		params    map[string]runtime.Value
		logging   runtime.LogSettings
	}

	EnvironmentOption func(env *environmentBuilder)

	Environment struct {
		Functions *runtime.Functions
		Params    map[string]runtime.Value
		Logging   runtime.LogSettings
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
		Logging: runtime.LogSettings{
			Writer: os.Stdout,
			Level:  runtime.ErrorLevel,
		},
	}
}

func NewEnvironment(opts []EnvironmentOption) (*Environment, error) {
	envBuilder := &environmentBuilder{
		functions: runtime.NewFunctionsBuilder(),
		params:    make(map[string]runtime.Value),
		logging: runtime.LogSettings{
			Writer: os.Stdout,
			Level:  runtime.ErrorLevel,
		},
	}

	for _, opt := range opts {
		opt(envBuilder)
	}

	funcs, err := envBuilder.functions.Build()

	if err != nil {
		return nil, err
	}

	return &Environment{
		Functions: funcs,
		Params:    envBuilder.params,
		Logging:   envBuilder.logging,
	}, nil
}

func MergeEnvironments(envs ...*Environment) (*Environment, error) {
	if len(envs) == 0 {
		return NewDefaultEnvironment(), nil
	}

	if len(envs) == 1 {
		return envs[0], nil
	}

	merged := NewDefaultEnvironment()
	funcsToMerge := make([]*runtime.Functions, 0, len(envs))

	for _, env := range envs {
		if env == nil {
			continue
		}

		// merge functions
		// TODO: Resolve conflicts between functions with the same name?
		funcsToMerge = append(funcsToMerge, env.Functions)

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

	builder := runtime.NewFunctionsBuilderFrom(funcsToMerge...)
	funcs, err := builder.Build()

	if err != nil {
		return nil, err
	}

	merged.Functions = funcs

	return merged, nil
}
