package values

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
)

func NewArray(cap int) core.List {
	return internal.NewArray(cap)
}
