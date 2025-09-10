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

func WithFunctions(funcs runtime.Functions) EnvironmentOption {
	return func(env *Environment) {
		if funcs != nil {
			env.Functions = runtime.NewFunctionsFrom(env.Functions, funcs)
		}
	}
}

func WithFunction(name string, function runtime.Function) EnvironmentOption {
	return func(env *Environment) {
		env.Functions = runtime.NewFunctionsBuilder().SetFrom(env.Functions).Set(name, function).Build()
	}
}

func WithNamespace(ns runtime.Namespace) EnvironmentOption {
	return func(env *Environment) {
		if ns != nil {
			env.Functions = runtime.NewFunctionsFrom(env.Functions, ns.Functions().Build())
		}
	}
}

func WithFunctionsBuilder(setter func(fns runtime.FunctionsBuilder)) EnvironmentOption {
	return func(env *Environment) {
		if setter != nil {
			builder := runtime.NewFunctionsBuilder()
			setter(builder)
			builder.SetFrom(env.Functions)

			env.Functions = builder.Build()
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
