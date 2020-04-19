package drivers

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// WaitEvent is an enum that represents what event is needed to wait for
	WaitEvent int

	// Node is an interface from which a number of DOM API object types inherit.
	// It allows those types to be treated similarly;
	// for example, inheriting the same set of methods, or being tested in the same way.
	HTMLNode interface {
		core.Value
		core.Iterable
		core.Getter
		core.Setter
		collections.Measurable
		io.Closer

		IsDetached() values.Boolean

		GetNodeType() values.Int

		GetNodeName() values.String

		GetChildNodes(ctx context.Context) (*values.Array, error)

		GetChildNode(ctx context.Context, idx values.Int) (core.Value, error)

		QuerySelector(ctx context.Context, selector values.String) (core.Value, error)

		QuerySelectorAll(ctx context.Context, selector values.String) (*values.Array, error)

		CountBySelector(ctx context.Context, selector values.String) (values.Int, error)

		ExistsBySelector(ctx context.Context, selector values.String) (values.Boolean, error)

		XPath(ctx context.Context, expression values.String) (core.Value, error)
	}

	// HTMLElement is the most general base interface which most objects in a GetMainFrame implement.
	HTMLElement interface {
		HTMLNode

		GetInnerText(ctx context.Context) (values.String, error)

		SetInnerText(ctx context.Context, innerText values.String) error

		GetInnerHTML(ctx context.Context) (values.String, error)

		SetInnerHTML(ctx context.Context, innerHTML values.String) error

		GetValue(ctx context.Context) (core.Value, error)

		SetValue(ctx context.Context, value core.Value) error

		GetStyles(ctx context.Context) (*values.Object, error)

		GetStyle(ctx context.Context, name values.String) (core.Value, error)

		SetStyles(ctx context.Context, values *values.Object) error

		SetStyle(ctx context.Context, name values.String, value core.Value) error

		RemoveStyle(ctx context.Context, name ...values.String) error

		GetAttributes(ctx context.Context) (*values.Object, error)

		GetAttribute(ctx context.Context, name values.String) (core.Value, error)

		SetAttributes(ctx context.Context, values *values.Object) error

		SetAttribute(ctx context.Context, name, value values.String) error

		RemoveAttribute(ctx context.Context, name ...values.String) error

		GetInnerHTMLBySelector(ctx context.Context, selector values.String) (values.String, error)

		SetInnerHTMLBySelector(ctx context.Context, selector, innerHTML values.String) error

		GetInnerHTMLBySelectorAll(ctx context.Context, selector values.String) (*values.Array, error)

		GetInnerTextBySelector(ctx context.Context, selector values.String) (values.String, error)

		SetInnerTextBySelector(ctx context.Context, selector, innerText values.String) error

		GetInnerTextBySelectorAll(ctx context.Context, selector values.String) (*values.Array, error)

		Click(ctx context.Context, count values.Int) error

		ClickBySelector(ctx context.Context, selector values.String, count values.Int) error

		ClickBySelectorAll(ctx context.Context, selector values.String, count values.Int) error

		Clear(ctx context.Context) error

		ClearBySelector(ctx context.Context, selector values.String) error

		Input(ctx context.Context, value core.Value, delay values.Int) error

		InputBySelector(ctx context.Context, selector values.String, value core.Value, delay values.Int) error

		Select(ctx context.Context, value *values.Array) (*values.Array, error)

		SelectBySelector(ctx context.Context, selector values.String, value *values.Array) (*values.Array, error)

		ScrollIntoView(ctx context.Context, options ScrollOptions) error

		Focus(ctx context.Context) error

		FocusBySelector(ctx context.Context, selector values.String) error

		Blur(ctx context.Context) error

		BlurBySelector(ctx context.Context, selector values.String) error

		Hover(ctx context.Context) error

		HoverBySelector(ctx context.Context, selector values.String) error

		WaitForAttribute(ctx context.Context, name values.String, value core.Value, when WaitEvent) error

		WaitForStyle(ctx context.Context, name values.String, value core.Value, when WaitEvent) error

		WaitForClass(ctx context.Context, class values.String, when WaitEvent) error
	}

	HTMLDocument interface {
		HTMLNode

		GetTitle() values.String

		GetElement() HTMLElement

		GetURL() values.String

		GetName() values.String

		GetParentDocument(ctx context.Context) (HTMLDocument, error)

		GetChildDocuments(ctx context.Context) (*values.Array, error)

		ScrollTop(ctx context.Context, options ScrollOptions) error

		ScrollBottom(ctx context.Context, options ScrollOptions) error

		ScrollBySelector(ctx context.Context, selector values.String, options ScrollOptions) error

		ScrollByXY(ctx context.Context, x, y values.Float, options ScrollOptions) error

		MoveMouseByXY(ctx context.Context, x, y values.Float) error

		WaitForElement(ctx context.Context, selector values.String, when WaitEvent) error

		WaitForAttributeBySelector(ctx context.Context, selector, name values.String, value core.Value, when WaitEvent) error

		WaitForAttributeBySelectorAll(ctx context.Context, selector, name values.String, value core.Value, when WaitEvent) error

		WaitForStyleBySelector(ctx context.Context, selector, name values.String, value core.Value, when WaitEvent) error

		WaitForStyleBySelectorAll(ctx context.Context, selector, name values.String, value core.Value, when WaitEvent) error

		WaitForClassBySelector(ctx context.Context, selector, class values.String, when WaitEvent) error

		WaitForClassBySelectorAll(ctx context.Context, selector, class values.String, when WaitEvent) error
	}

	// HTMLPage interface represents any web page loaded in the browser
	// and serves as an entry point into the web page's content
	HTMLPage interface {
		core.Value
		core.Iterable
		core.Getter
		core.Setter
		collections.Measurable
		io.Closer

		IsClosed() values.Boolean

		GetURL() values.String

		GetMainFrame() HTMLDocument

		GetFrames(ctx context.Context) (*values.Array, error)

		GetFrame(ctx context.Context, idx values.Int) (core.Value, error)

		GetCookies(ctx context.Context) (HTTPCookies, error)

		SetCookies(ctx context.Context, cookies HTTPCookies) error

		DeleteCookies(ctx context.Context, cookies HTTPCookies) error

		GetResponse(ctx context.Context) (HTTPResponse, error)

		PrintToPDF(ctx context.Context, params PDFParams) (values.Binary, error)

		CaptureScreenshot(ctx context.Context, params ScreenshotParams) (values.Binary, error)

		WaitForNavigation(ctx context.Context, targetURL values.String) error

		Navigate(ctx context.Context, url values.String) error

		NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error)

		NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error)
	}
)

const (
	// Event indicating to wait for value to appear
	WaitEventPresence = 0

	// Event indicating to wait for value to disappear
	WaitEventAbsence = 1
)
