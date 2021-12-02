package http

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"hash/fnv"

	"github.com/PuerkitoBio/goquery"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLPage struct {
	document *HTMLDocument
	cookies  *drivers.HTTPCookies
	frames   *values.Array
	response drivers.HTTPResponse
}

func NewHTMLPage(
	qdoc *goquery.Document,
	url string,
	response drivers.HTTPResponse,
	cookies *drivers.HTTPCookies,
) (*HTMLPage, error) {
	doc, err := NewRootHTMLDocument(qdoc, url)

	if err != nil {
		return nil, err
	}

	p := new(HTMLPage)
	p.document = doc
	p.cookies = cookies
	p.frames = nil
	p.response = response

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
	return p.document
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
	var cookies *drivers.HTTPCookies

	if p.cookies != nil {
		cookies = p.cookies.Copy().(*drivers.HTTPCookies)
	}

	page, err := NewHTMLPage(
		p.document.doc,
		p.document.GetURL().String(),
		p.response,
		cookies,
	)

	if err != nil {
		return values.None
	}

	return page
}

func (p *HTMLPage) Iterate(ctx context.Context) (core.Iterator, error) {
	return p.document.Iterate(ctx)
}

func (p *HTMLPage) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	return common.GetInPage(ctx, path, p)
}

func (p *HTMLPage) SetIn(ctx context.Context, path []core.Value, value core.Value) core.PathError {
	return common.SetInPage(ctx, path, p, value)
}

func (p *HTMLPage) Length() values.Int {
	return p.document.Length()
}

func (p *HTMLPage) Close() error {
	return nil
}

func (p *HTMLPage) IsClosed() values.Boolean {
	return values.True
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
			return values.NewArray(0), err
		}

		p.frames = arr
	}

	return p.frames, nil
}

func (p *HTMLPage) GetFrame(ctx context.Context, idx values.Int) (core.Value, error) {
	if p.frames == nil {
		arr := values.NewArray(10)

		err := common.CollectFrames(ctx, arr, p.document)

		if err != nil {
			return values.None, err
		}

		p.frames = arr
	}

	return p.frames.Get(idx), nil
}

func (p *HTMLPage) GetCookies(_ context.Context) (*drivers.HTTPCookies, error) {
	res := drivers.NewHTTPCookies()

	if p.cookies != nil {
		p.cookies.ForEach(func(value drivers.HTTPCookie, _ values.String) bool {
			res.Set(value)

			return true
		})
	}

	return res, nil
}

func (p *HTMLPage) GetResponse(_ context.Context) (drivers.HTTPResponse, error) {
	return p.response, nil
}

func (p *HTMLPage) SetCookies(_ context.Context, _ *drivers.HTTPCookies) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) DeleteCookies(_ context.Context, _ *drivers.HTTPCookies) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) PrintToPDF(_ context.Context, _ drivers.PDFParams) (values.Binary, error) {
	return nil, core.ErrNotSupported
}

func (p *HTMLPage) CaptureScreenshot(_ context.Context, _ drivers.ScreenshotParams) (values.Binary, error) {
	return nil, core.ErrNotSupported
}

func (p *HTMLPage) WaitForNavigation(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) WaitForFrameNavigation(_ context.Context, _ drivers.HTMLDocument, _ values.String) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) Navigate(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (p *HTMLPage) NavigateBack(_ context.Context, _ values.Int) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (p *HTMLPage) NavigateForward(_ context.Context, _ values.Int) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (p *HTMLPage) Subscribe(_ context.Context, _ events.Subscription) (events.Stream, error) {
	return nil, core.ErrNotSupported
}
