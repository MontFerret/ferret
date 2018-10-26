package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Collection interface {
		Length() values.Int
	}

	IndexedCollection interface {
		Collection
		Get(idx values.Int) core.Value
		Set(idx values.Int, value core.Value) error
	}

	KeyedCollection interface {
		Collection
		Keys() []string
		Get(key values.String) (core.Value, values.Boolean)
		Set(key values.String, value core.Value)
	}
)
