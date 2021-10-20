package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// Measurable represents an interface of a value that can has length.
	Measurable interface {
		Length() values.Int
	}

	IndexedCollection interface {
		core.Value
		Measurable
		Get(idx values.Int) core.Value
		Set(idx values.Int, value core.Value) error
	}

	KeyedCollection interface {
		core.Value
		Measurable
		Keys() []values.String
		Get(key values.String) (core.Value, values.Boolean)
		Set(key values.String, value core.Value)
	}
)
