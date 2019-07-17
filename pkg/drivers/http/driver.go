package http

import (
	"bytes"
	"context"
	"net/http"
	"net/url"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
)

const DriverName = "http"

type Driver struct {
	client  *pester.Client
	options *Options
}

func NewDriver(opts ...Option) *Driver {
	drv := new(Driver)
	drv.options = newOptions(opts)

	if drv.options.proxy == "" {
		drv.client = pester.New()
	} else {
		client, err := newClientWithProxy(drv.options)

		if err != nil {
			drv.client = pester.New()
		} else {
			drv.client = pester.NewExtendedClient(client)
		}
	}

	drv.client.Concurrency = drv.options.concurrency
	drv.client.MaxRetries = drv.options.maxRetries
	drv.client.Backoff = drv.options.backoff

	return drv
}

func newClientWithProxy(options *Options) (*http.Client, error) {
	proxyURL, err := url.Parse(options.proxy)

	if err != nil {
		return nil, err
	}

	proxy := http.ProxyURL(proxyURL)
	tr := &http.Transport{Proxy: proxy}

	return &http.Client{Transport: tr}, nil
}

func (drv *Driver) Name() string {
	return DriverName
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

	if params.Header != nil {
		for k := range params.Header {
			req.Header.Add(k, params.Header.Get(k))

			logger.
				Debug().
				Timestamp().
				Str("header", k).
				Msg("set header")
		}
	}

	if params.Cookies != nil {
		for _, c := range params.Cookies {
			req.AddCookie(&http.Cookie{
				Name:  c.Name,
				Value: c.Value,
			})

			logger.
				Debug().
				Timestamp().
				Str("cookie", c.Name).
				Msg("set cookie")
		}
	}

	req = req.WithContext(ctx)

	var ua string

	if params.UserAgent != "" {
		ua = common.GetUserAgent(params.UserAgent)
	} else {
		ua = common.GetUserAgent(drv.options.userAgent)
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse a document %s", params.URL)
	}

	return NewHTMLPage(doc, params.URL, params.Cookies)
}

func (drv *Driver) Parse(_ context.Context, str values.String) (drivers.HTMLPage, error) {
	buf := bytes.NewBuffer([]byte(str))

	doc, err := goquery.NewDocumentFromReader(buf)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse a document")
	}

	return NewHTMLPage(doc, "#blank", nil)
}

func (drv *Driver) Close() error {
	drv.client = nil

	return nil
}
