package drivers

import "github.com/MontFerret/ferret/pkg/runtime/values"

// PDFParams represents the arguments for PrintToPDF function.
type PDFParams struct {
	// Paper orientation. Defaults to false.
	Landscape values.Boolean
	// Display values and footer. Defaults to false.
	DisplayHeaderFooter values.Boolean
	// Print background graphics. Defaults to false.
	PrintBackground values.Boolean
	// Scale of the webpage rendering. Defaults to 1.
	Scale values.Float
	// Paper width in inches. Defaults to 8.5 inches.
	PaperWidth values.Float
	// Paper height in inches. Defaults to 11 inches.
	PaperHeight values.Float
	// Top margin in inches. Defaults to 1cm (~0.4 inches).
	MarginTop values.Float
	// Bottom margin in inches. Defaults to 1cm (~0.4 inches).
	MarginBottom values.Float
	// Left margin in inches. Defaults to 1cm (~0.4 inches).
	MarginLeft values.Float
	// Right margin in inches. Defaults to 1cm (~0.4 inches).
	MarginRight values.Float
	// Paper ranges to print, e.g., '1-5, 8, 11-13'. Defaults to the empty string, which means print all pages.
	PageRanges values.String
	// HTML template for the print values. Should be valid HTML markup with following classes used to inject printing values into them: - `date`: formatted print date - `title`: document title - `url`: document location - `pageNumber`: current page number - `totalPages`: total pages in the document
	// For example, `<span class=title></span>` would generate span containing the title.
	HeaderTemplate values.String
	// HTML template for the print footer. Should use the same format as the `headerTemplate`.
	FooterTemplate values.String
	// Whether or not to prefer page size as defined by css.
	// Defaults to false, in which case the content will be scaled to fit the paper size.
	PreferCSSPageSize values.Boolean
}
