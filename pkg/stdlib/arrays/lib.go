package arrays

import "github.com/MontFerret/ferret/pkg/runtime/core"

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"APPEND":       Append,
		"FIRST":        First,
		"FLATTEN":      Flatten,
		"INTERSECTION": Intersection,
		"LAST":         Last,
		"MINUS":        Minus,
		"NTH":          Nth,
	}
}
