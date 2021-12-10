package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/gobwas/glob"

	"golang.org/x/net/html/charset"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sethgrid/pester"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
)

const DriverName = "http"

type Driver struct {
	client  *pester.Client
	options *Options
}

func NewDriver(opts ...Option) *Driver {
	drv := new(Driver)
	drv.options = NewOptions(opts)

	drv.client = newHTTPClient(drv.options)

	return drv
}

func newHTTPClient(options *Options) (httpClient *pester.Client) {
	httpClient = pester.New()

	httpClient.Concurrency = options.Concurrency
	httpClient.MaxRetries = options.MaxRetries
	httpClient.Backoff = options.Backoff
	httpClient.Timeout = options.Timeout

	if options.HTTPTransport != nil {
		httpClient.Transport = options.HTTPTransport
	}

	if options.Proxy == "" {
		return
	}

	if err := addProxy(httpClient, options.Proxy); err != nil {
		return
	}

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

	params = drivers.SetDefaultParams(drv.options.Options, params)

	drv.makeRequest(ctx, req, params)

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

	body := io.Reader(resp.Body)
	if params.Charset != "" {
		body, err = drv.convertToUTF8(body, params.Charset)
		if err != nil {
			return nil, errors.Wrapf(err, "failed convert to UTF-8 a document %s", params.URL)
		}
	}

	doc, err := goquery.NewDocumentFromReader(body)
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
		Headers:    drivers.NewHTTPHeadersWith(resp.Header),
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

func (drv *Driver) convertToUTF8(reader io.Reader, srcCharset string) (data io.Reader, err error) {
	data, err = charset.NewReader(reader, srcCharset)
	if err != nil {
		return nil, err
	}

	return
}

func (drv *Driver) makeRequest(ctx context.Context, req *http.Request, params drivers.Params) {
	logger := logging.FromContext(ctx)

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	if params.Headers != nil {
		params.Headers.ForEach(func(value []string, key string) bool {
			v := params.Headers.Get(key)

			req.Header.Set(key, v)

			logger.
				Debug().
				Timestamp().
				Str("header", key).
				Msg("set header")

			return true
		})
	}

	if params.Cookies != nil {
		params.Cookies.ForEach(func(value drivers.HTTPCookie, key values.String) bool {
			v, exist := params.Cookies.Get(key)
			if !exist {
				return false
			}

			req.AddCookie(fromDriverCookie(v))

			logger.
				Debug().
				Timestamp().
				Str("cookie", key.String()).
				Msg("set cookie")

			return true
		})
	}

	ua := common.GetUserAgent(params.UserAgent)
	logger.
		Debug().
		Timestamp().
		Str("user-agent", ua).
		Msg("using User-Agent")

	if ua != "" {
		req.Header.Set("User-Agent", ua)
	}

	req = req.WithContext(ctx)
}
