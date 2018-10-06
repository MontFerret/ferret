package static

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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

func (el *HTMLDocument) Type() core.Type {
	return core.HTMLDocumentType
}

func (el *HTMLDocument) Compare(other core.Value) int {
	switch other.Type() {
	case core.HTMLDocumentType:
		otherDoc := other.(values.HTMLDocument)

		return el.url.Compare(otherDoc.URL())
	default:
		if other.Type() > core.HTMLDocumentType {
			return -1
		}

		return 1
	}
}

func (el *HTMLDocument) URL() core.Value {
	return el.url
}
