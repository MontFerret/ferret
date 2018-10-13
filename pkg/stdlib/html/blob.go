package html

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mafredri/cdp/protocol/page"

	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func ValidateDocument(ctx context.Context, value core.Value) (core.Value, error) {
	err := core.ValidateType(value, core.HTMLDocumentType, core.StringType)
	if err != nil {
		return values.None, err
	}

	var doc *dynamic.HTMLDocument
	var ok bool
	if value.Type() == core.StringType {
		buf, err := Document(ctx, value, values.NewBoolean(true))
		if err != nil {
			return values.None, err
		}
		doc, ok = buf.(*dynamic.HTMLDocument)
	} else {
		doc, ok = value.(*dynamic.HTMLDocument)
	}

	if !ok {
		return nil, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return doc, nil
}

/*
 * Take a screenshot of the current page.
 * @param source (Document) - Document.
 * @param params (Object) - Optional, An object containing the following properties :
 * 		x (Float|Int) - Optional, X position of the viewport.
 * 		x (Float|Int) - Optional,Y position of the viewport.
 * 		width (Float|Int) - Optional, Width of the viewport.
 * 		height (Float|Int) - Optional, Height of the viewport.
 * 		format (String) - Optional, Either "jpeg" or "png".
 * 		quality (Int) - Optional, Quality, in [0, 100], only for jpeg format.
 * @returns data (Binary) - Returns a base64 encoded string in binary format.
 */
func Screenshot(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)
	if err != nil {
		return values.None, err
	}

	arg1 := args[0]

	err = core.ValidateType(arg1, core.HTMLDocumentType, core.StringType)
	if err != nil {
		return values.None, err
	}

	val, err := ValidateDocument(ctx, arg1)
	if err != nil {
		return values.None, err
	}
	doc := val.(*dynamic.HTMLDocument)
	defer doc.Close()

	screenshotParams := &dynamic.ScreenshotArgs{
		X:       0,
		Y:       0,
		Width:   -1,
		Height:  -1,
		Format:  "jpeg",
		Quality: 100,
	}
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

		format, found := params.Get("format")
		if found {
			err = core.ValidateType(format, core.StringType)
			if err != nil {
				return values.None, err
			}
			if !dynamic.IsScreenshotFormatValid(format.String()) {
				return values.None, core.Error(
					core.ErrInvalidArgument,
					fmt.Sprintf("format is not valid, expected jpeg or png, but got %s", format.String()))
			}
			screenshotParams.Format = dynamic.ScreenshotFormat(format.String())
		}
		x, found := params.Get("x")
		if found {
			err = core.ValidateType(x, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if x.Type() == core.IntType {
				x = values.Float(x.(values.Int))
			}
			screenshotParams.X = x.Unwrap().(float64)
		}
		y, found := params.Get("y")
		if found {
			err = core.ValidateType(y, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if y.Type() == core.IntType {
				y = values.Float(y.(values.Int))
			}
			screenshotParams.Y = y.Unwrap().(float64)
		}
		width, found := params.Get("width")
		if found {
			err = core.ValidateType(width, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if width.Type() == core.IntType {
				width = values.Float(width.(values.Int))
			}
			screenshotParams.Width = width.Unwrap().(float64)
		}
		height, found := params.Get("height")
		if found {
			err = core.ValidateType(height, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if height.Type() == core.IntType {
				height = values.Float(height.(values.Int))
			}
			screenshotParams.Height = height.Unwrap().(float64)
		}
		quality, found := params.Get("quality")
		if found {
			err = core.ValidateType(quality, core.IntType)
			if err != nil {
				return values.None, err
			}
			screenshotParams.Quality = quality.Unwrap().(int)
		}
	}

	scr, err := doc.CaptureScreenshot(screenshotParams)
	if err != nil {
		return values.None, err
	}

	return scr, nil
}

func ValidatePageRanges(pageRanges string) (bool, error) {
	match, err := regexp.Match(`^(([1-9][0-9]*|[1-9][0-9]*)(\s*-\s*|\s*,\s*|))*$`, []byte(pageRanges))
	if err != nil {
		return false, err
	}
	return match, nil
}

/*
 * Print a PDF of the current page.
 * @param source (Document) - Document.
 * @param params (Object) - Optional, An object containing the following properties :
 *   Landscape (Bool) - Paper orientation. Defaults to false.
 *   DisplayHeaderFooter (Bool) - Display header and footer. Defaults to false.
 *   PrintBackground (Bool) - Print background graphics. Defaults to false.
 *   Scale (Float64) - Scale of the webpage rendering. Defaults to 1.
 *   PaperWidth (Float64) - Paper width in inches. Defaults to 8.5 inches.
 *   PaperHeight (Float64) - Paper height in inches. Defaults to 11 inches.
 *   MarginTop (Float64) - Top margin in inches. Defaults to 1cm (~0.4 inches).
 *   MarginBottom (Float64) - Bottom margin in inches. Defaults to 1cm (~0.4 inches).
 *   MarginLeft (Float64) - Left margin in inches. Defaults to 1cm (~0.4 inches).
 *   MarginRight (Float64) - Right margin in inches. Defaults to 1cm (~0.4 inches).
 *   PageRanges (String) - Paper ranges to print, e.g., '1-5, 8, 11-13'. Defaults to the empty string, which means print all pages.
 *   IgnoreInvalidPageRanges (Bool) - to silently ignore invalid but successfully parsed page ranges, such as '3-2'. Defaults to false.
 *   HeaderTemplate (String) - HTML template for the print header. Should be valid HTML markup with following classes used to inject printing values into them: - `date`: formatted print date - `title`: document title - `url`: document location - `pageNumber`: current page number - `totalPages`: total pages in the document For example, `<span class=title></span>` would generate span containing the title.
 *   FooterTemplate (String) - HTML template for the print footer. Should use the same format as the `headerTemplate`.
 *   PreferCSSPageSize (Bool) - Whether or not to prefer page size as defined by css. Defaults to false, in which case the content will be scaled to fit the paper size. *
 * @returns data (Binary) - Returns a base64 encoded string in binary format.
 */
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
	doc := val.(*dynamic.HTMLDocument)
	defer doc.Close()

	pdfParams := page.NewPrintToPDFArgs()

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
			pdfParams.SetLandscape(landscape.Unwrap().(bool))
		}
		displayHeaderFooter, found := params.Get("displayHeaderFooter")
		if found {
			err = core.ValidateType(displayHeaderFooter, core.BooleanType)
			if err != nil {
				return values.None, err
			}
			pdfParams.SetDisplayHeaderFooter(displayHeaderFooter.Unwrap().(bool))
		}
		printBackground, found := params.Get("printBackground")
		if found {
			err = core.ValidateType(printBackground, core.BooleanType)
			if err != nil {
				return values.None, err
			}
			pdfParams.SetPrintBackground(printBackground.Unwrap().(bool))
		}
		scale, found := params.Get("scale")
		if found {
			err = core.ValidateType(scale, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if scale.Type() == core.IntType {
				scale = values.Float(scale.(values.Int))
			}
			pdfParams.SetScale(scale.Unwrap().(float64))
		}
		paperWidth, found := params.Get("paperWidth")
		if found {
			err = core.ValidateType(paperWidth, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if paperWidth.Type() == core.IntType {
				paperWidth = values.Float(paperWidth.(values.Int))
			}
			pdfParams.SetPaperWidth(paperWidth.Unwrap().(float64))
		}
		paperHeight, found := params.Get("paperHeight")
		if found {
			err = core.ValidateType(paperHeight, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if paperHeight.Type() == core.IntType {
				paperHeight = values.Float(paperHeight.(values.Int))
			}
			pdfParams.SetPaperHeight(paperHeight.Unwrap().(float64))
		}
		marginTop, found := params.Get("marginTop")
		if found {
			err = core.ValidateType(marginTop, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if marginTop.Type() == core.IntType {
				marginTop = values.Float(marginTop.(values.Int))
			}
			pdfParams.SetMarginTop(marginTop.Unwrap().(float64))
		}
		marginBottom, found := params.Get("marginBottom")
		if found {
			err = core.ValidateType(marginBottom, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if marginBottom.Type() == core.IntType {
				marginBottom = values.Float(marginBottom.(values.Int))
			}
			pdfParams.SetMarginBottom(marginBottom.Unwrap().(float64))
		}
		marginLeft, found := params.Get("marginLeft")
		if found {
			err = core.ValidateType(marginLeft, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if marginLeft.Type() == core.IntType {
				marginLeft = values.Float(marginLeft.(values.Int))
			}
			pdfParams.SetMarginLeft(marginLeft.Unwrap().(float64))
		}
		marginRight, found := params.Get("marginRight")
		if found {
			err = core.ValidateType(marginRight, core.FloatType, core.IntType)
			if err != nil {
				return values.None, err
			}
			if marginRight.Type() == core.IntType {
				marginRight = values.Float(marginRight.(values.Int))
			}
			pdfParams.SetMarginRight(marginRight.Unwrap().(float64))
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
			pdfParams.SetPageRanges(pageRanges.String())
		}
		ignoreInvalidPageRanges, found := params.Get("ignoreInvalidPageRanges")
		if found {
			err = core.ValidateType(ignoreInvalidPageRanges, core.BooleanType)
			if err != nil {
				return values.None, err
			}
			pdfParams.SetIgnoreInvalidPageRanges(ignoreInvalidPageRanges.Unwrap().(bool))
		}
		headerTemplate, found := params.Get("headerTemplate")
		if found {
			err = core.ValidateType(headerTemplate, core.StringType)
			if err != nil {
				return values.None, err
			}
			pdfParams.SetHeaderTemplate(headerTemplate.String())
		}
		footerTemplate, found := params.Get("footerTemplate")
		if found {
			err = core.ValidateType(footerTemplate, core.StringType)
			if err != nil {
				return values.None, err
			}
			pdfParams.SetFooterTemplate(footerTemplate.String())
		}
		preferCSSPageSize, found := params.Get("preferCSSPageSize")
		if found {
			err = core.ValidateType(preferCSSPageSize, core.BooleanType)
			if err != nil {
				return values.None, err
			}
			pdfParams.SetPreferCSSPageSize(preferCSSPageSize.Unwrap().(bool))
		}
	}

	pdf, err := doc.PrintToPDF(pdfParams)
	if err != nil {
		return values.None, err
	}

	return pdf, nil
}
