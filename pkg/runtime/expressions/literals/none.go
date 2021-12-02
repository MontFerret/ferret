package literals

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type noneLiteral struct{}

var None = &noneLiteral{}

func (l noneLiteral) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return values.None, nil
}

func IsNone(exp core.Expression) bool {
	_, is := exp.(*noneLiteral)

	return is
}
