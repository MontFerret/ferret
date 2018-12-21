package http

import (
	"bytes"
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"net/http"
	"net/url"

	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
)

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

func (drv *Driver) GetDocument(ctx context.Context, targetURL values.String) (values.HTMLDocument, error) {
	u := targetURL.String()
	req, err := http.NewRequest(http.MethodGet, u, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	req = req.WithContext(ctx)

	ua := common.GetUserAgent(drv.options.userAgent)

	logger := logging.FromContext(ctx)
	logger.
		Debug().
		Str("user-agent", ua).
		Msg("using User-Agent")

	// use custom user agent
	if ua != "" {
		req.Header.Set("User-Agent", uarand.GetRandom())
	}

	resp, err := drv.client.Do(req)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve a document %s", u)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse a document %s", u)
	}

	return NewHTMLDocument(u, doc)
}

func (drv *Driver) ParseDocument(_ context.Context, str values.String) (values.HTMLDocument, error) {
	buf := bytes.NewBuffer([]byte(str))

	doc, err := goquery.NewDocumentFromReader(buf)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse a document")
	}

	return NewHTMLDocument("#string", doc)
}

func (drv *Driver) Close() error {
	drv.client = nil

	return nil
}
