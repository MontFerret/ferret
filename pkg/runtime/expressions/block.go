package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

type BlockExpression struct {
	statements []core.Expression
	expression core.Expression
}

func NewBlockExpression(size int) *BlockExpression {
	return &BlockExpression{make([]core.Expression, 0, size), nil}
}

func NewBlockExpressionWith(elements ...core.Expression) *BlockExpression {
	block := NewBlockExpression(len(elements))

	for _, el := range elements {
		block.Add(el)
	}

	return block
}

func (b *BlockExpression) Add(exp core.Expression) error {
	switch exp.(type) {
	case *ForExpression, *ReturnExpression:
		// return an error?
		if !core.IsNil(b.expression) {
			return errors.Wrap(core.ErrInvalidOperation, "return expression is already defined")
		}

		b.expression = exp

		break
	default:
		b.statements = append(b.statements, exp)
	}

	return nil
}

func (b *BlockExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	for _, exp := range b.statements {
		if _, err := exp.Exec(ctx, scope); err != nil {
			return values.None, err
		}
	}

	if !core.IsNil(b.expression) {
		return b.expression.Exec(ctx, scope)
	}

	return values.None, nil
}
