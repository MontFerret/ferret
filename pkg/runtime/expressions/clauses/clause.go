package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type baseClause struct {
	src        core.SourceMap
	dataSource collections.IterableExpression
}

func (clause *baseClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	return src, nil
}

func (clause *baseClause) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	iterator, err := clause.Iterate(ctx, scope)

	if err != nil {
		return values.None, err
	}

	return collections.ToArray(iterator)
}
