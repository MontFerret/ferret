package vm

import (
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func WithParams(params map[string]runtime.Value) EnvironmentOption {
	return func(env *Environment) {
		env.Params = params
	}
}

func WithParam(name string, value interface{}) EnvironmentOption {
	return func(options *Environment) {
		options.Params[name] = runtime.Parse(value)
	}
}

func WithFunctions(functions runtime.Functions) EnvironmentOption {
	return func(env *Environment) {
		env.Functions.SetAll(functions)
	}
}

func WithFunction(name string, function runtime.Function) EnvironmentOption {
	return func(env *Environment) {
		env.Functions.F().Set(name, function)
	}
}

func WithFunctionSetter(setter func(fns runtime.Functions)) EnvironmentOption {
	return func(env *Environment) {
		if setter != nil {
			setter(env.Functions)
		}
	}
}

func WithLog(writer io.Writer) EnvironmentOption {
	return func(options *Environment) {
		options.Logging.Writer = writer
	}
}

func WithLogLevel(lvl runtime.LogLevel) EnvironmentOption {
	return func(options *Environment) {
		options.Logging.Level = lvl
	}
}

func WithLogFields(fields map[string]interface{}) EnvironmentOption {
	return func(options *Environment) {
		options.Logging.Fields = fields
	}
}
