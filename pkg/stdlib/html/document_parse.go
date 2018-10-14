package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/html"
	"github.com/MontFerret/ferret/pkg/html/static"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DocumentParse parses a given HTML string and returns a HTML document.
// Returned HTML document is always static.
// @param html (String) - Target HTML string.
// @returns (HTMLDocument) - Parsed HTML static document.
func DocumentParse(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.StringType)

	if err != nil {
		return values.None, err
	}

	drv, err := html.FromContext(ctx, html.Static)

	if err != nil {
		return values.None, err
	}

	str := args[0].(values.String)

	return drv.(*static.Driver).ParseDocument(ctx, str)
}
