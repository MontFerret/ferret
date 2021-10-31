package drivers

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	HTTPRequestType   = core.NewType("ferret.drivers.HTTPRequest")
	HTTPResponseType  = core.NewType("ferret.drivers.HTTPResponse")
	HTTPHeaderType    = core.NewType("ferret.drivers.HTTPHeaders")
	HTTPCookieType    = core.NewType("ferret.drivers.HTTPCookie")
	HTTPCookiesType   = core.NewType("ferret.drivers.HTTPCookies")
	HTMLElementType   = core.NewType("ferret.drivers.HTMLElement")
	HTMLDocumentType  = core.NewType("ferret.drivers.HTMLDocument")
	HTMLPageType      = core.NewType("ferret.drivers.HTMLPageType")
	QuerySelectorType = core.NewType("ferret.drivers.QuerySelector")
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
