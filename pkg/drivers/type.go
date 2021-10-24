package drivers

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	HTTPRequestType   = core.NewType("drivers.HTTPRequest")
	HTTPResponseType  = core.NewType("drivers.HTTPResponse")
	HTTPHeaderType    = core.NewType("drivers.HTTPHeaders")
	HTTPCookieType    = core.NewType("drivers.HTTPCookie")
	HTTPCookiesType   = core.NewType("drivers.HTTPCookies")
	HTMLElementType   = core.NewType("drivers.HTMLElement")
	HTMLDocumentType  = core.NewType("drivers.HTMLDocument")
	HTMLPageType      = core.NewType("drivers.HTMLPageType")
	QuerySelectorType = core.NewType("drivers.QuerySelector")
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
