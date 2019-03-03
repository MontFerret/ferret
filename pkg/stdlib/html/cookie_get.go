package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// CookieSet gets a cookie from a given document by name.
// @param doc (HTMLDocument) - Target HTMLDocument.
// @param name (String) - Cookie or cookie name to delete.
func CookieGet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], drivers.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	doc := args[0].(drivers.HTMLDocument)
	name := args[1].(values.String)

	if err != nil {
		return values.None, err
	}

	cookies, err := doc.GetCookies(ctx)

	if err != nil {
		return values.None, err
	}

	found, _ := cookies.Find(func(value core.Value, _ int) bool {
		cookie, ok := value.(drivers.HTTPCookie)

		if !ok {
			return ok
		}

		return cookie.Name == name.String()
	})

	return found, nil
}
