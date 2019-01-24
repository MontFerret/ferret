package html

import (
	"context"
	"fmt"
	"regexp"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func ValidatePageRanges(pageRanges string) (bool, error) {
	match, err := regexp.Match(`^(([1-9][0-9]*|[1-9][0-9]*)(\s*-\s*|\s*,\s*|))*$`, []byte(pageRanges))

	if err != nil {
		return false, err
	}

	return match, nil
}

// PDF print a PDF of the current page.
// @param source (Document) - Document.
// @param params (Object) - Optional, An object containing the following properties :
//   Landscape (Bool) - Paper orientation. Defaults to false.
//   DisplayHeaderFooter (Bool) - Display header and footer. Defaults to false.
//   PrintBackground (Bool) - Print background graphics. Defaults to false.
//   Scale (Float64) - Scale of the webpage rendering. Defaults to 1.
//   PaperWidth (Float64) - Paper width in inches. Defaults to 8.5 inches.
//   PaperHeight (Float64) - Paper height in inches. Defaults to 11 inches.
//   MarginTop (Float64) - Top margin in inches. Defaults to 1cm (~0.4 inches).
//   MarginBottom (Float64) - Bottom margin in inches. Defaults to 1cm (~0.4 inches).
//   MarginLeft (Float64) - Left margin in inches. Defaults to 1cm (~0.4 inches).
//   MarginRight (Float64) - Right margin in inches. Defaults to 1cm (~0.4 inches).
//   PageRanges (String) - Paper ranges to print, e.g., '1-5, 8, 11-13'. Defaults to the empty string, which means print all pages.
//   IgnoreInvalidPageRanges (Bool) - to silently ignore invalid but successfully parsed page ranges, such as '3-2'. Defaults to false.
//   HeaderTemplate (String) - HTML template for the print header. Should be valid HTML markup with following classes used to inject printing values into them: - `date`: formatted print date - `title`: document title - `url`: document location - `pageNumber`: current page number - `totalPages`: total pages in the document For example, `<span class=title></span>` would generate span containing the title.
//   FooterTemplate (String) - HTML template for the print footer. Should use the same format as the `headerTemplate`.
//   PreferCSSPageSize (Bool) - Whether or not to prefer page size as defined by css. Defaults to false, in which case the content will be scaled to fit the paper size. *
// @returns data (Binary) - Returns a base64 encoded string in binary format.
func PDF(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	arg1 := args[0]
	val, err := ValidateDocument(ctx, arg1)

	if err != nil {
		return values.None, err
	}

	doc := val.(drivers.DHTMLDocument)
	defer doc.Close()

	pdfParams := drivers.PDFParams{}

	if len(args) == 2 {
		arg2 := args[1]
		err = core.ValidateType(arg2, core.ObjectType)

		if err != nil {
			return values.None, err
		}

		params, ok := arg2.(*values.Object)

		if !ok {
			return values.None, core.Error(core.ErrInvalidType, "expected object")
		}

		landscape, found := params.Get("landscape")

		if found {
			err = core.ValidateType(landscape, core.BooleanType)

			if err != nil {
				return values.None, err
			}

			pdfParams.Landscape = landscape.(values.Boolean)
		}

		displayHeaderFooter, found := params.Get("displayHeaderFooter")

		if found {
			err = core.ValidateType(displayHeaderFooter, core.BooleanType)

			if err != nil {
				return values.None, err
			}

			pdfParams.DisplayHeaderFooter = displayHeaderFooter.(values.Boolean)
		}

		printBackground, found := params.Get("printBackground")

		if found {
			err = core.ValidateType(printBackground, core.BooleanType)

			if err != nil {
				return values.None, err
			}

			pdfParams.PrintBackground = printBackground.(values.Boolean)
		}

		scale, found := params.Get("scale")

		if found {
			err = core.ValidateType(scale, core.FloatType, core.IntType)

			if err != nil {
				return values.None, err
			}

			if scale.Type() == core.IntType {
				pdfParams.Scale = values.Float(scale.(values.Int))
			} else {
				pdfParams.Scale = scale.(values.Float)
			}
		}

		paperWidth, found := params.Get("paperWidth")

		if found {
			err = core.ValidateType(paperWidth, core.FloatType, core.IntType)

			if err != nil {
				return values.None, err
			}

			if paperWidth.Type() == core.IntType {
				pdfParams.PaperWidth = values.Float(paperWidth.(values.Int))
			} else {
				pdfParams.PaperWidth = paperWidth.(values.Float)
			}
		}

		paperHeight, found := params.Get("paperHeight")

		if found {
			err = core.ValidateType(paperHeight, core.FloatType, core.IntType)

			if err != nil {
				return values.None, err
			}

			if paperHeight.Type() == core.IntType {
				pdfParams.PaperHeight = values.Float(paperHeight.(values.Int))
			} else {
				pdfParams.PaperHeight = paperHeight.(values.Float)
			}
		}

		marginTop, found := params.Get("marginTop")

		if found {
			err = core.ValidateType(marginTop, core.FloatType, core.IntType)

			if err != nil {
				return values.None, err
			}

			if marginTop.Type() == core.IntType {
				pdfParams.MarginTop = values.Float(marginTop.(values.Int))
			} else {
				pdfParams.MarginTop = marginTop.(values.Float)
			}
		}

		marginBottom, found := params.Get("marginBottom")

		if found {
			err = core.ValidateType(marginBottom, core.FloatType, core.IntType)

			if err != nil {
				return values.None, err
			}

			if marginBottom.Type() == core.IntType {
				pdfParams.MarginBottom = values.Float(marginBottom.(values.Int))
			} else {
				pdfParams.MarginBottom = marginBottom.(values.Float)
			}
		}

		marginLeft, found := params.Get("marginLeft")

		if found {
			err = core.ValidateType(marginLeft, core.FloatType, core.IntType)

			if err != nil {
				return values.None, err
			}

			if marginLeft.Type() == core.IntType {
				pdfParams.MarginLeft = values.Float(marginLeft.(values.Int))
			} else {
				pdfParams.MarginLeft = marginLeft.(values.Float)
			}
		}

		marginRight, found := params.Get("marginRight")

		if found {
			err = core.ValidateType(marginRight, core.FloatType, core.IntType)

			if err != nil {
				return values.None, err
			}

			if marginRight.Type() == core.IntType {
				pdfParams.MarginRight = values.Float(marginRight.(values.Int))
			} else {
				pdfParams.MarginRight = marginRight.(values.Float)
			}
		}

		pageRanges, found := params.Get("pageRanges")

		if found {
			err = core.ValidateType(pageRanges, core.StringType)

			if err != nil {
				return values.None, err
			}

			validate, err := ValidatePageRanges(pageRanges.String())

			if err != nil {
				return values.None, err
			}

			if !validate {
				return values.None, core.Error(core.ErrInvalidArgument, fmt.Sprintf(`page ranges "%s", not valid`, pageRanges.String()))
			}

			pdfParams.PageRanges = pageRanges.(values.String)
		}

		ignoreInvalidPageRanges, found := params.Get("ignoreInvalidPageRanges")

		if found {
			err = core.ValidateType(ignoreInvalidPageRanges, core.BooleanType)

			if err != nil {
				return values.None, err
			}

			pdfParams.IgnoreInvalidPageRanges = ignoreInvalidPageRanges.(values.Boolean)
		}

		headerTemplate, found := params.Get("headerTemplate")

		if found {
			err = core.ValidateType(headerTemplate, core.StringType)

			if err != nil {
				return values.None, err
			}

			pdfParams.HeaderTemplate = headerTemplate.(values.String)
		}

		footerTemplate, found := params.Get("footerTemplate")

		if found {
			err = core.ValidateType(footerTemplate, core.StringType)

			if err != nil {
				return values.None, err
			}

			pdfParams.FooterTemplate = footerTemplate.(values.String)
		}

		preferCSSPageSize, found := params.Get("preferCSSPageSize")

		if found {
			err = core.ValidateType(preferCSSPageSize, core.BooleanType)

			if err != nil {
				return values.None, err
			}

			pdfParams.PreferCSSPageSize = preferCSSPageSize.(values.Boolean)
		}
	}

	pdf, err := doc.PrintToPDF(pdfParams)

	if err != nil {
		return values.None, err
	}

	return pdf, nil
}
