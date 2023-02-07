package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type BlockExpression struct {
	values     collections.Iterable
	statements []core.Expression
}

func NewBlockExpression(values collections.Iterable) (*BlockExpression, error) {
	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "values")
	}

	return &BlockExpression{
		values:     values,
		statements: make([]core.Expression, 0, 5),
	}, nil
}

func (exp *BlockExpression) Add(stmt core.Expression) {
	exp.statements = append(exp.statements, stmt)
}

func (exp *BlockExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	select {
	case <-ctx.Done():
		return values.None, core.ErrTerminated
	default:
		for _, stmt := range exp.statements {
			_, err := stmt.Exec(ctx, scope)

			if err != nil {
				return values.None, err
			}
		}

		return values.None, nil
	}
}

func (exp *BlockExpression) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	select {
	case <-ctx.Done():
		return nil, core.ErrTerminated
	default:
		iter, err := exp.values.Iterate(ctx, scope)

		if err != nil {
			return nil, err
		}

		return collections.NewTapIterator(iter, exp)
	}
}
