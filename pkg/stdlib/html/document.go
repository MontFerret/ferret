package html

import (
	"context"
	"github.com/pkg/errors"
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
// @param {Object} [params] - An object containing the following properties :
// @param {String} [params.driver] - Driver name to use.
// @param {Int} [params.timeout=60000] - Page load timeout.
// @param {String} [params.userAgent] - Custom user agent.
// @param {Boolean} [params.keepCookies=False] - Boolean value indicating whether to use cookies from previous sessions i.e. not to open a page in the Incognito mode.
// @param {Object[] | Object} [params.cookies] - Set of HTTP cookies to use during page loading.
// @param {String} params.cookies.*.name - Cookie name.
// @param {String} params.cookies.*.value - Cookie value.
// @param {String} params.cookies.*.path - Cookie path.
// @param {String} params.cookies.*.domain - Cookie domain.
// @param {Int} [params.cookies.*.maxAge] - Cookie max age.
// @param {String|DateTime} [params.cookies.*.expires] - Cookie expiration date time.
// @param {String} [params.cookies.*.sameSite] - Cookie cross-origin policy.
// @param {Boolean} [params.cookies.*.httpOnly=false] - Cookie cannot be accessed through client side script.
// @param {Boolean} [params.cookies.*.secure=false] - Cookie sent to the server only with an encrypted request over the HTTPS protocol.
// @param {Object} [params.headers] - Set of HTTP headers to use during page loading.
// @param {Object} [params.ignore] - Set of parameters to ignore some page functionality or behavior.
// @param {Object[]} [params.ignore.resources] - Collection of rules to ignore resources during page load and navigation.
// @param {String} [params.ignore.resources.*.url] - Resource url pattern. If set, requests for matching urls will be blocked. Wildcards ('*' -> zero or more, '?' -> exactly one) are allowed. Escape character is backslash. Omitting is equivalent to "*".
// @param {String} [params.ignore.resources.*.type] - Resource type. If set, requests for matching resource types will be blocked.
// @param {Object[]} [params.ignore.statusCodes] - Collection of rules to ignore certain HTTP codes that can cause failures.
// @param {String} [params.ignore.statusCodes.*.url] - Url pattern. If set, codes for matching urls will be ignored. Wildcards ('*' -> zero or more, '?' -> exactly one) are allowed. Escape character is backslash. Omitting is equivalent to "*".
// @param {Int} [params.ignore.statusCodes.*.code] - HTTP code to ignore.
// @param {Object} [params.viewport] - Viewport params.
// @param {Int} [params.viewport.height] - Viewport height.
// @param {Int} [params.viewport.width] - Viewport width.
// @param {Float} [params.viewport.scaleFactor] - Viewport scale factor.
// @param {Boolean} [params.viewport.mobile] - Value that indicates whether to emulate mobile device.
// @param {Boolean} [params.viewport.landscape] - Value that indicates whether to render a page in landscape position.
// @param {String} [params.charset] - (only HTTPDriver) Source charset content to convert UTF-8.
// @return {HTMLPage} - Loaded HTML page.
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
				res.Cookies = drivers.NewHTTPCookies()
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

		ignore, exists := obj.Get(values.NewString("ignore"))

		if exists {
			ignore, err := parseIgnore(ignore)

			if err != nil {
				return res, err
			}

			res.Ignore = ignore
		}

		charset, exists := obj.Get(values.NewString("charset"))

		if exists {
			if err := core.ValidateType(charset, types.String); err != nil {
				return res, err
			}

			res.Charset = charset.String()
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

func parseCookieObject(obj *values.Object) (*drivers.HTTPCookies, error) {
	if obj == nil {
		return nil, errors.Wrap(core.ErrMissedArgument, "cookies")
	}

	var err error
	res := drivers.NewHTTPCookies()

	obj.ForEach(func(value core.Value, _ string) bool {
		cookie, e := parseCookie(value)

		if e != nil {
			err = e

			return false
		}

		res.Set(cookie)

		return true
	})

	return res, err
}

func parseCookieArray(arr *values.Array) (*drivers.HTTPCookies, error) {
	if arr == nil {
		return nil, errors.Wrap(core.ErrMissedArgument, "cookies")
	}

	var err error
	res := drivers.NewHTTPCookies()

	arr.ForEach(func(value core.Value, _ int) bool {
		cookie, e := parseCookie(value)

		if e != nil {
			err = e

			return false
		}

		res.Set(cookie)

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
		if err = core.ValidateType(expires, types.DateTime, types.String); err != nil {
			return drivers.HTTPCookie{}, err
		}

		if expires.Type() == types.DateTime {
			cookie.Expires = expires.(values.DateTime).Unwrap().(time.Time)
		} else {
			t, err := time.Parse(values.DefaultTimeLayout, expires.String())

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

func parseHeader(headers *values.Object) *drivers.HTTPHeaders {
	res := drivers.NewHTTPHeaders()

	headers.ForEach(func(value core.Value, key string) bool {
		if value.Type() == types.Array {
			value := value.(*values.Array)

			keyValues := make([]string, 0, value.Length())

			value.ForEach(func(v core.Value, idx int) bool {
				keyValues = append(keyValues, v.String())

				return true
			})

			res.SetArr(key, keyValues)
		} else {
			res.Set(key, value.String())
		}

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

func parseIgnore(value core.Value) (*drivers.Ignore, error) {
	if err := core.ValidateType(value, types.Object); err != nil {
		return nil, err
	}

	res := &drivers.Ignore{}

	ignore := value.(*values.Object)

	resources, exists := ignore.Get("resources")

	if exists {
		if err := core.ValidateType(resources, types.Array); err != nil {
			return nil, err
		}

		resources := resources.(*values.Array)

		res.Resources = make([]drivers.ResourceFilter, 0, resources.Length())

		var e error

		resources.ForEach(func(el core.Value, idx int) bool {
			if e = core.ValidateType(el, types.Object); e != nil {
				return false
			}

			pattern := el.(*values.Object)

			url, urlExists := pattern.Get("url")
			resType, resTypeExists := pattern.Get("type")

			// ignore element
			if !urlExists && !resTypeExists {
				return true
			}

			res.Resources = append(res.Resources, drivers.ResourceFilter{
				URL:  url.String(),
				Type: resType.String(),
			})

			return true
		})

		if e != nil {
			return nil, e
		}
	}

	statusCodes, exists := ignore.Get("statusCodes")

	if exists {
		if err := core.ValidateType(statusCodes, types.Array); err != nil {
			return nil, err
		}

		statusCodes := statusCodes.(*values.Array)

		res.StatusCodes = make([]drivers.StatusCodeFilter, 0, statusCodes.Length())

		var e error

		statusCodes.ForEach(func(el core.Value, idx int) bool {
			if e = core.ValidateType(el, types.Object); e != nil {
				return false
			}

			pattern := el.(*values.Object)

			url := pattern.MustGetOr("url", values.NewString(""))
			code, codeExists := pattern.Get("code")

			// ignore element
			if !codeExists {
				e = errors.New("http code is required")
				return false
			}

			res.StatusCodes = append(res.StatusCodes, drivers.StatusCodeFilter{
				URL:  url.String(),
				Code: int(values.ToInt(code)),
			})

			return true
		})
	}

	return res, nil
}
