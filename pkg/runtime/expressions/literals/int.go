package literals

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type IntLiteral int

func NewIntLiteral(value int) IntLiteral {
	return IntLiteral(value)
}

func (l IntLiteral) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return values.NewInt(int(l)), nil
}
