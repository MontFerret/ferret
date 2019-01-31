package html

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type DocumentLoadParams struct {
	Dynamic values.Boolean
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

	if params.Dynamic {
		drv, err := drivers.DynamicFrom(ctx)

		if err != nil {
			return values.None, err
		}

		return drv.GetDocument(ctx, url)
	}

	drv, err := drivers.StaticFrom(ctx)

	if err != nil {
		return values.None, err
	}

	return drv.GetDocument(ctx, url)
}

func newDefaultDocLoadParams() DocumentLoadParams {
	return DocumentLoadParams{
		Dynamic: false,
		Timeout: time.Second * 30,
	}
}

func newDocLoadParams(arg core.Value) (DocumentLoadParams, error) {
	res := newDefaultDocLoadParams()

	if err := core.ValidateType(arg, types.Boolean, types.Object); err != nil {
		return res, err
	}

	if arg.Type() == types.Boolean {
		res.Dynamic = arg.(values.Boolean)

		return res, nil
	}

	obj := arg.(*values.Object)

	isDynamic, exists := obj.Get(values.NewString("dynamic"))

	if exists {
		if err := core.ValidateType(isDynamic, types.Boolean); err != nil {
			return res, err
		}

		res.Dynamic = isDynamic.(values.Boolean)
	}

	timeout, exists := obj.Get(values.NewString("timeout"))

	if exists {
		if err := core.ValidateType(timeout, types.Int); err != nil {
			return res, err
		}

		res.Timeout = time.Duration(timeout.(values.Int)) + time.Millisecond
	}

	return res, nil
}
