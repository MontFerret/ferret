package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type FilterClause struct {
	src        core.SourceMap
	dataSource collections.Iterable
	predicate  core.Expression
	variables  collections.Variables
}

func NewFilterClause(
	src core.SourceMap,
	dataSource collections.Iterable,
	predicate core.Expression,
) (collections.Iterable, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "dataSource source")
	}

	if predicate == nil {
		return nil, core.Error(core.ErrMissedArgument, "predicate")
	}

	return &FilterClause{
		src, dataSource,
		predicate,
		dataSource.Variables(),
	}, nil
}

func (clause *FilterClause) Variables() collections.Variables {
	return clause.dataSource.Variables()
}

func (clause *FilterClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewFilterIterator(src, clause.filter)
}

func (clause *FilterClause) filter(ctx context.Context, scope *core.Scope, set collections.DataSet) (bool, error) {
	innerScope := scope.Fork()

	err := set.Apply(innerScope, clause.variables)

	if err != nil {
		return false, core.SourceError(clause.src, err)
	}

	ret, err := clause.predicate.Exec(ctx, innerScope)

	if err != nil {
		return false, err
	}

	if ret == values.True {
		return true, nil
	}

	return false, nil
}
