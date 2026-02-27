package vm

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func WithParams(params map[string]runtime.Value) EnvironmentOption {
	return func(env *environmentBuilder) {
		if params != nil {
			env.params = params
		}
	}
}

func WithParam(name string, value interface{}) EnvironmentOption {
	return func(env *environmentBuilder) {
		if env.params == nil {
			env.params = make(map[string]runtime.Value)
		}

		env.params[name] = runtime.Parse(value)
	}
}

func WithFunctions(funcs *runtime.Functions) EnvironmentOption {
	return func(env *environmentBuilder) {
		if funcs != nil {
			env.functions.From(runtime.NewFunctionsBuilderFrom(funcs))
		}
	}
}

func WithFunction(name string, function runtime.Function) EnvironmentOption {
	return func(env *environmentBuilder) {
		if name != "" && function != nil {
			env.functions.Var().Add(name, function)
		}
	}
}

func WithNamespace(ns runtime.Namespace) EnvironmentOption {
	return func(env *environmentBuilder) {
		if ns != nil {
			env.functions.From(ns.Function())
		}
	}
}

func WithFunctionsBuilder(builder runtime.FunctionDefs) EnvironmentOption {
	return func(env *environmentBuilder) {
		if builder != nil {
			env.functions.From(builder)
		}
	}
}

func WithFunctionsRegistrar(setter func(fns runtime.FunctionDefs)) EnvironmentOption {
	return func(env *environmentBuilder) {
		if setter != nil {
			setter(env.functions)
		}
	}
}

func WithLog(writer io.Writer) EnvironmentOption {
	return func(env *environmentBuilder) {
		env.logging.Writer = writer
	}
}

func WithLogLevel(lvl runtime.LogLevel) EnvironmentOption {
	return func(env *environmentBuilder) {
		env.logging.Level = lvl
	}
}

func WithLogFields(fields map[string]interface{}) EnvironmentOption {
	return func(env *environmentBuilder) {
		env.logging.Fields = fields
	}
}
