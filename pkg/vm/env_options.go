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

func WithFunctions(functions runtime.Functions) EnvironmentOption {
	return func(env *Environment) {
		env.functions.SetAll(functions)
	}
}

func WithFunction(name string, function runtime.Function) EnvironmentOption {
	return func(env *Environment) {
		env.functions.F().Set(name, function)
	}
}

func WithFunctionSetter(setter func(fns runtime.Functions)) EnvironmentOption {
	return func(env *Environment) {
		if setter != nil {
			setter(env.functions)
		}
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
