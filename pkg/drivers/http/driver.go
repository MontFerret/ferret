package http

import (
	"bytes"
	"context"
	"github.com/gobwas/glob"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sethgrid/pester"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
)

const DriverName = "http"

type Driver struct {
	client  *pester.Client
	options *Options
}

func NewDriver(opts ...Option) *Driver {
	drv := new(Driver)
	drv.options = newOptions(opts)

	drv.client = newHTTPClient(drv.options)
	drv.client.Concurrency = drv.options.Concurrency
	drv.client.MaxRetries = drv.options.MaxRetries
	drv.client.Backoff = drv.options.Backoff

	return drv
}

func newHTTPClient(options *Options) (httpClient *pester.Client) {
	httpClient = pester.New()

	if options.HTTPTransport != nil {
		httpClient.Transport = options.HTTPTransport
	}

	if options.Proxy == "" {
		return
	}

	if err := addProxy(httpClient, options.Proxy); err != nil {
		return
	}

	httpClient = pester.NewExtendedClient(&http.Client{Transport: httpClient.Transport})

	return
}

func addProxy(httpClient *pester.Client, proxyStr string) error {
	if proxyStr == "" {
		return nil
	}

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return err
	}

	proxy := http.ProxyURL(proxyURL)

	if httpClient.Transport != nil {
		httpClient.Transport.(*http.Transport).Proxy = proxy

		return nil
	}

	httpClient.Transport = &http.Transport{Proxy: proxy}

	return nil
}

func (drv *Driver) Name() string {
	return drv.options.Name
}

func (drv *Driver) Open(ctx context.Context, params drivers.Params) (drivers.HTMLPage, error) {
	req, err := http.NewRequest(http.MethodGet, params.URL, nil)

	if err != nil {
		return nil, err
	}

	logger := logging.FromContext(ctx)

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	if drv.options.Headers != nil && params.Headers == nil {
		params.Headers = make(drivers.HTTPHeaders)
	}

	// Set default headers
	for k, v := range drv.options.Headers {
		_, exists := params.Headers[k]

		// do not override user's set values
		if !exists {
			params.Headers[k] = v
		}
	}

	for k := range params.Headers {
		req.Header.Add(k, params.Headers.Get(k))

		logger.
			Debug().
			Timestamp().
			Str("header", k).
			Msg("set header")
	}

	if drv.options.Cookies != nil && params.Cookies == nil {
		params.Cookies = make(drivers.HTTPCookies)
	}

	// set default cookies
	for k, v := range drv.options.Cookies {
		_, exists := params.Cookies[k]

		// do not override user's set values
		if !exists {
			params.Cookies[k] = v
		}
	}

	for _, c := range params.Cookies {
		req.AddCookie(fromDriverCookie(c))

		logger.
			Debug().
			Timestamp().
			Str("cookie", c.Name).
			Msg("set cookie")
	}

	req = req.WithContext(ctx)

	var ua string

	if params.UserAgent != "" {
		ua = common.GetUserAgent(params.UserAgent)
	} else {
		ua = common.GetUserAgent(drv.options.UserAgent)
	}

	logger.
		Debug().
		Timestamp().
		Str("user-agent", ua).
		Msg("using User-Agent")

	if ua != "" {
		req.Header.Set("User-Agent", ua)
	}

	resp, err := drv.client.Do(req)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve a document %s", params.URL)
	}

	defer resp.Body.Close()

	var queryFilters []drivers.StatusCodeFilter

	if params.Ignore != nil {
		queryFilters = params.Ignore.StatusCodes
	}

	if !drv.responseCodeAllowed(resp, queryFilters) {
		return nil, errors.New(resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse a document %s", params.URL)
	}

	cookies, err := toDriverCookies(resp.Cookies())

	if err != nil {
		return nil, err
	}

	r := drivers.HTTPResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    drivers.HTTPHeaders(resp.Header),
	}

	return NewHTMLPage(doc, params.URL, r, cookies)
}

func (drv *Driver) Parse(_ context.Context, params drivers.ParseParams) (drivers.HTMLPage, error) {
	buf := bytes.NewBuffer(params.Content)

	doc, err := goquery.NewDocumentFromReader(buf)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse a document")
	}

	return NewHTMLPage(doc, "#blank", drivers.HTTPResponse{}, nil)
}

func (drv *Driver) Close() error {
	drv.client = nil

	return nil
}

func (drv *Driver) responseCodeAllowed(resp *http.Response, additional []drivers.StatusCodeFilter) bool {
	var allowed bool
	reqURL := resp.Request.URL.String()

	// OK is by default
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return true
	}

	// Try to use those that are passed within a query
	for _, filter := range additional {
		allowed = filter.Code == resp.StatusCode

		// check url
		if allowed && filter.URL != "" {
			allowed = glob.MustCompile(filter.URL).Match(reqURL)
		}

		if allowed {
			break
		}
	}

	// if still not allowed, try the default ones
	if !allowed {
		for _, filter := range drv.options.HTTPCodesFilter {
			allowed = filter.Code == resp.StatusCode

			if allowed && filter.URL != nil {
				allowed = filter.URL.Match(reqURL)
			}

			if allowed {
				break
			}
		}
	}

	return allowed
}
