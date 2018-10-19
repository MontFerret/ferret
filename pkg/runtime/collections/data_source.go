package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	DataSource interface {
		Iterate(ctx context.Context, scope *core.Scope) (Iterator, error)
	}
)
