package base

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	FnState[T any] interface {
		Get() T
		Set(value T)
	}

	fnState[T any] struct {
		value T
	}
)

func (f *fnState[T]) Get() T {
	return f.value
}

func (f *fnState[T]) Set(value T) {
	f.value = value
}

var globalState map[context.Context]any

func init() {
	globalState = make(map[context.Context]any)
}

func StateFn[T any](fn runtime.Function, factory func(ctx context.Context) T) runtime.Function {
	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		_, exists := globalState[ctx]

		if !exists {
			val := factory(ctx)
			globalState[ctx] = &fnState[T]{val}
		}

		return fn(ctx, args...)
	}
}

func GetFnState[T any](ctx context.Context) FnState[T] {
	state, exists := globalState[ctx]

	if !exists {
		panic("fn state is not initialized")
	}

	stateT, ok := state.(FnState[T])

	if !ok {
		panic("fn state has wrong type")
	}

	return stateT
}

func UpdateFnState[T any](ctx context.Context, updater func(current T) T) {
	state := GetFnState[T](ctx)
	state.Set(updater(state.Get()))
}
