package literals

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type FloatLiteral float64

func NewFloatLiteral(value float64) FloatLiteral {
	return FloatLiteral(value)
}

func (l FloatLiteral) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return values.NewFloat(float64(l)), nil
}
