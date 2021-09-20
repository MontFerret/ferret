package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type SuppressibleExpression struct {
	exp core.Expression
}

func SuppressErrors(exp core.Expression) (core.Expression, error) {
	if exp == nil {
		return nil, core.Error(core.ErrMissedArgument, "expression")
	}

	return &SuppressibleExpression{exp}, nil
}

func (exp *SuppressibleExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	return exp.Maybe(exp.exp.Exec(ctx, scope))
}

func (exp *SuppressibleExpression) Maybe(value core.Value, err error) (core.Value, error) {
	if err != nil {
		return values.None, nil
	}

	return value, nil
}
