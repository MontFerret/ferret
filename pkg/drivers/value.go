package drivers

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"io"
)

type (
	// HTMLNode is a HTML Node
	HTMLNode interface {
		collections.Collection
		collections.IterableCollection
		core.Getter
		core.Setter

		NodeType() values.Int

		NodeName() values.String

		InnerText() values.String

		InnerHTML() values.String

		Value() core.Value

		GetAttributes() core.Value

		GetAttribute(name values.String) core.Value

		GetChildNodes() core.Value

		GetChildNode(idx values.Int) core.Value

		QuerySelector(selector values.String) core.Value

		QuerySelectorAll(selector values.String) core.Value

		InnerHTMLBySelector(selector values.String) values.String

		InnerHTMLBySelectorAll(selector values.String) *values.Array

		InnerTextBySelector(selector values.String) values.String

		InnerTextBySelectorAll(selector values.String) *values.Array

		CountBySelector(selector values.String) values.Int

		ExistsBySelector(selector values.String) values.Boolean
	}

	DHTMLNode interface {
		HTMLNode
		io.Closer

		Click() (values.Boolean, error)

		Input(value core.Value, delay values.Int) error

		Select(value *values.Array) (*values.Array, error)

		ScrollIntoView() error

		Hover() error

		WaitForClass(class values.String, timeout values.Int) error
	}

	// HTMLDocument is a HTML Document
	HTMLDocument interface {
		HTMLNode

		URL() core.Value
	}

	// DHTMLDocument is a Dynamic HTML Document
	DHTMLDocument interface {
		HTMLDocument
		io.Closer

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

		WaitForClass(selector, class values.String, timeout values.Int) error

		WaitForClassAll(selector, class values.String, timeout values.Int) error
	}
)
