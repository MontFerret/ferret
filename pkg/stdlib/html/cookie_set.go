package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// COOKIE_SET sets cookies to a given page
// @param page (HTMLPage) - Target page.
// @param cookie... (HTTPCookie) - Target cookies.
func CookieSet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	cookies := make(drivers.HTTPCookies)

	for _, c := range args[1:] {
		cookie, err := parseCookie(c)

		if err != nil {
			return values.None, err
		}

		cookies[cookie.Name] = cookie
	}

	return values.None, page.SetCookies(ctx, cookies)
}
