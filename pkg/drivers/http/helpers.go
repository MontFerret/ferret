package http

import (
	"bytes"
	HTTP "net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func parseXPathNode(nav *htmlquery.NodeNavigator) (core.Value, error) {
	switch nav.NodeType() {
	case xpath.ElementNode:
		node := nav.Current()

		if node == nil {
			return values.None, nil
		}

		return NewHTMLElement(&goquery.Selection{Nodes: []*html.Node{node}})
	case xpath.RootNode:
		node := nav.Current()

		if node == nil {
			return values.None, nil
		}

		url := htmlquery.SelectAttr(node, "url")
		return NewHTMLDocument(goquery.NewDocumentFromNode(node), url, nil)
	default:
		return values.NewString(nav.Value()), nil
	}
}

func outerHTML(s *goquery.Selection) (string, error) {
	var buf bytes.Buffer

	if len(s.Nodes) > 0 {
		c := s.Nodes[0]

		err := html.Render(&buf, c)

		if err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func toDriverCookies(cookies []*HTTP.Cookie) (drivers.HTTPCookies, error) {
	res := make(drivers.HTTPCookies)

	for _, c := range cookies {
		dc, err := toDriverCookie(c)

		if err != nil {
			return nil, err
		}

		res[dc.Name] = dc
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
