package http

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
	"hash/fnv"
)

type HTMLPage struct {
	document *HTMLDocument
	cookies  []drivers.HTTPCookie
	frames   *values.Array
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
	p.frames = nil

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
	tc := drivers.Compare(p.Type(), other.Type())

	if tc != 0 {
		return tc
	}

	httpPage, ok := other.(*HTMLPage)

	if !ok {
		return 1
	}

	return p.document.GetURL().Compare(httpPage.GetURL())
}

func (p *HTMLPage) Unwrap() interface{} {
	return p
}

func (p *HTMLPage) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("HTTP"))
	h.Write([]byte(p.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(p.document.GetURL()))

	return h.Sum64()
}

func (p *HTMLPage) Copy() core.Value {
	return values.None
}

func (p *HTMLPage) Iterate(ctx context.Context) (core.Iterator, error) {
	return p.document.Iterate(ctx)
}

func (p *HTMLPage) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInPage(ctx, p, path)
}

func (p *HTMLPage) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInPage(ctx, p, path, value)
}

func (p *HTMLPage) Length() values.Int {
	return p.document.Length()
}

func (p *HTMLPage) Close() error {
	return nil
}

func (p *HTMLPage) IsClosed() values.Boolean {
	return values.False
}

func (p *HTMLPage) GetURL() values.String {
	return p.document.GetURL()
}

func (p *HTMLPage) GetMainFrame() drivers.HTMLDocument {
	return p.document
}

func (p *HTMLPage) GetFrames(ctx context.Context) (*values.Array, error) {
	if p.frames == nil {
		arr := values.NewArray(10)

		err := common.CollectFrames(ctx, arr, p.document)

		if err != nil {
			return nil, err
		}

		p.frames = arr
	}

	return p.frames, nil
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

func (p *HTMLPage) unfoldFrames(ctx context.Context) (core.Value, error) {
	res := values.NewArray(10)

	err := common.CollectFrames(ctx, res, p.document)

	if err != nil {
		return nil, err
	}

	return res, nil
}
