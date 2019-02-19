package drivers

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	HTMLElementType  = core.NewType("HTMLElement")
	HTMLDocumentType = core.NewType("HTMLDocument")
)

// Comparison table of builtin types
var typeComparisonTable = map[core.Type]uint64{
	HTMLElementType:  0,
	HTMLDocumentType: 1,
}

func Compare(first, second core.Type) int64 {
	f, ok := typeComparisonTable[first]

	if !ok {
		return -1
	}

	s, ok := typeComparisonTable[second]

	if !ok {
		return -1
	}

	if f == s {
		return 0
	}

	if f > s {
		return 1
	}

	return -1
}
