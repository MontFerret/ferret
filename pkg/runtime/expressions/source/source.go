package source

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	DataSource interface {
		Variables() []string
		Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error)
	}
)
