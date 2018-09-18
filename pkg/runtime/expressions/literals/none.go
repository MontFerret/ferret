package literals

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type noneLiteral struct{}

var None = &noneLiteral{}

func (l noneLiteral) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	return values.None, nil
}
