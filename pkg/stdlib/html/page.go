package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/html"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"time"
)

type LoadPageArgs struct {
	Dynamic bool
	Timeout time.Duration
}

// Page loads a HTML document by a given url.
// By default, loads a document by http call - resulted document does not support any interactions.
// If passed "true" as a second argument, headless browser is used for loading the document which support interactions.
// @param url (String) - Target url string. If passed "about:blank" for dynamic document - it will open an empty page.
// @param dynamicOrTimeout (Boolean|Int, optional) - If boolean value is passed, it indicates whether to use dynamic document.
// If integer values is passed it sets a custom timeout.
// @param timeout (Int, optional) - Sets a custom timeout.
// @returns (HTMLDocument) - Returns loaded HTML document.
func Page(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 3)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.StringType)

	if err != nil {
		return values.None, err
	}

	url := args[0].(values.String)

	params, err := parseLoadPageArgs(args)

	if err != nil {
		return values.None, err
	}

	var drv html.Driver

	ctx, cancel := context.WithTimeout(ctx, params.Timeout)
	defer cancel()

	if params.Dynamic == false {
		drv, err = html.FromContext(ctx, html.Static)
	} else {
		drv, err = html.FromContext(ctx, html.Dynamic)
	}

	if err != nil {
		return values.None, err
	}

	return drv.GetDocument(ctx, url)
}

func parseLoadPageArgs(args []core.Value) (LoadPageArgs, error) {
	res := LoadPageArgs{
		Timeout: time.Second * 30,
	}

	if len(args) == 3 {
		err := core.ValidateType(args[1], core.BooleanType)

		if err != nil {
			return res, err
		}

		res.Dynamic = bool(args[1].(values.Boolean))

		err = core.ValidateType(args[2], core.IntType)

		if err != nil {
			return res, err
		}

		res.Timeout = time.Duration(args[2].(values.Int)) * time.Millisecond
	} else if len(args) == 2 {
		err := core.ValidateType(args[1], core.BooleanType, core.IntType)

		if err != nil {
			return res, err
		}

		if args[1].Type() == core.BooleanType {
			res.Dynamic = bool(args[1].(values.Boolean))
		} else {
			res.Timeout = time.Duration(args[1].(values.Int)) * time.Millisecond
		}
	}

	return res, nil
}
