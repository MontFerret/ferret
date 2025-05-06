package vm

import (
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/logging"
)

func WithParams(params map[string]runtime.Value) EnvironmentOption {
	return func(env *Environment) {
		env.params = params
	}
}

func WithParam(name string, value interface{}) EnvironmentOption {
	return func(options *Environment) {
		options.params[name] = runtime.Parse(value)
	}
}

func WithFunctions(functions map[string]runtime.Function) EnvironmentOption {
	return func(env *Environment) {
		if env.functions == nil {
			env.functions = make(map[string]runtime.Function)
		}

		for name, function := range functions {
			env.functions[name] = function
		}
	}
}

func WithFunction(name string, function runtime.Function) EnvironmentOption {
	return func(env *Environment) {
		env.functions[name] = function
	}
}

func WithLog(writer io.Writer) EnvironmentOption {
	return func(options *Environment) {
		options.logging.Writer = writer
	}
}

func WithLogLevel(lvl logging.Level) EnvironmentOption {
	return func(options *Environment) {
		options.logging.Level = lvl
	}
}

func WithLogFields(fields map[string]interface{}) EnvironmentOption {
	return func(options *Environment) {
		options.logging.Fields = fields
	}
}
