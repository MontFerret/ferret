package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// COOKIE_SET sets cookies to a given page
// @param {HTMLPage} page - Target page.
// @param {HTTPCookie, repeated} cookies - Target cookies.
func CookieSet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return core.None, err
	}

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return core.None, err
	}

	cookies := drivers.NewHTTPCookies()

	for _, c := range args[1:] {
		cookie, err := parseCookie(c)

		if err != nil {
			return core.None, err
		}

		cookies.Set(cookie)
	}

	return core.None, page.SetCookies(ctx, cookies)
}
