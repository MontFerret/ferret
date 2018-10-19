package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type ExpressionDataSource struct {
	exp core.Expression
}

func NewExpressionDataSource(exp core.Expression) (DataSource, error) {
	if exp == nil {
		return nil, core.Error(core.ErrMissedArgument, "expression")
	}

	return &ExpressionDataSource{exp}, nil
}

func (ds *ExpressionDataSource) Iterate(ctx context.Context, scope *core.Scope) (Iterator, error) {
	data, err := ds.exp.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return ToIterator(data)
}
