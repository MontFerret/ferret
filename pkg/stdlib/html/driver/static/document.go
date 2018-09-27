package static

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/PuerkitoBio/goquery"
)

type HtmlDocument struct {
	*HtmlElement
	url string
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

	return &HtmlDocument{el, url}, nil
}

func (el *HtmlDocument) Type() core.Type {
	return core.HtmlDocumentType
}

func (el *HtmlDocument) Compare(other core.Value) int {
	switch other.Type() {
	case core.HtmlDocumentType:
		// TODO: complete the comparison
		return -1
	default:
		if other.Type() > core.HtmlDocumentType {
			return -1
		}

		return 1
	}
}
