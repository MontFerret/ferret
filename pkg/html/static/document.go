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

func (doc *HTMLDocument) Type() core.Type {
	return core.HTMLDocumentType
}

func (doc *HTMLDocument) Compare(other core.Value) int {
	switch other.Type() {
	case core.HTMLDocumentType:
		otherDoc := other.(values.HTMLDocument)

		return doc.url.Compare(otherDoc.URL())
	default:
		if other.Type() > core.HTMLDocumentType {
			return -1
		}

		return 1
	}
}

func (doc *HTMLDocument) URL() core.Value {
	return doc.url
}

func (doc *HTMLDocument) InnerHTMLBySelector(selector values.String) values.String {
	selection := doc.selection.Find(selector.String())

	str, err := selection.Html()

	// TODO: log error
	if err != nil {
		return values.EmptyString
	}

	return values.NewString(str)
}

func (doc *HTMLDocument) InnerHTMLBySelectorAll(selector values.String) *values.Array {
	selection := doc.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		str, err := selection.Html()

		// TODO: log error
		if err == nil {
			arr.Push(values.NewString(str))
		}
	})

	return arr
}

func (doc *HTMLDocument) InnerTextBySelector(selector values.String) values.String {
	selection := doc.selection.Find(selector.String())

	return values.NewString(selection.Text())
}

func (doc *HTMLDocument) InnerTextBySelectorAll(selector values.String) *values.Array {
	selection := doc.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		arr.Push(values.NewString(selection.Text()))
	})

	return arr
}
