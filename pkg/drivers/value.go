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

		NodeType() values.Int

		NodeName() values.String

		GetChildNodes(ctx context.Context) core.Value

		GetChildNode(ctx context.Context, idx values.Int) core.Value

		QuerySelector(ctx context.Context, selector values.String) core.Value

		QuerySelectorAll(ctx context.Context, selector values.String) core.Value

		CountBySelector(ctx context.Context, selector values.String) values.Int

		ExistsBySelector(ctx context.Context, selector values.String) values.Boolean
	}

	// HTMLElement is the most general base interface which most objects in a Document implement.
	HTMLElement interface {
		HTMLNode

		InnerText(ctx context.Context) values.String

		InnerHTML(ctx context.Context) values.String

		GetValue(ctx context.Context) core.Value

		SetValue(ctx context.Context, value core.Value) error

		GetStyles(ctx context.Context) (*values.Object, error)

		GetStyle(ctx context.Context, name values.String) (core.Value, error)

		SetStyles(ctx context.Context, values *values.Object) error

		SetStyle(ctx context.Context, name values.String, value core.Value) error

		RemoveStyle(ctx context.Context, name ...values.String) error

		GetAttributes(ctx context.Context) *values.Object

		GetAttribute(ctx context.Context, name values.String) core.Value

		SetAttributes(ctx context.Context, values *values.Object) error

		SetAttribute(ctx context.Context, name, value values.String) error

		RemoveAttribute(ctx context.Context, name ...values.String) error

		InnerHTMLBySelector(ctx context.Context, selector values.String) values.String

		InnerHTMLBySelectorAll(ctx context.Context, selector values.String) *values.Array

		InnerTextBySelector(ctx context.Context, selector values.String) values.String

		InnerTextBySelectorAll(ctx context.Context, selector values.String) *values.Array

		Click(ctx context.Context) (values.Boolean, error)

		Input(ctx context.Context, value core.Value, delay values.Int) error

		Select(ctx context.Context, value *values.Array) (*values.Array, error)

		ScrollIntoView(ctx context.Context) error

		Hover(ctx context.Context) error

		WaitForAttribute(ctx context.Context, name values.String, value core.Value, when WaitEvent) error

		WaitForStyle(ctx context.Context, name values.String, value core.Value, when WaitEvent) error

		WaitForClass(ctx context.Context, class values.String, when WaitEvent) error
	}

	HTMLDocument interface {
		HTMLNode

		Element() HTMLElement

		GetURL() core.Value

		Name() values.String

		GetParentDocument() HTMLDocument

		GetChildDocuments(ctx context.Context) (*values.Array, error)

		ClickBySelector(ctx context.Context, selector values.String) (values.Boolean, error)

		ClickBySelectorAll(ctx context.Context, selector values.String) (values.Boolean, error)

		InputBySelector(ctx context.Context, selector values.String, value core.Value, delay values.Int) (values.Boolean, error)

		SelectBySelector(ctx context.Context, selector values.String, value *values.Array) (*values.Array, error)

		ScrollTop(ctx context.Context) error

		ScrollBottom(ctx context.Context) error

		ScrollBySelector(ctx context.Context, selector values.String) error

		ScrollByXY(ctx context.Context, x, y values.Float) error

		MoveMouseByXY(ctx context.Context, x, y values.Float) error

		MoveMouseBySelector(ctx context.Context, selector values.String) error

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

		Document() HTMLDocument

		GetCookies(ctx context.Context) (*values.Array, error)

		SetCookies(ctx context.Context, cookies ...HTTPCookie) error

		DeleteCookies(ctx context.Context, cookies ...HTTPCookie) error

		PrintToPDF(ctx context.Context, params PDFParams) (values.Binary, error)

		CaptureScreenshot(ctx context.Context, params ScreenshotParams) (values.Binary, error)

		WaitForNavigation(ctx context.Context) error

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
