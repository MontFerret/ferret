package clauses

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type LimitClause struct {
	src        core.SourceMap
	dataSource collections.Iterable
	count      core.Expression
	offset     core.Expression
}

func NewLimitClause(
	src core.SourceMap,
	dataSource collections.Iterable,
	count core.Expression,
	offset core.Expression,
) (collections.Iterable, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "dataSource source")
	}

	return &LimitClause{src, dataSource, count, offset}, nil
}

func (clause *LimitClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, core.SourceError(clause.src, err)
	}

	count, err := clause.count.Exec(ctx, scope)

	if err != nil {
		return nil, core.SourceError(clause.src, err)
	}

	offset, err := clause.offset.Exec(ctx, scope)

	if err != nil {
		return nil, core.SourceError(clause.src, err)
	}

	countInt, err := clause.parseValue(count)

	if err != nil {
		return nil, err
	}

	offsetInt, err := clause.parseValue(offset)

	if err != nil {
		return nil, err
	}

	iterator, err := collections.NewLimitIterator(src, countInt, offsetInt)

	if err != nil {
		return nil, core.SourceError(clause.src, err)
	}

	return iterator, nil
}

func (clause *LimitClause) parseValue(val core.Value) (int, error) {
	if val.Type() == types.Int {
		return val.Unwrap().(int), nil
	}

	if val.Type() == types.Float {
		return int(val.Unwrap().(float64)), nil
	}

	return -1, core.TypeError(val.Type(), types.Int, types.Float)
}
