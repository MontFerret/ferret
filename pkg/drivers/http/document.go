package http

import (
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type HTMLDocument struct {
	*HTMLElement
	url values.String
}

func NewHTMLDocument(
	url string,
	node *goquery.Document,
) (*HTMLDocument, error) {
	if url == "" {
		return nil, core.Error(core.ErrMissedArgument, "document url")
	}

	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "document root selection")
	}

	el, err := NewHTMLElement(node.Selection)

	if err != nil {
		return nil, err
	}

	return &HTMLDocument{el, values.NewString(url)}, nil
}

func (doc *HTMLDocument) Type() core.Type {
	return drivers.HTMLElementType
}

func (doc *HTMLDocument) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLElementType:
		otherDoc := other.(drivers.HTMLDocument)

		return doc.url.Compare(otherDoc.GetURL())
	default:
		return drivers.Compare(doc.Type(), other.Type())
	}
}

func (doc *HTMLDocument) GetURL() core.Value {
	return doc.url
}

func (doc *HTMLDocument) SetURL(url values.String) error {
	return errors.New("Url is read-only")
}
