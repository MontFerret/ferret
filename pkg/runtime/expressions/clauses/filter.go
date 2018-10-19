package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/datasource"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type FilterClause struct {
	src        core.SourceMap
	dataSource datasource.DataSource
	predicate  core.Expression
}

func NewFilterClause(
	src core.SourceMap,
	dataSource datasource.DataSource,
	predicate core.Expression,
) *FilterClause {
	return &FilterClause{
		src, dataSource,
		predicate,
	}
}

func (clause *FilterClause) Variables() datasource.Variables {
	return clause.dataSource.Variables()
}

func (clause *FilterClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	variables := clause.dataSource.Variables()

	return collections.NewFilterIterator(src, func(set collections.ResultSet) (bool, error) {
		innerScope := scope.Fork()

		err := variables.Apply(innerScope, set)

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
	})
}
