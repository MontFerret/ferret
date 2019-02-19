package drivers

import (
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

		GetChildNodes() core.Value

		GetChildNode(idx values.Int) core.Value

		QuerySelector(selector values.String) core.Value

		QuerySelectorAll(selector values.String) core.Value

		CountBySelector(selector values.String) values.Int

		ExistsBySelector(selector values.String) values.Boolean
	}

	// HTMLElement is the most general base interface which most objects in a Document implement.
	HTMLElement interface {
		HTMLNode

		InnerText() values.String

		InnerHTML() values.String

		GetValue() core.Value

		SetValue(value core.Value) error

		GetAttributes() core.Value

		GetAttribute(name values.String) core.Value

		SetAttribute(name, value values.String) error

		InnerHTMLBySelector(selector values.String) values.String

		InnerHTMLBySelectorAll(selector values.String) *values.Array

		InnerTextBySelector(selector values.String) values.String

		InnerTextBySelectorAll(selector values.String) *values.Array

		Click() (values.Boolean, error)

		Input(value core.Value, delay values.Int) error

		Select(value *values.Array) (*values.Array, error)

		ScrollIntoView() error

		Hover() error

		WaitForClass(class values.String, timeout values.Int) error
	}

	// The Document interface represents any web page loaded in the browser
	// and serves as an entry point into the web page's content, which is the DOM tree.
	HTMLDocument interface {
		HTMLNode

		DocumentElement() HTMLElement

		GetURL() core.Value

		SetURL(url values.String) error

		Navigate(url values.String, timeout values.Int) error

		NavigateBack(skip values.Int, timeout values.Int) (values.Boolean, error)

		NavigateForward(skip values.Int, timeout values.Int) (values.Boolean, error)

		ClickBySelector(selector values.String) (values.Boolean, error)

		ClickBySelectorAll(selector values.String) (values.Boolean, error)

		InputBySelector(selector values.String, value core.Value, delay values.Int) (values.Boolean, error)

		SelectBySelector(selector values.String, value *values.Array) (*values.Array, error)

		HoverBySelector(selector values.String) error

		PrintToPDF(params PDFParams) (values.Binary, error)

		CaptureScreenshot(params ScreenshotParams) (values.Binary, error)

		ScrollTop() error

		ScrollBottom() error

		ScrollBySelector(selector values.String) error

		WaitForNavigation(timeout values.Int) error

		WaitForSelector(selector values.String, timeout values.Int) error

		WaitForClassBySelector(selector, class values.String, timeout values.Int) error

		WaitForClassBySelectorAll(selector, class values.String, timeout values.Int) error
	}
)
