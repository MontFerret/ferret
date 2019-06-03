package http

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
)

type HTMLPage struct {
	document *HTMLDocument
	cookies  []drivers.HTTPCookie
}

func NewHTMLPage(
	qdoc *goquery.Document,
	url string,
	cookies []drivers.HTTPCookie,
) (*HTMLPage, error) {
	doc, err := NewRootHTMLDocument(qdoc, url)

	if err != nil {
		return nil, err
	}

	p := new(HTMLPage)
	p.document = doc
	p.cookies = cookies

	return p, nil
}

func (p *HTMLPage) MarshalJSON() ([]byte, error) {
	return p.document.MarshalJSON()
}

func (p *HTMLPage) Type() core.Type {
	return drivers.HTMLPageType
}

func (p *HTMLPage) String() string {
	return p.document.GetURL().String()
}

func (p *HTMLPage) Compare(other core.Value) int64 {
	panic("implement me")
}

func (p *HTMLPage) Unwrap() interface{} {
	return p
}

func (p *HTMLPage) Hash() uint64 {
	panic("implement me")
}

func (p *HTMLPage) Copy() core.Value {
	panic("implement me")
}

func (p *HTMLPage) Iterate(ctx context.Context) (core.Iterator, error) {
	panic("implement me")
}

func (p *HTMLPage) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	panic("implement me")
}

func (p *HTMLPage) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	panic("implement me")
}

func (p *HTMLPage) Length() values.Int {
	panic("implement me")
}

func (p *HTMLPage) Close() error {
	panic("implement me")
}

func (p *HTMLPage) IsClosed() values.Boolean {
	panic("implement me")
}

func (p *HTMLPage) MainFrame() drivers.HTMLDocument {
	return p.document
}

func (p *HTMLPage) Frames(ctx context.Context) (*values.Array, error) {
	panic("implement me")
}

func (p *HTMLPage) GetCookies(ctx context.Context) (*values.Array, error) {
	if p.cookies == nil {
		return values.NewArray(0), nil
	}

	arr := values.NewArray(len(p.cookies))

	for _, c := range p.cookies {
		arr.Push(c)
	}

	return arr, nil
}

func (p *HTMLPage) SetCookies(ctx context.Context, cookies ...drivers.HTTPCookie) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) DeleteCookies(ctx context.Context, cookies ...drivers.HTTPCookie) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) PrintToPDF(ctx context.Context, params drivers.PDFParams) (values.Binary, error) {
	return nil, core.ErrNotSupported
}

func (p *HTMLPage) CaptureScreenshot(ctx context.Context, params drivers.ScreenshotParams) (values.Binary, error) {
	return nil, core.ErrNotSupported
}

func (p *HTMLPage) WaitForNavigation(ctx context.Context) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) Navigate(ctx context.Context, url values.String) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (p *HTMLPage) NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error) {
	return false, core.ErrNotSupported
}
