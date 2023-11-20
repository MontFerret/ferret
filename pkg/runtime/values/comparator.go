package values

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func Compare(a, b core.Value) int64 {
	aComparable, ok := a.(core.Comparable)

	if ok {
		return aComparable.Compare(b)
	}

	bComparable, ok := b.(core.Comparable)

	if ok {
		return -bComparable.Compare(a)
	}

	return types.Compare(core.Reflect(a), core.Reflect(b))
}

