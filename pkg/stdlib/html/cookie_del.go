package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// COOKIE_DEL gets a cookie from a given page by name.
// @param page (HTMLPage) - Target page.
// @param cookie (...HTTPCookie|String) - Cookie or cookie name to delete.
func CookieDel(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	inputs := args[1:]
	var currentCookies drivers.HTTPCookies
	cookies := make(drivers.HTTPCookies)

	for _, c := range inputs {
		switch cookie := c.(type) {
		case values.String:
			if currentCookies == nil {
				current, err := page.GetCookies(ctx)

				if err != nil {
					return values.None, err
				}

				currentCookies = current
			}

			found, isFound := currentCookies[cookie.String()]

			if isFound {
				cookies[cookie.String()] = found
			}

		case drivers.HTTPCookie:
			cookies[cookie.Name] = cookie
		default:
			return values.None, core.TypeError(c.Type(), types.String, drivers.HTTPCookieType)
		}
	}

	return values.None, page.DeleteCookies(ctx, cookies)
}
