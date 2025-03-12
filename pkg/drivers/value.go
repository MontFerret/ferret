package drivers

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
)

type (
	// WaitEvent is an enum that represents what event is needed to wait for
	WaitEvent int

	// HTMLNode is an interface from which a number of DOM API object types inherit.
	// It allows those types to be treated similarly;
	// for example, inheriting the same set of methods, or being tested in the same way.
	HTMLNode interface {
		core.Value
		core.Iterable
		core.Keyed
		core.Measurable
		io.Closer

		GetNodeType(ctx context.Context) (core.Int, error)

		GetNodeName(ctx context.Context) (core.String, error)

		GetChildNodes(ctx context.Context) (*internal.Array, error)

		GetChildNode(ctx context.Context, idx core.Int) (core.Value, error)

		QuerySelector(ctx context.Context, selector QuerySelector) (core.Value, error)

		QuerySelectorAll(ctx context.Context, selector QuerySelector) (*internal.Array, error)

		CountBySelector(ctx context.Context, selector QuerySelector) (core.Int, error)

		ExistsBySelector(ctx context.Context, selector QuerySelector) (core.Boolean, error)

		XPath(ctx context.Context, expression core.String) (core.Value, error)
	}

	// HTMLElement is the most general base interface which most objects in a GetMainFrame implement.
	HTMLElement interface {
		HTMLNode

		GetInnerText(ctx context.Context) (core.String, error)

		SetInnerText(ctx context.Context, innerText core.String) error

		GetInnerHTML(ctx context.Context) (core.String, error)

		SetInnerHTML(ctx context.Context, innerHTML core.String) error

		GetValue(ctx context.Context) (core.Value, error)

		SetValue(ctx context.Context, value core.Value) error

		GetStyles(ctx context.Context) (*internal.Object, error)

		GetStyle(ctx context.Context, name core.String) (core.Value, error)

		SetStyles(ctx context.Context, values *internal.Object) error

		SetStyle(ctx context.Context, name, value core.String) error

		RemoveStyle(ctx context.Context, name ...core.String) error

		GetAttributes(ctx context.Context) (*internal.Object, error)

		GetAttribute(ctx context.Context, name core.String) (core.Value, error)

		SetAttributes(ctx context.Context, values *internal.Object) error

		SetAttribute(ctx context.Context, name, value core.String) error

		RemoveAttribute(ctx context.Context, name ...core.String) error

		GetInnerHTMLBySelector(ctx context.Context, selector QuerySelector) (core.String, error)

		SetInnerHTMLBySelector(ctx context.Context, selector QuerySelector, innerHTML core.String) error

		GetInnerHTMLBySelectorAll(ctx context.Context, selector QuerySelector) (*internal.Array, error)

		GetInnerTextBySelector(ctx context.Context, selector QuerySelector) (core.String, error)

		SetInnerTextBySelector(ctx context.Context, selector QuerySelector, innerText core.String) error

		GetInnerTextBySelectorAll(ctx context.Context, selector QuerySelector) (*internal.Array, error)

		GetPreviousElementSibling(ctx context.Context) (core.Value, error)

		GetNextElementSibling(ctx context.Context) (core.Value, error)

		GetParentElement(ctx context.Context) (core.Value, error)

		Click(ctx context.Context, count core.Int) error

		ClickBySelector(ctx context.Context, selector QuerySelector, count core.Int) error

		ClickBySelectorAll(ctx context.Context, selector QuerySelector, count core.Int) error

		Clear(ctx context.Context) error

		ClearBySelector(ctx context.Context, selector QuerySelector) error

		Input(ctx context.Context, value core.Value, delay core.Int) error

		InputBySelector(ctx context.Context, selector QuerySelector, value core.Value, delay core.Int) error

		Press(ctx context.Context, keys []core.String, count core.Int) error

		PressBySelector(ctx context.Context, selector QuerySelector, keys []core.String, count core.Int) error

		Select(ctx context.Context, value *internal.Array) (*internal.Array, error)

		SelectBySelector(ctx context.Context, selector QuerySelector, value *internal.Array) (*internal.Array, error)

		ScrollIntoView(ctx context.Context, options ScrollOptions) error

		Focus(ctx context.Context) error

		FocusBySelector(ctx context.Context, selector QuerySelector) error

		Blur(ctx context.Context) error

		BlurBySelector(ctx context.Context, selector QuerySelector) error

		Hover(ctx context.Context) error

		HoverBySelector(ctx context.Context, selector QuerySelector) error

		WaitForElement(ctx context.Context, selector QuerySelector, when WaitEvent) error

		WaitForElementAll(ctx context.Context, selector QuerySelector, when WaitEvent) error

		WaitForAttribute(ctx context.Context, name core.String, value core.Value, when WaitEvent) error

		WaitForAttributeBySelector(ctx context.Context, selector QuerySelector, name core.String, value core.Value, when WaitEvent) error

		WaitForAttributeBySelectorAll(ctx context.Context, selector QuerySelector, name core.String, value core.Value, when WaitEvent) error

		WaitForStyle(ctx context.Context, name core.String, value core.Value, when WaitEvent) error

		WaitForStyleBySelector(ctx context.Context, selector QuerySelector, name core.String, value core.Value, when WaitEvent) error

		WaitForStyleBySelectorAll(ctx context.Context, selector QuerySelector, name core.String, value core.Value, when WaitEvent) error

		WaitForClass(ctx context.Context, class core.String, when WaitEvent) error

		WaitForClassBySelector(ctx context.Context, selector QuerySelector, class core.String, when WaitEvent) error

		WaitForClassBySelectorAll(ctx context.Context, selector QuerySelector, class core.String, when WaitEvent) error
	}

	HTMLDocument interface {
		HTMLNode

		GetTitle() core.String

		GetElement() HTMLElement

		GetURL() core.String

		GetName() core.String

		GetParentDocument(ctx context.Context) (HTMLDocument, error)

		GetChildDocuments(ctx context.Context) (*internal.Array, error)

		Scroll(ctx context.Context, options ScrollOptions) error

		ScrollTop(ctx context.Context, options ScrollOptions) error

		ScrollBottom(ctx context.Context, options ScrollOptions) error

		ScrollBySelector(ctx context.Context, selector QuerySelector, options ScrollOptions) error

		MoveMouseByXY(ctx context.Context, x, y core.Float) error
	}

	// HTMLPage interface represents any web page loaded in the browser
	// and serves as an entry point into the web page's content
	HTMLPage interface {
		core.Value
		core.Iterable
		core.Keyed
		core.Measurable
		events.Observable
		io.Closer

		IsClosed() core.Boolean

		GetURL() core.String

		GetMainFrame() HTMLDocument

		GetFrames(ctx context.Context) (*internal.Array, error)

		GetFrame(ctx context.Context, idx core.Int) (core.Value, error)

		GetCookies(ctx context.Context) (*HTTPCookies, error)

		SetCookies(ctx context.Context, cookies *HTTPCookies) error

		DeleteCookies(ctx context.Context, cookies *HTTPCookies) error

		GetResponse(ctx context.Context) (HTTPResponse, error)

		PrintToPDF(ctx context.Context, params PDFParams) (core.Binary, error)

		CaptureScreenshot(ctx context.Context, params ScreenshotParams) (core.Binary, error)

		WaitForNavigation(ctx context.Context, targetURL core.String) error

		WaitForFrameNavigation(ctx context.Context, frame HTMLDocument, targetURL core.String) error

		Navigate(ctx context.Context, url core.String) error

		NavigateBack(ctx context.Context, skip core.Int) (core.Boolean, error)

		NavigateForward(ctx context.Context, skip core.Int) (core.Boolean, error)
	}
)

const (
	// WaitEventPresence indicating to wait for value to appear
	WaitEventPresence = 0

	// WaitEventAbsence indicating to wait for value to disappear
	WaitEventAbsence = 1
)
