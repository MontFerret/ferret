package html

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type DocumentLoadParams struct {
	Driver  string
	Timeout time.Duration
}

// Document loads a HTML document by a given url.
// By default, loads a document by http call - resulted document does not support any interactions.
// If passed "true" as a second argument, headless browser is used for loading the document which support interactions.
// @param url (String) - Target url string. If passed "about:blank" for dynamic document - it will open an empty page.
// @param isDynamicOrParams (Boolean|DocumentLoadParams) - Either a boolean value that indicates whether to use dynamic page
// or an object with the following properties :
// 		dynamic (Boolean) - Optional, indicates whether to use dynamic page.
// 		timeout (Int) - Optional, Document load timeout.
// @returns (HTMLDocument) - Returns loaded HTML document.
func Document(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return values.None, err
	}

	url := args[0].(values.String)

	var params DocumentLoadParams

	if len(args) == 1 {
		params = newDefaultDocLoadParams()
	} else {
		p, err := newDocLoadParams(args[1])

		if err != nil {
			return values.None, err
		}

		params = p
	}

	ctx, cancel := context.WithTimeout(ctx, params.Timeout)
	defer cancel()

	drv, err := drivers.FromContext(ctx, params.Driver)

	if err != nil {
		return values.None, err
	}

	return drv.GetDocument(ctx, url)
}

func newDefaultDocLoadParams() DocumentLoadParams {
	return DocumentLoadParams{
		Timeout: time.Second * 30,
	}
}

func newDocLoadParams(arg core.Value) (DocumentLoadParams, error) {
	res := newDefaultDocLoadParams()

	if err := core.ValidateType(arg, types.Boolean, types.String, types.Object); err != nil {
		return res, err
	}

	switch arg.Type() {
	case types.Object:
		obj := arg.(*values.Object)

		driver, exists := obj.Get(values.NewString("driver"))

		if exists {
			if err := core.ValidateType(driver, types.String); err != nil {
				return res, err
			}

			res.Driver = driver.(values.String).String()
		}

		timeout, exists := obj.Get(values.NewString("timeout"))

		if exists {
			if err := core.ValidateType(timeout, types.Int); err != nil {
				return res, err
			}

			res.Timeout = time.Duration(timeout.(values.Int)) + time.Millisecond
		}

		break
	case types.String:
		res.Driver = arg.(values.String).String()

		break
	case types.Boolean:
		b := arg.(values.Boolean)

		// fallback
		if b {
			res.Driver = cdp.DriverName
		}

		break
	}

	return res, nil
}
