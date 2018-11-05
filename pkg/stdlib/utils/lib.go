package utils

import "github.com/MontFerret/ferret/pkg/runtime/core"

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"WAIT":  Wait,
		"PRINT": Print,
	}
}
