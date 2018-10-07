package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/static"
)

/*
 * Parses a given HTML string and returns a HTML document.
 * Returned HTML document is always static.
 * @param html (String) - Target HTML string.
 * @returns (HTMLDocument) - Parsed HTML static document.
 */
func DocumentParse(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.StringType)

	if err != nil {
		return values.None, err
	}

	drv, err := driver.FromContext(ctx, driver.Static)

	if err != nil {
		return values.None, err
	}

	str := args[0].(values.String)

	return drv.(*static.Driver).ParseDocument(ctx, str)
}
