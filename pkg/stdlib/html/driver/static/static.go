package static

import (
	"bytes"
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
	httpx "net/http"
)

type Driver struct {
	client *pester.Client
}

func NewDriver(setters ...Option) *Driver {
	client := pester.New()
	client.Concurrency = 3
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialBackoff

	for _, setter := range setters {
		setter(client)
	}

	return &Driver{client}
}

func (d *Driver) GetDocument(_ context.Context, url string) (values.HTMLNode, error) {
	req, err := httpx.NewRequest(httpx.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("User-Agent", uarand.GetRandom())

	resp, err := d.client.Do(req)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve a document %s", url)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse a document %s", url)
	}

	return NewHTMLDocument(url, doc)
}

func (d *Driver) ParseDocument(_ context.Context, str string) (values.HTMLNode, error) {
	buf := bytes.NewBuffer([]byte(str))

	doc, err := goquery.NewDocumentFromReader(buf)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse a document")
	}

	return NewHTMLDocument("#string", doc)
}

func (d *Driver) Close() error {
	d.client = nil

	return nil
}
