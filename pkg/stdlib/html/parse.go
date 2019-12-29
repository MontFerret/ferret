package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ParseParams struct {
	drivers.ParseParams
	Driver string
}

// PARSE loads an HTML page from a given string or byte array
// @param params (Object) - Optional, an object containing the following properties :
// 	driver (String) - Optional, driver name.
//      keepCookies (Boolean) - Optional, boolean value indicating whether to use cookies from previous sessions.
//      	i.e. not to open a page in the Incognito mode.
//      cookies (HTTPCookies) - Optional, set of HTTP cookies.
//      headers (HTTPHeaders) - Optional, HTTP headers.
//      viewport (Viewport) - Optional, viewport params.
// @returns (HTMLPage) - Returns parsed and loaded HTML page.
func Parse(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 2); err != nil {
		return values.None, err
	}

	arg1 := args[0]

	if err := core.ValidateType(arg1, types.String, types.Binary); err != nil {
		return values.None, err
	}

	var content []byte

	if arg1.Type() == types.String {
		content = []byte(arg1.(values.String))
	} else {
		content = []byte(arg1.(values.Binary))
	}

	var params ParseParams

	if len(args) > 1 {
		if err := core.ValidateType(args[1], types.Object); err != nil {
			return values.None, err
		}

		p, err := parseParseParams(content, args[1].(*values.Object))

		if err != nil {
			return values.None, err
		}

		params = p
	} else {
		params = defaultParseParams(content)
	}

	drv, err := drivers.FromContext(ctx, params.Driver)

	if err != nil {
		return values.None, err
	}

	return drv.Parse(ctx, params.ParseParams)
}

func defaultParseParams(content []byte) ParseParams {
	return ParseParams{
		ParseParams: drivers.ParseParams{
			Content: content,
		},
		Driver: "",
	}
}

func parseParseParams(content []byte, arg *values.Object) (ParseParams, error) {
	res := defaultParseParams(content)

	if arg.Has("driver") {
		driverName := arg.MustGet("driver")

		if err := core.ValidateType(driverName, types.String); err != nil {
			return ParseParams{}, errors.Wrap(err, ".driver")
		}

		res.Driver = driverName.String()
	}

	if arg.Has("keepCookies") {
		keepCookies := arg.MustGet("keepCookies")

		if err := core.ValidateType(keepCookies, types.Boolean); err != nil {
			return ParseParams{}, errors.Wrap(err, ".keepCookies")
		}

		res.KeepCookies = bool(keepCookies.(values.Boolean))
	}

	if arg.Has("cookies") {
		cookies := arg.MustGet("cookies")

		if err := core.ValidateType(cookies, types.Array, types.Object); err != nil {
			return res, err
		}

		switch c := cookies.(type) {
		case *values.Array:
			cookies, err := parseCookieArray(c)

			if err != nil {
				return ParseParams{}, errors.Wrap(err, ".cookies")
			}

			res.Cookies = cookies
		case *values.Object:
			cookies, err := parseCookieObject(c)

			if err != nil {
				return ParseParams{}, errors.Wrap(err, ".cookies")
			}

			res.Cookies = cookies
		default:
			res.Cookies = make(drivers.HTTPCookies)
		}
	}

	if arg.Has("headers") {
		headers := arg.MustGet("headers")

		if err := core.ValidateType(headers, types.Object); err != nil {
			return ParseParams{}, errors.Wrap(err, ".headers")
		}

		res.Headers = parseHeader(headers.(*values.Object))
	}

	if arg.Has("viewport") {
		viewport, err := parseViewport(arg.MustGet("viewport"))

		if err != nil {
			return ParseParams{}, errors.Wrap(err, ".viewport")
		}

		res.Viewport = viewport
	}

	return res, nil
}
