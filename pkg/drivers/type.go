package drivers

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	HTTPResponseType  = core.NewType("HTTPResponse")
	HTTPHeaderType    = core.NewType("HTTPHeaders")
	HTTPCookieType    = core.NewType("HTTPCookie")
	HTTPCookiesType   = core.NewType("HTTPCookies")
	HTMLElementType   = core.NewType("HTMLElement")
	HTMLDocumentType  = core.NewType("HTMLDocument")
	HTMLPageType      = core.NewType("HTMLPageType")
	QuerySelectorType = core.NewType("QuerySelector")
)

// Comparison table of builtin types
var typeComparisonTable = map[core.Type]uint64{
	QuerySelectorType: 0,
	HTTPHeaderType:    1,
	HTTPCookieType:    2,
	HTTPCookiesType:   3,
	HTMLElementType:   4,
	HTMLDocumentType:  5,
	HTMLPageType:      6,
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
