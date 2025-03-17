package html

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"regexp"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func ValidatePageRanges(pageRanges string) (bool, error) {
	match, err := regexp.Match(`^(([1-9][0-9]*|[1-9][0-9]*)(\s*-\s*|\s*,\s*|))*$`, []byte(pageRanges))

	if err != nil {
		return false, err
	}

	return match, nil
}

// PDF prints a PDF of the current page.
// @param {HTMLPage | String}target - Target page or url.
// @param {hashMap} [params] - An object containing the following properties:
// @param {Bool} [params.landscape=False] - Paper orientation.
// @param {Bool} [params.displayHeaderFooter=False] - Display header and footer.
// @param {Bool} [params.printBackground=False] - Print background graphics.
// @param {Float} [params.scale=1] - Scale of the webpage rendering.
// @param {Float} [params.paperWidth=22] - Paper width in inches.
// @param {Float} [params.paperHeight=28] - Paper height in inches.
// @param {Float} [params.marginTo=1] - Top margin in inches.
// @param {Float} [params.marginBottom=1] - Bottom margin in inches.
// @param {Float} [params.marginLeft=1] - Left margin in inches.
// @param {Float} [params.marginRight=1] - Right margin in inches.
// @param {String} [params.pageRanges] - Paper ranges to print, e.g., '1-5, 8, 11-13'.
// @param {String} [params.headerTemplate] - HTML template for the print header. Should be valid HTML markup with following classes used to inject printing values into them: - `date`: formatted print date - `title`: document title - `url`: document location - `pageNumber`: current page number - `totalPages`: total pages in the document For example, `<span class=title></span>` would generate span containing the title.
// @param {String} [params.footerTemplate] - HTML template for the print footer. Should use the same format as the `headerTemplate`.
// @param {Bool} [params.preferCSSPageSize=False] - Whether or not to prefer page size as defined by css. Defaults to false, in which case the content will be scaled to fit the paper size. *
// @return {Binary} - PDF document in binary format.
func PDF(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return core.None, err
	}

	arg1 := args[0]
	page, closeAfter, err := OpenOrCastPage(ctx, arg1)

	if err != nil {
		return core.None, err
	}

	defer func() {
		if closeAfter {
			page.Close()
		}
	}()

	pdfParams := drivers.PDFParams{}

	if len(args) == 2 {
		arg2 := args[1]
		err = core.ValidateType(arg2, types.Object)

		if err != nil {
			return core.None, err
		}

		params, ok := arg2.(*internal.Object)

		if !ok {
			return core.None, core.Error(core.ErrInvalidType, "expected object")
		}

		landscape, found := params.Get("landscape")

		if found {
			err = core.ValidateType(landscape, types.Boolean)

			if err != nil {
				return core.None, err
			}

			pdfParams.Landscape = landscape.(core.Boolean)
		}

		displayHeaderFooter, found := params.Get("displayHeaderFooter")

		if found {
			err = core.ValidateType(displayHeaderFooter, types.Boolean)

			if err != nil {
				return core.None, err
			}

			pdfParams.DisplayHeaderFooter = displayHeaderFooter.(core.Boolean)
		}

		printBackground, found := params.Get("printBackground")

		if found {
			err = core.ValidateType(printBackground, types.Boolean)

			if err != nil {
				return core.None, err
			}

			pdfParams.PrintBackground = printBackground.(core.Boolean)
		}

		scale, found := params.Get("scale")

		if found {
			err = core.ValidateType(scale, types.Float, types.Int)

			if err != nil {
				return core.None, err
			}

			if scale.Type() == types.Int {
				pdfParams.Scale = core.Float(scale.(core.Int))
			} else {
				pdfParams.Scale = scale.(core.Float)
			}
		}

		paperWidth, found := params.Get("paperWidth")

		if found {
			err = core.ValidateType(paperWidth, types.Float, types.Int)

			if err != nil {
				return core.None, err
			}

			if paperWidth.Type() == types.Int {
				pdfParams.PaperWidth = core.Float(paperWidth.(core.Int))
			} else {
				pdfParams.PaperWidth = paperWidth.(core.Float)
			}
		}

		paperHeight, found := params.Get("paperHeight")

		if found {
			err = core.ValidateType(paperHeight, types.Float, types.Int)

			if err != nil {
				return core.None, err
			}

			if paperHeight.Type() == types.Int {
				pdfParams.PaperHeight = core.Float(paperHeight.(core.Int))
			} else {
				pdfParams.PaperHeight = paperHeight.(core.Float)
			}
		}

		marginTop, found := params.Get("marginTop")

		if found {
			err = core.ValidateType(marginTop, types.Float, types.Int)

			if err != nil {
				return core.None, err
			}

			if marginTop.Type() == types.Int {
				pdfParams.MarginTop = core.Float(marginTop.(core.Int))
			} else {
				pdfParams.MarginTop = marginTop.(core.Float)
			}
		}

		marginBottom, found := params.Get("marginBottom")

		if found {
			err = core.ValidateType(marginBottom, types.Float, types.Int)

			if err != nil {
				return core.None, err
			}

			if marginBottom.Type() == types.Int {
				pdfParams.MarginBottom = core.Float(marginBottom.(core.Int))
			} else {
				pdfParams.MarginBottom = marginBottom.(core.Float)
			}
		}

		marginLeft, found := params.Get("marginLeft")

		if found {
			err = core.ValidateType(marginLeft, types.Float, types.Int)

			if err != nil {
				return core.None, err
			}

			if marginLeft.Type() == types.Int {
				pdfParams.MarginLeft = core.Float(marginLeft.(core.Int))
			} else {
				pdfParams.MarginLeft = marginLeft.(core.Float)
			}
		}

		marginRight, found := params.Get("marginRight")

		if found {
			err = core.ValidateType(marginRight, types.Float, types.Int)

			if err != nil {
				return core.None, err
			}

			if marginRight.Type() == types.Int {
				pdfParams.MarginRight = core.Float(marginRight.(core.Int))
			} else {
				pdfParams.MarginRight = marginRight.(core.Float)
			}
		}

		pageRanges, found := params.Get("pageRanges")

		if found {
			err = core.ValidateType(pageRanges, types.String)

			if err != nil {
				return core.None, err
			}

			validate, err := ValidatePageRanges(pageRanges.String())

			if err != nil {
				return core.None, err
			}

			if !validate {
				return core.None, core.Error(core.ErrInvalidArgument, fmt.Sprintf(`page ranges "%s", not valid`, pageRanges.String()))
			}

			pdfParams.PageRanges = pageRanges.(core.String)
		}

		headerTemplate, found := params.Get("headerTemplate")

		if found {
			err = core.ValidateType(headerTemplate, types.String)

			if err != nil {
				return core.None, err
			}

			pdfParams.HeaderTemplate = headerTemplate.(core.String)
		}

		footerTemplate, found := params.Get("footerTemplate")

		if found {
			err = core.ValidateType(footerTemplate, types.String)

			if err != nil {
				return core.None, err
			}

			pdfParams.FooterTemplate = footerTemplate.(core.String)
		}

		preferCSSPageSize, found := params.Get("preferCSSPageSize")

		if found {
			err = core.ValidateType(preferCSSPageSize, types.Boolean)

			if err != nil {
				return core.None, err
			}

			pdfParams.PreferCSSPageSize = preferCSSPageSize.(core.Boolean)
		}
	}

	pdf, err := page.PrintToPDF(ctx, pdfParams)

	if err != nil {
		return core.None, err
	}

	return pdf, nil
}
