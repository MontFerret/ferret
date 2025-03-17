package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type PageLoadParams struct {
	drivers.Params
	Driver  string
	Timeout time.Duration
}

// DOCUMENT opens an HTML page by a given url.
// By default, loads a page by http call - resulted page does not support any interactions.
// @param {hashMap} [params] - An object containing the following properties :
// @param {String} [params.driver] - Driver name to use.
// @param {Int} [params.timeout=60000] - Page load timeout.
// @param {String} [params.userAgent] - Custom user agent.
// @param {Boolean} [params.keepCookies=False] - Boolean value indicating whether to use cookies from previous sessions i.e. not to open a page in the Incognito mode.
// @param {hashMap[] | hashMap} [params.cookies] - Set of HTTP cookies to use during page loading.
// @param {String} params.cookies.*.name - Cookie name.
// @param {String} params.cookies.*.value - Cookie value.
// @param {String} params.cookies.*.path - Cookie path.
// @param {String} params.cookies.*.domain - Cookie domain.
// @param {Int} [params.cookies.*.maxAge] - Cookie max age.
// @param {String|DateTime} [params.cookies.*.expires] - Cookie expiration date time.
// @param {String} [params.cookies.*.sameSite] - Cookie cross-origin policy.
// @param {Boolean} [params.cookies.*.httpOnly=false] - Cookie cannot be accessed through client side script.
// @param {Boolean} [params.cookies.*.secure=false] - Cookie sent to the server only with an encrypted request over the HTTPS protocol.
// @param {hashMap} [params.headers] - Set of HTTP headers to use during page loading.
// @param {hashMap} [params.ignore] - Set of parameters to ignore some page functionality or behavior.
// @param {hashMap[]} [params.ignore.resources] - Collection of rules to ignore resources during page load and navigation.
// @param {String} [params.ignore.resources.*.url] - Resource url pattern. If set, requests for matching urls will be blocked. Wildcards ('*' -> zero or more, '?' -> exactly one) are allowed. Escape character is backslash. Omitting is equivalent to "*".
// @param {String} [params.ignore.resources.*.type] - Resource type. If set, requests for matching resource types will be blocked.
// @param {hashMap[]} [params.ignore.statusCodes] - Collection of rules to ignore certain HTTP codes that can cause failures.
// @param {String} [params.ignore.statusCodes.*.url] - Url pattern. If set, codes for matching urls will be ignored. Wildcards ('*' -> zero or more, '?' -> exactly one) are allowed. Escape character is backslash. Omitting is equivalent to "*".
// @param {Int} [params.ignore.statusCodes.*.code] - HTTP code to ignore.
// @param {hashMap} [params.viewport] - Viewport params.
// @param {Int} [params.viewport.height] - Viewport height.
// @param {Int} [params.viewport.width] - Viewport width.
// @param {Float} [params.viewport.scaleFactor] - Viewport scale factor.
// @param {Boolean} [params.viewport.mobile] - Second that indicates whether to emulate mobile device.
// @param {Boolean} [params.viewport.landscape] - Second that indicates whether to render a page in landscape position.
// @param {String} [params.charset] - (only HTTPDriver) Source charset content to convert UTF-8.
// @return {HTMLPage} - Loaded HTML page.
func Open(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return core.None, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return core.None, err
	}

	url := args[0].(core.String)

	var params PageLoadParams

	if len(args) == 1 {
		params = newDefaultDocLoadParams(url)
	} else {
		p, err := newPageLoadParams(url, args[1])

		if err != nil {
			return core.None, err
		}

		params = p
	}

	ctx, cancel := context.WithTimeout(ctx, params.Timeout)
	defer cancel()

	drv, err := drivers.FromContext(ctx, params.Driver)

	if err != nil {
		return core.None, err
	}

	return drv.Open(ctx, params.Params)
}

func newDefaultDocLoadParams(url core.String) PageLoadParams {
	return PageLoadParams{
		Params: drivers.Params{
			URL: url.String(),
		},
		Timeout: drivers.DefaultPageLoadTimeout * time.Millisecond,
	}
}

func newPageLoadParams(url core.String, arg core.Value) (PageLoadParams, error) {
	res := newDefaultDocLoadParams(url)

	if err := core.ValidateType(arg, types.Boolean, types.String, types.Object); err != nil {
		return res, err
	}

	switch arg.Type() {
	case types.Object:
		obj := arg.(*internal.Object)

		driver, exists := obj.Get(core.NewString("driver"))

		if exists {
			if err := core.ValidateType(driver, types.String); err != nil {
				return res, err
			}

			res.Driver = driver.(core.String).String()
		}

		timeout, exists := obj.Get(core.NewString("timeout"))

		if exists {
			if err := core.ValidateType(timeout, types.Int); err != nil {
				return res, err
			}

			res.Timeout = time.Duration(timeout.(core.Int)) * time.Millisecond
		}

		userAgent, exists := obj.Get(core.NewString("userAgent"))

		if exists {
			if err := core.ValidateType(userAgent, types.String); err != nil {
				return res, err
			}

			res.UserAgent = userAgent.String()
		}

		keepCookies, exists := obj.Get(core.NewString("keepCookies"))

		if exists {
			if err := core.ValidateType(keepCookies, types.Boolean); err != nil {
				return res, err
			}

			res.KeepCookies = bool(keepCookies.(core.Boolean))
		}

		cookies, exists := obj.Get(core.NewString("cookies"))

		if exists {
			if err := core.ValidateType(cookies, types.Array, types.Object); err != nil {
				return res, err
			}

			switch c := cookies.(type) {
			case *internal.Array:
				cookies, err := parseCookieArray(c)

				if err != nil {
					return res, err
				}

				res.Cookies = cookies
			case *internal.Object:
				cookies, err := parseCookieObject(c)

				if err != nil {
					return res, err
				}

				res.Cookies = cookies
			default:
				res.Cookies = drivers.NewHTTPCookies()
			}
		}

		headers, exists := obj.Get(core.NewString("headers"))

		if exists {
			if err := core.ValidateType(headers, types.Object); err != nil {
				return res, err
			}

			header := parseHeader(headers.(*internal.Object))
			res.Headers = header
		}

		viewport, exists := obj.Get(core.NewString("viewport"))

		if exists {
			viewport, err := parseViewport(viewport)

			if err != nil {
				return res, err
			}

			res.Viewport = viewport
		}

		ignore, exists := obj.Get(core.NewString("ignore"))

		if exists {
			ignore, err := parseIgnore(ignore)

			if err != nil {
				return res, err
			}

			res.Ignore = ignore
		}

		charset, exists := obj.Get(core.NewString("charset"))

		if exists {
			if err := core.ValidateType(charset, types.String); err != nil {
				return res, err
			}

			res.Charset = charset.String()
		}
	case types.String:
		res.Driver = arg.(core.String).String()
	case types.Boolean:
		b := arg.(core.Boolean)

		// fallback
		if b {
			res.Driver = cdp.DriverName
		}
	}

	return res, nil
}

func parseCookieObject(obj *internal.Object) (*drivers.HTTPCookies, error) {
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

func parseCookieArray(arr *internal.Array) (*drivers.HTTPCookies, error) {
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

	co := value.(*internal.Object)

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

		cookie.MaxAge = int(maxAge.(core.Int))
	}

	expires, exists := co.Get("expires")

	if exists {
		if err = core.ValidateType(expires, types.DateTime, types.String); err != nil {
			return drivers.HTTPCookie{}, err
		}

		if expires.Type() == types.DateTime {
			cookie.Expires = expires.(core.DateTime).Unwrap().(time.Time)
		} else {
			t, err := time.Parse(core.DefaultTimeLayout, expires.String())

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

		cookie.HTTPOnly = bool(httpOnly.(core.Boolean))
	}

	secure, exists := co.Get("secure")

	if exists {
		if err = core.ValidateType(secure, types.Boolean); err != nil {
			return drivers.HTTPCookie{}, err
		}

		cookie.Secure = bool(secure.(core.Boolean))
	}

	return cookie, err
}

func parseHeader(headers *internal.Object) *drivers.HTTPHeaders {
	res := drivers.NewHTTPHeaders()

	headers.ForEach(func(value core.Value, key string) bool {
		if value.Type() == types.Array {
			value := value.(*internal.Array)

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

	viewport := value.(*internal.Object)

	width, exists := viewport.Get(core.NewString("width"))

	if exists {
		if err := core.ValidateType(width, types.Int); err != nil {
			return nil, err
		}

		res.Width = int(internal.ToInt(width))
	}

	height, exists := viewport.Get(core.NewString("height"))

	if exists {
		if err := core.ValidateType(height, types.Int); err != nil {
			return nil, err
		}

		res.Height = int(internal.ToInt(height))
	}

	mobile, exists := viewport.Get(core.NewString("mobile"))

	if exists {
		res.Mobile = bool(internal.ToBoolean(mobile))
	}

	landscape, exists := viewport.Get(core.NewString("landscape"))

	if exists {
		res.Landscape = bool(internal.ToBoolean(landscape))
	}

	scaleFactor, exists := viewport.Get(core.NewString("scaleFactor"))

	if exists {
		res.ScaleFactor = float64(internal.ToFloat(scaleFactor))
	}

	return res, nil
}

func parseIgnore(value core.Value) (*drivers.Ignore, error) {
	if err := core.ValidateType(value, types.Object); err != nil {
		return nil, err
	}

	res := &drivers.Ignore{}

	ignore := value.(*internal.Object)

	resources, exists := ignore.Get("resources")

	if exists {
		if err := core.ValidateType(resources, types.Array); err != nil {
			return nil, err
		}

		resources := resources.(*internal.Array)

		res.Resources = make([]drivers.ResourceFilter, 0, resources.Length())

		var e error

		resources.ForEach(func(el core.Value, idx int) bool {
			if e = core.ValidateType(el, types.Object); e != nil {
				return false
			}

			pattern := el.(*internal.Object)

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

		statusCodes := statusCodes.(*internal.Array)

		res.StatusCodes = make([]drivers.StatusCodeFilter, 0, statusCodes.Length())

		var e error

		statusCodes.ForEach(func(el core.Value, idx int) bool {
			if e = core.ValidateType(el, types.Object); e != nil {
				return false
			}

			pattern := el.(*internal.Object)

			url := pattern.MustGetOr("url", core.NewString(""))
			code, codeExists := pattern.Get("code")

			// ignore element
			if !codeExists {
				e = errors.New("http code is required")
				return false
			}

			res.StatusCodes = append(res.StatusCodes, drivers.StatusCodeFilter{
				URL:  url.String(),
				Code: int(internal.ToInt(code)),
			})

			return true
		})
	}

	return res, nil
}
