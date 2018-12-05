package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/html"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"time"
)

type LoadPageParams struct {
	Dynamic   values.Boolean
	Timeout   time.Duration
	UserAgent values.String
}

// Page loads a HTML document by a given url.
// By default, loads a document by http call - resulted document does not support any interactions.
// If passed "true" as a second argument, headless browser is used for loading the document which support interactions.
// @param url (String) - Target url string. If passed "about:blank" for dynamic document - it will open an empty page.
// @param isDynamicOrParams (Boolean|LoadPageParams) - Either a boolean value that indicates whether to use dynamic page
// or an object with the following properties :
// 		dynamic (Boolean) - Optional, indicates whether to use dynamic page.
// 		userAgent (String) - Optional, custom user agent.
// 		timeout (Int) - Optional, Page load timeout.
// @returns (HTMLDocument) - Returns loaded HTML document.
func Page(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.StringType)

	if err != nil {
		return values.None, err
	}

	url := args[0].(values.String)

	var params LoadPageParams

	if len(args) == 1 {
		params = newDefaultLoadPageParams()
	} else {
		p, err := newLoadPageParams(args[1])

		if err != nil {
			return values.None, err
		}

		params = p
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

func newDefaultLoadPageParams() LoadPageParams {
	return LoadPageParams{
		Dynamic:   false,
		UserAgent: "",
		Timeout:   time.Second * 30,
	}
}

func newLoadPageParams(arg core.Value) (LoadPageParams, error) {
	res := newDefaultLoadPageParams()

	if err := core.ValidateType(arg, core.BooleanType, core.ObjectType); err != nil {
		return res, err
	}

	if arg.Type() == core.BooleanType {
		res.Dynamic = arg.(values.Boolean)

		return res, nil
	}

	obj := arg.(*values.Object)

	isDynamic, exists := obj.Get(values.NewString("dynamic"))

	if exists {
		if err := core.ValidateType(isDynamic, core.BooleanType); err != nil {
			return res, err
		}

		res.Dynamic = isDynamic.(values.Boolean)
	}

	userAgent, exists := obj.Get(values.NewString("userAgent"))

	if exists {
		if err := core.ValidateType(userAgent, core.StringType); err != nil {
			return res, err
		}

		res.UserAgent = userAgent.(values.String)
	}

	timeout, exists := obj.Get(values.NewString("timeout"))

	if exists {
		if err := core.ValidateType(timeout, core.IntType); err != nil {
			return res, err
		}

		res.Timeout = time.Duration(timeout.(values.Int)) + time.Millisecond
	}

	return res, nil
}
