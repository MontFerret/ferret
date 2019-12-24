package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// COOKIE_GET gets a cookie from a given page by name.
// @param page (HTMLPage) - Target page.
// @param name (String) - Cookie or cookie name to delete.
func CookieGet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	name := args[1].(values.String)

	cookies, err := page.GetCookies(ctx)

	if err != nil {
		return values.None, err
	}

	cookie, found := cookies[name.String()]

	if found {
		return cookie, nil
	}

	return values.None, nil
}
