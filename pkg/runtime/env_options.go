package runtime

import (
	"io"

	"github.com/MontFerret/ferret/pkg/logging"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func WithParams(params map[string]core.Value) EnvironmentOption {
	return func(env *Environment) {
		env.params = params
	}
}

func WithParam(name string, value interface{}) EnvironmentOption {
	return func(options *Environment) {
		options.params[name] = core.Parse(value)
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
