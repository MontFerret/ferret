package values

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"io"
)

const (
	HTMLScreenshotFormatPNG  HTMLScreenshotFormat = "png"
	HTMLScreenshotFormatJPEG HTMLScreenshotFormat = "jpeg"
)

type (
	// HTMLPDFParams represents the arguments for PrintToPDF function.
	HTMLPDFParams struct {
		// Paper orientation. Defaults to false.
		Landscape Boolean
		// Display header and footer. Defaults to false.
		DisplayHeaderFooter Boolean
		// Print background graphics. Defaults to false.
		PrintBackground Boolean
		// Scale of the webpage rendering. Defaults to 1.
		Scale Float
		// Paper width in inches. Defaults to 8.5 inches.
		PaperWidth Float
		// Paper height in inches. Defaults to 11 inches.
		PaperHeight Float
		// Top margin in inches. Defaults to 1cm (~0.4 inches).
		MarginTop Float
		// Bottom margin in inches. Defaults to 1cm (~0.4 inches).
		MarginBottom Float
		// Left margin in inches. Defaults to 1cm (~0.4 inches).
		MarginLeft Float
		// Right margin in inches. Defaults to 1cm (~0.4 inches).
		MarginRight Float
		// Paper ranges to print, e.g., '1-5, 8, 11-13'. Defaults to the empty string, which means print all pages.
		PageRanges String
		// Whether to silently ignore invalid but successfully parsed page ranges, such as '3-2'. Defaults to false.
		IgnoreInvalidPageRanges Boolean
		// HTML template for the print header. Should be valid HTML markup with following classes used to inject printing values into them: - `date`: formatted print date - `title`: document title - `url`: document location - `pageNumber`: current page number - `totalPages`: total pages in the document
		// For example, `<span class=title></span>` would generate span containing the title.
		HeaderTemplate String
		// HTML template for the print footer. Should use the same format as the `headerTemplate`.
		FooterTemplate String
		// Whether or not to prefer page size as defined by css.
		// Defaults to false, in which case the content will be scaled to fit the paper size.
		PreferCSSPageSize Boolean
	}

	HTMLScreenshotFormat string

	HTMLScreenshotParams struct {
		X       Float
		Y       Float
		Width   Float
		Height  Float
		Format  HTMLScreenshotFormat
		Quality Int
	}

	// HTMLNode is a HTML Node
	HTMLNode interface {
		core.Value

		NodeType() Int

		NodeName() String

		Length() Int

		InnerText() String

		InnerHTML() String

		Value() core.Value

		GetAttributes() core.Value

		GetAttribute(name String) core.Value

		GetChildNodes() core.Value

		GetChildNode(idx Int) core.Value

		QuerySelector(selector String) core.Value

		QuerySelectorAll(selector String) core.Value

		InnerHTMLBySelector(selector String) String

		InnerHTMLBySelectorAll(selector String) *Array

		InnerTextBySelector(selector String) String

		InnerTextBySelectorAll(selector String) *Array

		CountBySelector(selector String) Int

		ExistsBySelector(selector String) Boolean
	}

	DHTMLNode interface {
		HTMLNode
		io.Closer

		Click() (Boolean, error)

		Input(value core.Value, delay Int) error

		Select(value *Array) (*Array, error)

		ScrollIntoView() error

		Hover() error

		WaitForClass(class String, timeout Int) error
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

		Navigate(url String, timeout Int) error

		NavigateBack(skip Int, timeout Int) (Boolean, error)

		NavigateForward(skip Int, timeout Int) (Boolean, error)

		ClickBySelector(selector String) (Boolean, error)

		ClickBySelectorAll(selector String) (Boolean, error)

		InputBySelector(selector String, value core.Value, delay Int) (Boolean, error)

		SelectBySelector(selector String, value *Array) (*Array, error)

		HoverBySelector(selector String) error

		PrintToPDF(params HTMLPDFParams) (Binary, error)

		CaptureScreenshot(params HTMLScreenshotParams) (Binary, error)

		ScrollTop() error

		ScrollBottom() error

		ScrollBySelector(selector String) error

		WaitForNavigation(timeout Int) error

		WaitForSelector(selector String, timeout Int) error

		WaitForClass(selector, class String, timeout Int) error

		WaitForClassAll(selector, class String, timeout Int) error
	}
)

func IsHTMLScreenshotFormatValid(format string) bool {
	value := HTMLScreenshotFormat(format)

	return value == HTMLScreenshotFormatPNG || value == HTMLScreenshotFormatJPEG
}

func NewDefaultHTMLPDFParams() HTMLPDFParams {
	return HTMLPDFParams{
		Landscape:               False,
		DisplayHeaderFooter:     False,
		PrintBackground:         False,
		Scale:                   Float(1),
		PaperWidth:              Float(8.5),
		PaperHeight:             Float(11),
		MarginTop:               Float(0.4),
		MarginBottom:            Float(0.4),
		MarginLeft:              Float(0.4),
		MarginRight:             Float(0.4),
		PageRanges:              EmptyString,
		IgnoreInvalidPageRanges: False,
		HeaderTemplate:          EmptyString,
		FooterTemplate:          EmptyString,
		PreferCSSPageSize:       False,
	}
}
