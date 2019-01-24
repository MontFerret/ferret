package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
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
		case core.ArrayType:
			return collections.NewIndexedIterator(ds.valVariable, ds.keyVariable, data.(collections.IndexedCollection))
		case core.ObjectType:
			return collections.NewKeyedIterator(ds.valVariable, ds.keyVariable, data.(collections.KeyedCollection))
		default:
			// fallback to user defined types
			switch data.(type) {
			case collections.IterableCollection:
				collection := data.(collections.IterableCollection)
				iterator, err := collection.Iterate(ctx)

				if err != nil {
					return nil, err
				}

				return collections.NewCollectionIterator(ds.valVariable, ds.keyVariable, iterator)
			case collections.KeyedCollection:
				return collections.NewIndexedIterator(ds.valVariable, ds.keyVariable, data.(collections.IndexedCollection))
			case collections.IndexedCollection:
				return collections.NewKeyedIterator(ds.valVariable, ds.keyVariable, data.(collections.KeyedCollection))
			default:
				return nil, core.TypeError(
					data.Type(),
					core.ArrayType,
					core.ObjectType,
					core.CustomType,
				)
			}
		}
	}
}
