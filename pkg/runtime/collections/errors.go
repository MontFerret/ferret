package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

var (
	ErrExhausted = core.Error(core.ErrInvalidOperation, "iterator has been exhausted")
)
