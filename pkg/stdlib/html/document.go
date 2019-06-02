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

type DocumentLoadParams struct {
	drivers.OpenPageParams
	Driver  string
	Timeout time.Duration
}

// Open opens an HTML page by a given url.
// By default, loads a document by http call - resulted document does not support any interactions.
// If passed "true" as a second argument, headless browser is used for loading the document which support interactions.
// @param url (String) - Target url string. If passed "about:blank" for dynamic document - it will open an empty page.
// @param isDynamicOrParams (Boolean|DocumentLoadParams) - Either a boolean value that indicates whether to use dynamic page
// or an object with the following properties :
// 		dynamic (Boolean) - Optional, indicates whether to use dynamic page.
// 		timeout (Int) - Optional, Open load timeout.
// @returns (HTMLDocument) - Returns loaded HTML document.
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

	var params DocumentLoadParams

	if len(args) == 1 {
		params = newDefaultDocLoadParams(url)
	} else {
		p, err := newDocLoadParams(url, args[1])

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

	return drv.Open(ctx, params.OpenPageParams)
}

func newDefaultDocLoadParams(url values.String) DocumentLoadParams {
	return DocumentLoadParams{
		OpenPageParams: drivers.OpenPageParams{
			URL: url.String(),
		},
		Timeout: time.Second * 30,
	}
}

func newDocLoadParams(url values.String, arg core.Value) (DocumentLoadParams, error) {
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

			res.Timeout = time.Duration(timeout.(values.Int)) + time.Millisecond
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
			if err := core.ValidateType(cookies, types.Array); err != nil {
				return res, err
			}

			cookies, err := parseCookies(cookies.(*values.Array))

			if err != nil {
				return res, err
			}

			res.Cookies = cookies
		}

		header, exists := obj.Get(values.NewString("header"))

		if exists {
			if err := core.ValidateType(header, types.Object); err != nil {
				return res, err
			}

			header := parseHeader(header.(*values.Object))
			res.Header = header
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

func parseCookies(arr *values.Array) ([]drivers.HTTPCookie, error) {
	var err error
	res := make([]drivers.HTTPCookie, 0, arr.Length())

	arr.ForEach(func(value core.Value, idx int) bool {
		cookie, e := parseCookie(value)

		if e != nil {
			err = e

			return false
		}

		res = append(res, cookie)

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

func parseHeader(header *values.Object) drivers.HTTPHeader {
	res := make(drivers.HTTPHeader)

	header.ForEach(func(value core.Value, key string) bool {
		res.Set(key, value.String())

		return true
	})

	return res
}
