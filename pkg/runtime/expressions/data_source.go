package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type DataSource struct {
	src       core.SourceMap
	variables collections.Variables
	exp       core.Expression
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
		collections.Variables{valVariable, keyVariable},
		exp,
	}, nil
}

func (ds *DataSource) Variables() collections.Variables {
	return ds.variables
}

func (ds *DataSource) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	data, err := ds.exp.Exec(ctx, scope)

	if err != nil {
		return nil, core.SourceError(ds.src, err)
	}

	valVar := ds.variables[0]
	keyVar := ds.variables[1]

	switch data.Type() {
	case core.ArrayType:
		return collections.NewIndexedIterator(valVar, keyVar, data.(collections.IndexedCollection))
	case core.ObjectType:
		return collections.NewKeyedIterator(valVar, keyVar, data.(collections.KeyedCollection))
	case core.HTMLElementType, core.HTMLDocumentType:
		return collections.NewHTMLNodeIterator(valVar, keyVar, data.(values.HTMLNode))
	default:
		// fallback to user defined types
		switch data.(type) {
		case collections.KeyedCollection:
			return collections.NewIndexedIterator(valVar, keyVar, data.(collections.IndexedCollection))
		case collections.IndexedCollection:
			return collections.NewKeyedIterator(valVar, keyVar, data.(collections.KeyedCollection))
		default:
			return nil, core.TypeError(
				data.Type(),
				core.ArrayType,
				core.ObjectType,
				core.HTMLDocumentType,
				core.HTMLElementType,
			)
		}
	}
}
