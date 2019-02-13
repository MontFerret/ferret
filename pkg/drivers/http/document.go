package http

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/PuerkitoBio/goquery"
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
	return types.HTMLDocument
}

func (doc *HTMLDocument) Compare(other core.Value) int64 {
	if other.Type() == types.HTMLDocument {
		otherDoc := other.(values.HTMLDocument)

		return doc.url.Compare(otherDoc.URL())
	}

	return types.Compare(other.Type(), types.HTMLDocument)
}

func (doc *HTMLDocument) URL() core.Value {
	return doc.url
}
