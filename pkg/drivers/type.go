package drivers

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	HTTPResponseType = core.NewType("HTTPResponse")
	HTTPHeaderType   = core.NewType("HTTPHeaders")
	HTTPCookieType   = core.NewType("HTTPCookie")
	HTTPCookiesType  = core.NewType("HTTPCookies")
	HTMLElementType  = core.NewType("HTMLElement")
	HTMLDocumentType = core.NewType("HTMLDocument")
	HTMLPageType     = core.NewType("HTMLPageType")
)

// Comparison table of builtin types
var typeComparisonTable = map[core.Type]uint64{
	HTTPHeaderType:   0,
	HTTPCookieType:   1,
	HTTPCookiesType:  2,
	HTMLElementType:  3,
	HTMLDocumentType: 4,
	HTMLPageType:     5,
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
