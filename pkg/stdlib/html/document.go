package html

import (
	"context"
	"strings"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type PageLoadParams struct {
	drivers.Params
	Driver  string
	Timeout time.Duration
}

// DOCUMENT opens an HTML page by a given url.
// By default, loads a page by http call - resulted page does not support any interactions.
// @param params (Object) - Optional, An object containing the following properties :
// 		driver (String) - Optional, driver name.
//      timeout (Int) - Optional, timeout.
//      userAgent (String) - Optional, user agent.
//      keepCookies (Boolean) - Optional, boolean value indicating whether to use cookies from previous sessions.
//      	i.e. not to open a page in the Incognito mode.
//      cookies (HTTPCookies) - Optional, set of HTTP cookies.
//      headers (HTTPHeaders) - Optional, HTTP headers.
//      viewport (Viewport) - Optional, viewport params.
// @returns (HTMLPage) - Returns loaded HTML page.
func Open(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return values.None, err
	}

	url := args[0].(values.String)

	var params PageLoadParams

	if len(args) == 1 {
		params = newDefaultDocLoadParams(url)
	} else {
		p, err := newPageLoadParams(url, args[1])

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

	return drv.Open(ctx, params.Params)
}

func newDefaultDocLoadParams(url values.String) PageLoadParams {
	return PageLoadParams{
		Params: drivers.Params{
			URL: url.String(),
		},
		Timeout: drivers.DefaultPageLoadTimeout * time.Millisecond,
	}
}

func newPageLoadParams(url values.String, arg core.Value) (PageLoadParams, error) {
	res := newDefaultDocLoadParams(url)

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

			res.Timeout = time.Duration(timeout.(values.Int)) * time.Millisecond
		}

		userAgent, exists := obj.Get(values.NewString("userAgent"))

		if exists {
			if err := core.ValidateType(userAgent, types.String); err != nil {
				return res, err
			}

			res.UserAgent = userAgent.String()
		}

		keepCookies, exists := obj.Get(values.NewString("keepCookies"))

		if exists {
			if err := core.ValidateType(keepCookies, types.Boolean); err != nil {
				return res, err
			}

			res.KeepCookies = bool(keepCookies.(values.Boolean))
		}

		cookies, exists := obj.Get(values.NewString("cookies"))

		if exists {
			if err := core.ValidateType(cookies, types.Array, types.Object); err != nil {
				return res, err
			}

			switch c := cookies.(type) {
			case *values.Array:
				cookies, err := parseCookieArray(c)

				if err != nil {
					return res, err
				}

				res.Cookies = cookies
			case *values.Object:
				cookies, err := parseCookieObject(c)

				if err != nil {
					return res, err
				}

				res.Cookies = cookies
			default:
				res.Cookies = make(drivers.HTTPCookies)
			}
		}

		headers, exists := obj.Get(values.NewString("headers"))

		if exists {
			if err := core.ValidateType(headers, types.Object); err != nil {
				return res, err
			}

			header := parseHeader(headers.(*values.Object))
			res.Headers = header
		}

		viewport, exists := obj.Get(values.NewString("viewport"))

		if exists {
			viewport, err := parseViewport(viewport)

			if err != nil {
				return res, err
			}

			res.Viewport = viewport
		}
	case types.String:
		res.Driver = arg.(values.String).String()
	case types.Boolean:
		b := arg.(values.Boolean)

		// fallback
		if b {
			res.Driver = cdp.DriverName
		}
	}

	return res, nil
}

func parseCookieObject(obj *values.Object) (drivers.HTTPCookies, error) {
	var err error
	res := make(drivers.HTTPCookies)

	obj.ForEach(func(value core.Value, _ string) bool {
		cookie, e := parseCookie(value)

		if e != nil {
			err = e

			return false
		}

		res[cookie.Name] = cookie

		return true
	})

	return res, err
}

func parseCookieArray(arr *values.Array) (drivers.HTTPCookies, error) {
	var err error
	res := make(drivers.HTTPCookies)

	arr.ForEach(func(value core.Value, _ int) bool {
		cookie, e := parseCookie(value)

		if e != nil {
			err = e

			return false
		}

		res[cookie.Name] = cookie

		return true
	})

	return res, err
}

func parseCookie(value core.Value) (drivers.HTTPCookie, error) {
	err := core.ValidateType(value, types.Object, drivers.HTTPCookieType)

	if err != nil {
		return drivers.HTTPCookie{}, err
	}

	if value.Type() == drivers.HTTPCookieType {
		return value.(drivers.HTTPCookie), nil
	}

	co := value.(*values.Object)

	cookie := drivers.HTTPCookie{
		Name:   co.MustGet("name").String(),
		Value:  co.MustGet("value").String(),
		Path:   co.MustGet("path").String(),
		Domain: co.MustGet("domain").String(),
	}

	maxAge, exists := co.Get("maxAge")

	if exists {
		if err = core.ValidateType(maxAge, types.Int); err != nil {
			return drivers.HTTPCookie{}, err
		}

		cookie.MaxAge = int(maxAge.(values.Int))
	}

	expires, exists := co.Get("expires")

	if exists {
		if err = core.ValidateType(maxAge, types.DateTime, types.String); err != nil {
			return drivers.HTTPCookie{}, err
		}

		if expires.Type() == types.DateTime {
			cookie.Expires = expires.(values.DateTime).Unwrap().(time.Time)
		} else {
			t, err := time.Parse(expires.String(), values.DefaultTimeLayout)

			if err != nil {
				return drivers.HTTPCookie{}, err
			}

			cookie.Expires = t
		}
	}

	sameSite, exists := co.Get("sameSite")

	if exists {
		sameSite := strings.ToLower(sameSite.String())

		switch sameSite {
		case "lax":
			cookie.SameSite = drivers.SameSiteLaxMode
		case "strict":
			cookie.SameSite = drivers.SameSiteStrictMode
		default:
			cookie.SameSite = drivers.SameSiteDefaultMode
		}
	}

	httpOnly, exists := co.Get("httpOnly")

	if exists {
		if err = core.ValidateType(httpOnly, types.Boolean); err != nil {
			return drivers.HTTPCookie{}, err
		}

		cookie.HTTPOnly = bool(httpOnly.(values.Boolean))
	}

	secure, exists := co.Get("secure")

	if exists {
		if err = core.ValidateType(secure, types.Boolean); err != nil {
			return drivers.HTTPCookie{}, err
		}

		cookie.Secure = bool(secure.(values.Boolean))
	}

	return cookie, err
}

func parseHeader(headers *values.Object) drivers.HTTPHeaders {
	res := make(drivers.HTTPHeaders)

	headers.ForEach(func(value core.Value, key string) bool {
		res.Set(key, value.String())

		return true
	})

	return res
}

func parseViewport(value core.Value) (*drivers.Viewport, error) {
	if err := core.ValidateType(value, types.Object); err != nil {
		return nil, err
	}

	res := &drivers.Viewport{}

	viewport := value.(*values.Object)

	width, exists := viewport.Get(values.NewString("width"))

	if exists {
		if err := core.ValidateType(width, types.Int); err != nil {
			return nil, err
		}

		res.Width = int(values.ToInt(width))
	}

	height, exists := viewport.Get(values.NewString("height"))

	if exists {
		if err := core.ValidateType(height, types.Int); err != nil {
			return nil, err
		}

		res.Height = int(values.ToInt(height))
	}

	mobile, exists := viewport.Get(values.NewString("mobile"))

	if exists {
		res.Mobile = bool(values.ToBoolean(mobile))
	}

	landscape, exists := viewport.Get(values.NewString("landscape"))

	if exists {
		res.Landscape = bool(values.ToBoolean(landscape))
	}

	scaleFactor, exists := viewport.Get(values.NewString("scaleFactor"))

	if exists {
		res.ScaleFactor = float64(values.ToFloat(scaleFactor))
	}

	return res, nil
}
