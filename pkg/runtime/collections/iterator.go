package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Iterator interface {
		Next(ctx context.Context, scope *core.Scope) (*core.Scope, error)
	}

	Iterable interface {
		Iterate(ctx context.Context, scope *core.Scope) (Iterator, error)
	}
)
