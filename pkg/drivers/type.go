package drivers

import "github.com/MontFerret/ferret/pkg/runtime/core"

func createType(name string) core.Type {
	return core.NewType("ferret.drivers", name)
}

var (
	HTTPRequestType   = createType("HTTPRequest")
	HTTPResponseType  = createType("HTTPResponse")
	HTTPHeaderType    = createType("HTTPHeaders")
	HTTPCookieType    = createType("HTTPCookie")
	HTTPCookiesType   = createType("HTTPCookies")
	HTMLElementType   = createType("HTMLElement")
	HTMLDocumentType  = createType("HTMLDocument")
	HTMLPageType      = createType("HTMLPageType")
	QuerySelectorType = createType("QuerySelector")
)

// Comparison table of builtin types
var typeComparisonTable = map[core.Type]uint64{
	QuerySelectorType: 0,
	HTTPHeaderType:    1,
	HTTPCookieType:    2,
	HTTPCookiesType:   3,
	HTTPRequestType:   4,
	HTTPResponseType:  5,
	HTMLElementType:   6,
	HTMLDocumentType:  7,
	HTMLPageType:      8,
}

func CompareTypes(first, second core.Type) int64 {
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
