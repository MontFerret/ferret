package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type DataSource struct {
	src         core.SourceMap
	valVariable string
	keyVariable string
	exp         core.Expression
}

func NewDataSource(
	src core.SourceMap,
	valVariable,
	keyVariable string,
	exp core.Expression,
) (collections.Iterable, error) {
	if exp == nil {
		return nil, core.Error(core.ErrMissedArgument, "expression")
	}

	return &DataSource{
		src,
		valVariable,
		keyVariable,
		exp,
	}, nil
}

func (ds *DataSource) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	select {
	case <-ctx.Done():
		return nil, core.ErrTerminated
	default:
		data, err := ds.exp.Exec(ctx, scope)

		if err != nil {
			return nil, core.SourceError(ds.src, err)
		}

		switch data.Type() {
		case types.Array:
			return collections.NewIndexedIterator(ds.valVariable, ds.keyVariable, data.(collections.IndexedCollection))
		case types.Object:
			return collections.NewKeyedIterator(ds.valVariable, ds.keyVariable, data.(collections.KeyedCollection))
		default:
			// fallback to user defined types
			switch collection := data.(type) {
			case core.Iterable:
				iterator, err := collection.Iterate(ctx)

				if err != nil {
					return nil, err
				}

				return collections.NewCoreIterator(ds.valVariable, ds.keyVariable, iterator)
			case collections.KeyedCollection:
				return collections.NewKeyedIterator(ds.valVariable, ds.keyVariable, collection)
			case collections.IndexedCollection:
				return collections.NewIndexedIterator(ds.valVariable, ds.keyVariable, collection)
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
