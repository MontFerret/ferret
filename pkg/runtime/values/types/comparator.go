package types

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	// Comparison table of builtin types
	typeComparisonTable = map[int64]uint64{
		None.ID():     0,
		Boolean.ID():  1,
		Int.ID():      2,
		Float.ID():    3,
		String.ID():   4,
		DateTime.ID(): 5,
		Regexp.ID():   6,
		Range.ID():    7,
		Array.ID():    8,
		Object.ID():   9,
		Binary.ID():   10,
	}
)

func Compare(first, second core.Type) int64 {
	f, ok := typeComparisonTable[first.ID()]

	// custom type
	if !ok {
		return -1
	}

	s, ok := typeComparisonTable[second.ID()]

	// custom type
	if !ok {
		return 1
	}

	if f == s {
		return 0
	}

	if f > s {
		return 1
	}

	return -1
}

func Equal(first, second core.Type) bool {
	return Compare(first, second) == 0
}
