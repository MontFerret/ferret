package source

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type ExpressionDataSource struct {
	variables Variables
	exp       core.Expression
}

func NewExpressionDataSource(
	variables Variables,
	exp core.Expression,
) (DataSource, error) {
	if variables == nil {
		return nil, core.Error(core.ErrMissedArgument, "variables")
	}

	if len(variables) == 0 {
		return nil, core.Error(core.ErrInvalidArgumentNumber, "variables array expected to be have at least one item")
	}

	if exp == nil {
		return nil, core.Error(core.ErrMissedArgument, "expression")
	}

	return &ExpressionDataSource{variables, exp}, nil
}

func (ds *ExpressionDataSource) Variables() Variables {
	return ds.variables
}

func (ds *ExpressionDataSource) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	data, err := ds.exp.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.ToIterator(data)
}
