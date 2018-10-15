package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/html"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Document loads a document by a given url.
// By default, loads a document by http call - resulted document does not support any interactions.
// If passed "true" as a second argument, headless browser is used for loading the document which support interactions.
// @param url (String) - Target url string. If passed "about:blank" for dynamic document - it will open an empty page.
// @param dynamic (Boolean) - Optional boolean value indicating whether to use dynamic document.
// @returns (HTMLDocument) - Returns loaded HTML document.
func Document(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.StringType)

	url := args[0].(values.String)
	dynamic := values.False

	if len(args) == 2 {
		err = core.ValidateType(args[1], core.BooleanType)

		if err != nil {
			return values.None, err
		}

		dynamic = args[1].(values.Boolean)
	}

	var drv html.Driver

	if !dynamic {
		drv, err = html.FromContext(ctx, html.Static)
	} else {
		drv, err = html.FromContext(ctx, html.Dynamic)
	}

	if err != nil {
		return values.None, err
	}

	return drv.GetDocument(ctx, url)
}
