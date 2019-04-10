package types

import "github.com/MontFerret/ferret/pkg/runtime/core"

// Comparison table of builtin types
var typeComparisonTable = map[core.Type]uint64{
	None:     0,
	Boolean:  1,
	Int:      2,
	Float:    3,
	String:   4,
	DateTime: 5,
	Array:    6,
	Object:   7,
	Binary:   8,
}

func Compare(first, second core.Type) int64 {
	f, ok := typeComparisonTable[first]

	// custom type
	if !ok {
		return -1
	}

	s, ok := typeComparisonTable[second]

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
