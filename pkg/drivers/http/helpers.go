package http

import (
	HTTP "net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func fromSelectionToNode(selection *goquery.Selection) *html.Node {
	if selection.Length() > 0 {
		return selection.Get(0)
	}

	return nil
}

func toDriverCookies(cookies []*HTTP.Cookie) (*drivers.HTTPCookies, error) {
	res := drivers.NewHTTPCookies()

	for _, c := range cookies {
		dc, err := toDriverCookie(c)

		if err != nil {
			return nil, err
		}

		res.Set(dc)
	}

	return res, nil
}

func toDriverCookie(cookie *HTTP.Cookie) (drivers.HTTPCookie, error) {
	res := drivers.HTTPCookie{}

	if cookie == nil {
		return res, core.Error(core.ErrMissedArgument, "cookie")
	}

	res.Name = cookie.Name
	res.Value = cookie.Value
	res.Path = cookie.Path
	res.Domain = cookie.Domain
	res.Expires = cookie.Expires
	res.MaxAge = cookie.MaxAge
	res.Secure = cookie.Secure
	res.HTTPOnly = cookie.HttpOnly

	switch cookie.SameSite {
	case HTTP.SameSiteLaxMode:
		res.SameSite = drivers.SameSiteLaxMode
	case HTTP.SameSiteStrictMode:
		res.SameSite = drivers.SameSiteStrictMode
	default:
		res.SameSite = drivers.SameSiteDefaultMode
	}

	return res, nil
}

func fromDriverCookie(cookie drivers.HTTPCookie) *HTTP.Cookie {
	res := &HTTP.Cookie{}

	res.Name = cookie.Name
	res.Value = cookie.Value
	res.Path = cookie.Path
	res.Domain = cookie.Domain
	res.Expires = cookie.Expires
	res.MaxAge = cookie.MaxAge
	res.Secure = cookie.Secure
	res.HttpOnly = cookie.HTTPOnly

	switch cookie.SameSite {
	case drivers.SameSiteLaxMode:
		res.SameSite = HTTP.SameSiteLaxMode
	case drivers.SameSiteStrictMode:
		res.SameSite = HTTP.SameSiteStrictMode
	default:
		res.SameSite = HTTP.SameSiteDefaultMode
	}

	return res
}
