package drivers

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
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

		GetAttributes(ctx context.Context) core.Value

		GetAttribute(ctx context.Context, name values.String) core.Value

		SetAttribute(ctx context.Context, name, value values.String) error

		InnerHTMLBySelector(ctx context.Context, selector values.String) values.String

		InnerHTMLBySelectorAll(ctx context.Context, selector values.String) *values.Array

		InnerTextBySelector(ctx context.Context, selector values.String) values.String

		InnerTextBySelectorAll(ctx context.Context, selector values.String) *values.Array

		Click(ctx context.Context) (values.Boolean, error)

		Input(ctx context.Context, value core.Value, delay values.Int) error

		Select(ctx context.Context, value *values.Array) (*values.Array, error)

		ScrollIntoView(ctx context.Context) error

		Hover(ctx context.Context) error

		WaitForClass(ctx context.Context, class values.String) error
	}

	// The Document interface represents any web page loaded in the browser
	// and serves as an entry point into the web page's content, which is the DOM tree.
	HTMLDocument interface {
		HTMLNode

		DocumentElement() HTMLElement

		GetURL() core.Value

		SetURL(ctx context.Context, url values.String) error

		Navigate(ctx context.Context, url values.String) error

		NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error)

		NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error)

		ClickBySelector(ctx context.Context, selector values.String) (values.Boolean, error)

		ClickBySelectorAll(ctx context.Context, selector values.String) (values.Boolean, error)

		InputBySelector(ctx context.Context, selector values.String, value core.Value, delay values.Int) (values.Boolean, error)

		SelectBySelector(ctx context.Context, selector values.String, value *values.Array) (*values.Array, error)

		HoverBySelector(ctx context.Context, selector values.String) error

		PrintToPDF(ctx context.Context, params PDFParams) (values.Binary, error)

		CaptureScreenshot(ctx context.Context, params ScreenshotParams) (values.Binary, error)

		ScrollTop(ctx context.Context) error

		ScrollBottom(ctx context.Context) error

		ScrollBySelector(ctx context.Context, selector values.String) error

		WaitForNavigation(ctx context.Context) error

		WaitForSelector(ctx context.Context, selector values.String) error

		WaitForClassBySelector(ctx context.Context, selector, class values.String) error

		WaitForClassBySelectorAll(ctx context.Context, selector, class values.String) error
	}
)
