package http

import (
	"context"
	"hash/fnv"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
)

type HTMLDocument struct {
	docNode *goquery.Document
	element drivers.HTMLElement
	url     values.String
	cookies []drivers.HTTPCookie
}

func NewHTMLDocument(
	node *goquery.Document,
	url string,
	cookies []drivers.HTTPCookie,
) (drivers.HTMLDocument, error) {
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

	return &HTMLDocument{node, el, values.NewString(url), cookies}, nil
}

func (doc *HTMLDocument) MarshalJSON() ([]byte, error) {
	return doc.element.MarshalJSON()
}

func (doc *HTMLDocument) Type() core.Type {
	return drivers.HTMLDocumentType
}

func (doc *HTMLDocument) String() string {
	str, err := doc.docNode.Html()

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
	return doc.docNode
}

func (doc *HTMLDocument) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(doc.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(doc.url))

	return h.Sum64()
}

func (doc *HTMLDocument) Copy() core.Value {
	cp, err := NewHTMLDocument(doc.docNode, string(doc.url), doc.cookies)

	if err != nil {
		return values.None
	}

	return cp
}

func (doc *HTMLDocument) Clone() core.Value {
	var cookies []drivers.HTTPCookie

	if doc.cookies != nil {
		cookies = make([]drivers.HTTPCookie, 0, len(doc.cookies))

		for i, c := range doc.cookies {
			cookies[i] = c
		}
	}

	cp, err := NewHTMLDocument(goquery.CloneDocument(doc.docNode), string(doc.url), cookies)

	if err != nil {
		return values.None
	}

	return cp
}

func (doc *HTMLDocument) Length() values.Int {
	return values.NewInt(doc.docNode.Length())
}

func (doc *HTMLDocument) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(doc.element)
}

func (doc *HTMLDocument) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInDocument(ctx, doc, path)
}

func (doc *HTMLDocument) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInDocument(ctx, doc, path, value)
}

func (doc *HTMLDocument) NodeType() values.Int {
	return 9
}

func (doc *HTMLDocument) NodeName() values.String {
	return "#document"
}

func (doc *HTMLDocument) GetChildNodes(ctx context.Context) core.Value {
	return doc.element.GetChildNodes(ctx)
}

func (doc *HTMLDocument) GetChildNode(ctx context.Context, idx values.Int) core.Value {
	return doc.element.GetChildNode(ctx, idx)
}

func (doc *HTMLDocument) QuerySelector(ctx context.Context, selector values.String) core.Value {
	return doc.element.QuerySelector(ctx, selector)
}

func (doc *HTMLDocument) QuerySelectorAll(ctx context.Context, selector values.String) core.Value {
	return doc.element.QuerySelectorAll(ctx, selector)
}

func (doc *HTMLDocument) CountBySelector(ctx context.Context, selector values.String) values.Int {
	return doc.element.CountBySelector(ctx, selector)
}

func (doc *HTMLDocument) ExistsBySelector(ctx context.Context, selector values.String) values.Boolean {
	return doc.element.ExistsBySelector(ctx, selector)
}

func (doc *HTMLDocument) DocumentElement() drivers.HTMLElement {
	return doc.element
}

func (doc *HTMLDocument) GetURL() core.Value {
	return doc.url
}

func (doc *HTMLDocument) SetURL(_ context.Context, _ values.String) error {
	return core.ErrInvalidOperation
}

func (doc *HTMLDocument) GetCookies(_ context.Context) (*values.Array, error) {
	if doc.cookies == nil {
		return values.NewArray(0), nil
	}

	arr := values.NewArray(len(doc.cookies))

	for _, c := range doc.cookies {
		arr.Push(c)
	}

	return arr, nil
}

func (doc *HTMLDocument) SetCookies(_ context.Context, _ ...drivers.HTTPCookie) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) DeleteCookies(_ context.Context, _ ...drivers.HTTPCookie) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) Navigate(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) NavigateBack(_ context.Context, _ values.Int) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (doc *HTMLDocument) NavigateForward(_ context.Context, _ values.Int) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (doc *HTMLDocument) ClickBySelector(_ context.Context, _ values.String) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (doc *HTMLDocument) ClickBySelectorAll(_ context.Context, _ values.String) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (doc *HTMLDocument) InputBySelector(_ context.Context, _ values.String, _ core.Value, _ values.Int) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (doc *HTMLDocument) SelectBySelector(_ context.Context, _ values.String, _ *values.Array) (*values.Array, error) {
	return nil, core.ErrNotSupported
}

func (doc *HTMLDocument) PrintToPDF(_ context.Context, _ drivers.PDFParams) (values.Binary, error) {
	return nil, core.ErrNotSupported
}

func (doc *HTMLDocument) CaptureScreenshot(_ context.Context, _ drivers.ScreenshotParams) (values.Binary, error) {
	return nil, core.ErrNotSupported
}

func (doc *HTMLDocument) ScrollTop(_ context.Context) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) ScrollBottom(_ context.Context) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) ScrollBySelector(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) ScrollByXY(_ context.Context, _, _ values.Float) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) MoveMouseBySelector(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) MoveMouseByXY(_ context.Context, _, _ values.Float) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) WaitForNavigation(_ context.Context) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) WaitForElement(_ context.Context, _ values.String, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) WaitForClassBySelector(_ context.Context, _, _ values.String, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) WaitForClassBySelectorAll(_ context.Context, _, _ values.String, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) WaitForAttributeBySelector(_ context.Context, _, _ values.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) WaitForAttributeBySelectorAll(_ context.Context, _, _ values.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) WaitForStyleBySelector(_ context.Context, _, _ values.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) WaitForStyleBySelectorAll(_ context.Context, _, _ values.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (doc *HTMLDocument) Close() error {
	return nil
}
