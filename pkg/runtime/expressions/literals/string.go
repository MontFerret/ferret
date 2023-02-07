package literals

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type StringLiteral string

func NewStringLiteral(str string) StringLiteral {
	return StringLiteral(str)
}

func (l StringLiteral) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return values.NewString(string(l)), nil
}
