package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type SessionOption = vm.EnvironmentOption

func WithSessionParams(params map[string]runtime.Value) SessionOption {
	return func(env *vm.Environment) {
		if params == nil {
			return
		}

		env.Params = params
	}
}

func WithSessionParam(name string, value interface{}) SessionOption {
	return func(options *vm.Environment) {
		options.Params[name] = runtime.Parse(value)
	}
}
