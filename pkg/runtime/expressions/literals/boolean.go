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

func (l BooleanLiteral) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	if l {
		return values.True, nil
	}

	return values.False, nil
}
