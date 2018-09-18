package collections

import "github.com/MontFerret/ferret/pkg/runtime/values"

type Collection interface {
	Length() values.Int
}
