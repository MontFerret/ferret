package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// CookieSet sets cookies to a given page
// @param page (HTMLPage) - Target page.
// @param cookie... (HTTPCookie) - Target cookies.
func CookieSet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], drivers.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	page, err := toPage(args[0])

	if err != nil {
		return values.None, err
	}

	cookies := make([]drivers.HTTPCookie, 0, len(args)-1)

	for _, c := range args[1:] {
		cookie, err := parseCookie(c)

		if err != nil {
			return values.None, err
		}

		cookies = append(cookies, cookie)
	}

	return values.None, page.SetCookies(ctx, cookies...)
}
