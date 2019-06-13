package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// CookieSet gets a cookie from a given page by name.
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
	var currentCookies *values.Array
	cookies := make([]drivers.HTTPCookie, 0, len(inputs))

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

			found, isFound := currentCookies.Find(func(value core.Value, _ int) bool {
				cv := value.(drivers.HTTPCookie)

				return cv.Name == cookie.String()
			})

			if isFound {
				cookies = append(cookies, found.(drivers.HTTPCookie))
			}

		case drivers.HTTPCookie:
			cookies = append(cookies, cookie)

		default:
			return values.None, core.TypeError(c.Type(), types.String, drivers.HTTPCookieType)
		}
	}

	return values.None, page.DeleteCookies(ctx, cookies...)
}
