package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type FilterClause struct {
	*baseClause
	valVar    string
	keyVar    string
	predicate core.Expression
}

func NewFilterClause(
	src core.SourceMap,
	dataSource collections.IterableExpression,
	valVar string,
	keyVar string,
	predicate core.Expression,
) *FilterClause {
	return &FilterClause{
		&baseClause{src, dataSource},
		valVar,
		keyVar,
		predicate,
	}
}

func (clause *FilterClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewFilterIterator(src, func(val core.Value, key core.Value) (bool, error) {
		innerScope := scope.Fork()

		innerScope.SetVariable(clause.valVar, val)

		if clause.keyVar != "" {
			innerScope.SetVariable(clause.keyVar, key)
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
