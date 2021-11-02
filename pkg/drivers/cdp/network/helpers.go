package network

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/rs/zerolog/log"

	"github.com/MontFerret/ferret/pkg/drivers"
)

var emptyExpires = time.Time{}

func toDriverBody(body *string) []byte {
	if body == nil {
		return nil
	}

	return []byte(*body)
}

func toDriverHeaders(headers network.Headers) *drivers.HTTPHeaders {
	result := drivers.NewHTTPHeaders()
	deserialized := make(map[string]string)

	if len(headers) > 0 {
		err := json.Unmarshal(headers, &deserialized)

		if err != nil {
			log.Trace().Err(err).Msg("failed to deserialize responseReceivedEvent headers")
		}
	}

	for key, value := range deserialized {
		result.Set(key, value)
	}

	return result
}

func toDriverResponse(resp network.Response, body []byte) *drivers.HTTPResponse {
	return &drivers.HTTPResponse{
		URL:          resp.URL,
		StatusCode:   resp.Status,
		Status:       resp.StatusText,
		Headers:      toDriverHeaders(resp.Headers),
		Body:         body,
		ResponseTime: float64(resp.ResponseTime),
	}
}

func toDriverRequest(req network.Request) *drivers.HTTPRequest {
	return &drivers.HTTPRequest{
		URL:     req.URL,
		Method:  req.Method,
		Headers: toDriverHeaders(req.Headers),
		Body:    toDriverBody(req.PostData),
	}
}

func fromDriverCookie(url string, cookie drivers.HTTPCookie) network.CookieParam {
	sameSite := network.CookieSameSiteNotSet

	switch cookie.SameSite {
	case drivers.SameSiteLaxMode:
		sameSite = network.CookieSameSiteLax
	case drivers.SameSiteStrictMode:
		sameSite = network.CookieSameSiteStrict
	}

	if cookie.Expires == emptyExpires {
		cookie.Expires = time.Now().Add(time.Duration(24) + time.Hour)
	}

	normalizedURL := normalizeCookieURL(url)

	return network.CookieParam{
		URL:      &normalizedURL,
		Name:     cookie.Name,
		Value:    cookie.Value,
		Secure:   &cookie.Secure,
		Path:     &cookie.Path,
		Domain:   &cookie.Domain,
		HTTPOnly: &cookie.HTTPOnly,
		SameSite: sameSite,
		Expires:  network.TimeSinceEpoch(cookie.Expires.Unix()),
	}
}

func fromDriverCookieDelete(url string, cookie drivers.HTTPCookie) *network.DeleteCookiesArgs {
	normalizedURL := normalizeCookieURL(url)

	return &network.DeleteCookiesArgs{
		URL:    &normalizedURL,
		Name:   cookie.Name,
		Path:   &cookie.Path,
		Domain: &cookie.Domain,
	}
}

func toDriverCookie(c network.Cookie) drivers.HTTPCookie {
	sameSite := drivers.SameSiteDefaultMode

	switch c.SameSite {
	case network.CookieSameSiteLax:
		sameSite = drivers.SameSiteLaxMode
	case network.CookieSameSiteStrict:
		sameSite = drivers.SameSiteStrictMode
	}

	return drivers.HTTPCookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Expires:  time.Unix(int64(c.Expires), 0),
		SameSite: sameSite,
		Secure:   c.Secure,
		HTTPOnly: c.HTTPOnly,
	}
}

func normalizeCookieURL(url string) string {
	const httpPrefix = "http://"
	const httpsPrefix = "https://"

	if strings.HasPrefix(url, httpPrefix) || strings.HasPrefix(url, httpsPrefix) {
		return url
	}

	return httpPrefix + url
}

func isURLMatched(url string, pattern *regexp.Regexp) bool {
	var matched bool

	// if a URL pattern is provided
	if pattern != nil {
		matched = pattern.MatchString(url)
	} else {
		// otherwise, just match
		matched = true
	}

	return matched
}

func isFrameMatched(current, target page.FrameID) bool {
	// if frameID is empty string or equals to the current one
	if len(target) == 0 {
		return true
	}

	return target == current
}
