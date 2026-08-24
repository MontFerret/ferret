package vm

import (
	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	environmentConfig struct {
		functions *runtime.FunctionsBuilder
		params    runtime.Params
	}

	// EnvironmentOption configures an Environment during construction.
	EnvironmentOption = options.Option[environmentConfig]

	Environment struct {
		Functions *runtime.Functions
		Params    runtime.Params
	}
)

var noopEnv = &Environment{
	Functions: runtime.NewFunctions(),
	Params:    runtime.NewParams(),
}

func NewDefaultEnvironment() *Environment {
	return &Environment{
		Functions: runtime.NewFunctions(),
		Params:    runtime.NewParams(),
	}
}

func defaultEnvironmentConfig() environmentConfig {
	return environmentConfig{
		functions: runtime.NewFunctionsBuilder(),
		params:    runtime.NewParams(),
	}
}

// NewEnvironment creates an environment from opts. Option errors are returned
// before the configured function registry is built.
func NewEnvironment(opts []EnvironmentOption) (*Environment, error) {
	cfg, err := options.ApplyTo(defaultEnvironmentConfig(), opts...)
	if err != nil {
		return nil, err
	}

	funcs, err := cfg.functions.Build()
	if err != nil {
		return nil, err
	}

	return &Environment{
		Functions: funcs,
		Params:    cfg.params,
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
	}

	builder := runtime.NewFunctionsBuilderFrom(funcsToMerge...)
	funcs, err := builder.Build()
	if err != nil {
		return nil, err
	}

	merged.Functions = funcs

	return merged, nil
}

// ExtendEnvironment applies opts to a new environment and merges it into env.
func ExtendEnvironment(env *Environment, opts []EnvironmentOption) (*Environment, error) {
	if len(opts) == 0 {
		return env, nil
	}

	newEnv, err := NewEnvironment(opts)
	if err != nil {
		return nil, err
	}

	return MergeEnvironments(env, newEnv)
}
