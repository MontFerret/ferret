package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Reducer interface {
		Reduce(collection core.Value, value core.Value) (core.Value, error)
	}

	Reducible interface {
		Reduce() Reducer
	}

	ReducibleExpression interface {
		core.Expression
		Reduce(ctx context.Context, scope *core.Scope) (Reducer, error)
	}
)
