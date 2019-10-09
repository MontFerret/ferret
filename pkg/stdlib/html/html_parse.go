package html

import (
	"context"
	"io/ioutil"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type HtmlParseParams struct {
	FileName string
	Driver   string
}

// HTML_PARSE loads an HTML page from a local file
// By default, loads a page by http call - resulted page does not support any interactions.
// @param params (Object) - Optional, An object containing the following properties :
// 		driver (String) - Optional, driver name.
// @returns (HTMLPage) - Returns loaded HTML page.
func HTMLParse(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return values.None, err
	}

	var params PageLoadParams

	if len(args) > 1 {
		p, err := newPageLoadParams("about:blank", args[1])

		if err != nil {
			return values.None, err
		}

		params = p
	}

	b, err := ioutil.ReadFile(args[0].String())
	if err != nil {
		return values.None, err
	}

	drv, err := drivers.FromContext(ctx, params.Driver)
	if err != nil {
		return values.None, err
	}

	return drv.Parse(ctx, b)
}
