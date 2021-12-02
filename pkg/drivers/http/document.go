package http

import (
	"context"
	"hash/fnv"

	"github.com/PuerkitoBio/goquery"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLDocument struct {
	doc      *goquery.Document
	element  drivers.HTMLElement
	url      values.String
	parent   drivers.HTMLDocument
	children *values.Array
}

func NewRootHTMLDocument(
	node *goquery.Document,
	url string,
) (*HTMLDocument, error) {
	return NewHTMLDocument(node, url, nil)
}

func NewHTMLDocument(
	node *goquery.Document,
	url string,
	parent drivers.HTMLDocument,
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

	doc := new(HTMLDocument)
	doc.doc = node
	doc.element = el
	doc.parent = parent
	doc.url = values.NewString(url)
	doc.children = values.NewArray(10)

	frames := node.Find("iframe")
	frames.Each(func(i int, selection *goquery.Selection) {
		child, _ := NewHTMLDocument(goquery.NewDocumentFromNode(selection.Nodes[0]), selection.AttrOr("src", url), doc)

		doc.children.Push(child)
	})

	return doc, nil
}

func (doc *HTMLDocument) MarshalJSON() ([]byte, error) {
	return doc.element.MarshalJSON()
}

func (doc *HTMLDocument) Type() core.Type {
	return drivers.HTMLDocumentType
}

func (doc *HTMLDocument) String() string {
	str, err := doc.doc.Html()

	if err != nil {
		return ""
	}

	return str
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

func (doc *HTMLDocument) Unwrap() interface{} {
	return doc.doc
}

func (doc *HTMLDocument) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(doc.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(doc.url))

	return h.Sum64()
}

func (doc *HTMLDocument) Copy() core.Value {
	cp, err := NewHTMLDocument(doc.doc, string(doc.url), doc.parent)

	if err != nil {
		return values.None
	}

	return cp
}

func (doc *HTMLDocument) Clone() core.Cloneable {
	cloned, err := NewHTMLDocument(doc.doc, doc.url.String(), doc.parent)

	if err != nil {
		return values.None
	}

	return cloned
}

func (doc *HTMLDocument) Length() values.Int {
	return values.NewInt(doc.doc.Length())
}

func (doc *HTMLDocument) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(doc.element)
}

func (doc *HTMLDocument) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	return common.GetInDocument(ctx, path, doc)
}

func (doc *HTMLDocument) SetIn(ctx context.Context, path []core.Value, value core.Value) core.PathError {
	return common.SetInDocument(ctx, path, doc, value)
}

func (doc *HTMLDocument) GetNodeType(_ context.Context) (values.Int, error) {
	return 9, nil
}

func (doc *HTMLDocument) GetNodeName(_ context.Context) (values.String, error) {
	return "#document", nil
}

func (doc *HTMLDocument) GetChildNodes(ctx context.Context) (*values.Array, error) {
	return doc.element.GetChildNodes(ctx)
}

func (doc *HTMLDocument) GetChildNode(ctx context.Context, idx values.Int) (core.Value, error) {
	return doc.element.GetChildNode(ctx, idx)
}

func (doc *HTMLDocument) QuerySelector(ctx context.Context, selector drivers.QuerySelector) (core.Value, error) {
	return doc.element.QuerySelector(ctx, selector)
}

func (doc *HTMLDocument) QuerySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*values.Array, error) {
	return doc.element.QuerySelectorAll(ctx, selector)
}

func (doc *HTMLDocument) CountBySelector(ctx context.Context, selector drivers.QuerySelector) (values.Int, error) {
	return doc.element.CountBySelector(ctx, selector)
}

func (doc *HTMLDocument) ExistsBySelector(ctx context.Context, selector drivers.QuerySelector) (values.Boolean, error) {
	return doc.element.ExistsBySelector(ctx, selector)
}

func (doc *HTMLDocument) XPath(ctx context.Context, expression values.String) (core.Value, error) {
	return doc.element.XPath(ctx, expression)
}

func (doc *HTMLDocument) GetTitle() values.String {
	title := doc.doc.Find("head > title")

	return values.NewString(title.Text())
}

func (doc *HTMLDocument) GetChildDocuments(_ context.Context) (*values.Array, error) {
	return doc.children.Clone().(*values.Array), nil
}

func (doc *HTMLDocument) GetURL() values.String {
	return doc.url
}

func (doc *HTMLDocument) GetElement() drivers.HTMLElement {
	return doc.element
}

func (doc *HTMLDocument) GetName() values.String {
	return ""
}

func (doc *HTMLDocument) GetParentDocument(_ context.Context) (drivers.HTMLDocument, error) {
	return doc.parent, nil
}

func (doc *HTMLDocument) ScrollTop(_ context.Context, _ drivers.ScrollOptions) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) ScrollBottom(_ context.Context, _ drivers.ScrollOptions) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) ScrollBySelector(_ context.Context, _ drivers.QuerySelector, _ drivers.ScrollOptions) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) Scroll(_ context.Context, _ drivers.ScrollOptions) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) MoveMouseByXY(_ context.Context, _, _ values.Float) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) Close() error {
	return nil
}
