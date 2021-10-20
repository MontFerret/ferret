package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type ForInIterableExpression struct {
	src         core.SourceMap
	valVariable string
	keyVariable string
	exp         core.Expression
}

func NewForInIterableExpression(
	src core.SourceMap,
	valVariable,
	keyVariable string,
	exp core.Expression,
) (collections.Iterable, error) {
	if exp == nil {
		return nil, core.Error(core.ErrMissedArgument, "expression")
	}

	return &ForInIterableExpression{
		src,
		valVariable,
		keyVariable,
		exp,
	}, nil
}

func (iterable *ForInIterableExpression) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	select {
	case <-ctx.Done():
		return nil, core.ErrTerminated
	default:
		data, err := iterable.exp.Exec(ctx, scope)

		if err != nil {
			return nil, core.SourceError(iterable.src, err)
		}

		switch data.Type() {
		case types.Array:
			return collections.NewIndexedIterator(iterable.valVariable, iterable.keyVariable, data.(collections.IndexedCollection))
		case types.Object:
			return collections.NewKeyedIterator(iterable.valVariable, iterable.keyVariable, data.(collections.KeyedCollection))
		default:
			// fallback to user defined types
			switch collection := data.(type) {
			case core.Iterable:
				iterator, err := collection.Iterate(ctx)

				if err != nil {
					return nil, err
				}

				return collections.FromCoreIterator(iterable.valVariable, iterable.keyVariable, iterator)
			case collections.KeyedCollection:
				return collections.NewKeyedIterator(iterable.valVariable, iterable.keyVariable, collection)
			case collections.IndexedCollection:
				return collections.NewIndexedIterator(iterable.valVariable, iterable.keyVariable, collection)
			default:
				return nil, core.TypeError(
					data.Type(),
					types.Array,
					types.Object,
					core.NewType("Iterable"),
				)
			}
		}
	}
}
