package vm

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func createContext(base context.Context, env *Environment) runtime.Context {
	alloc := runtime.NewAllocator()

	if env == nil {
		return runtime.NewContext(base, runtime.NewLogger(runtime.LogSettings{}), runtime.NewAllocator())
	}

	return runtime.NewContext(base, runtime.NewLogger(env.Logging), alloc)
}
