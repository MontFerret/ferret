package static

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
)

type HtmlDocument struct {
	*HtmlElement
	url values.String
}

func NewHtmlDocument(
	url string,
	node *goquery.Document,
) (*HtmlDocument, error) {
	if url == "" {
		return nil, core.Error(core.ErrMissedArgument, "document url")
	}

	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "document root selection")
	}

	el, err := NewHtmlElement(node.Selection)

	if err != nil {
		return nil, err
	}

	return &HtmlDocument{el, values.NewString(url)}, nil
}

func (el *HtmlDocument) Type() core.Type {
	return core.HtmlDocumentType
}

func (el *HtmlDocument) Compare(other core.Value) int {
	switch other.Type() {
	case core.HtmlDocumentType:
		otherDoc := other.(values.HtmlDocument)

		return el.url.Compare(otherDoc.Url())
	default:
		if other.Type() > core.HtmlDocumentType {
			return -1
		}

		return 1
	}
}

func (el *HtmlDocument) Url() core.Value {
	return el.url
}
