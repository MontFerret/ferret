package math

import "github.com/MontFerret/ferret/pkg/runtime/core"

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"ABS":  Abs,
		"ACOS": Acos,
		"ASIN": Asin,
	}
}
