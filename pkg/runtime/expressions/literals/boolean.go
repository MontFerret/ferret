package literals

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type BooleanLiteral bool

func NewBooleanLiteral(val bool) BooleanLiteral {
	return BooleanLiteral(val)
}

func (l BooleanLiteral) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	if l == true {
		return values.True, nil
	}

	return values.False, nil
}
