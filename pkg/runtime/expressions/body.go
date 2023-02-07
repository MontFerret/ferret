package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type BodyExpression struct {
	statements []core.Expression
	expression core.Expression
}

func NewBodyExpression(size int) *BodyExpression {
	return &BodyExpression{make([]core.Expression, 0, size), nil}
}

func (b *BodyExpression) Add(exp core.Expression) error {
	switch exp.(type) {
	case *ForExpression, *ReturnExpression:
		if b.expression != nil {
			return core.Error(core.ErrInvalidOperation, "return expression is already defined")
		}

		b.expression = exp
	default:
		b.statements = append(b.statements, exp)
	}

	return nil
}

func (b *BodyExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	select {
	case <-ctx.Done():
		return values.None, core.ErrTerminated
	default:
	}

	for _, exp := range b.statements {
		if _, err := exp.Exec(ctx, scope); err != nil {
			return values.None, err
		}
	}

	if b.expression != nil {
		return b.expression.Exec(ctx, scope)
	}

	return values.None, nil
}
