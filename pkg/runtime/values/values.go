package values

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func NewArray(cap int) core.List {
	return core.NewArray(cap)
}
