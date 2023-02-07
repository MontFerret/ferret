package drivers

import "github.com/MontFerret/ferret/pkg/runtime/values"

const (
	ScreenshotFormatPNG  ScreenshotFormat = "png"
	ScreenshotFormatJPEG ScreenshotFormat = "jpeg"
)

type (
	ScreenshotFormat string

	ScreenshotParams struct {
		X       values.Float
		Y       values.Float
		Width   values.Float
		Height  values.Float
		Format  ScreenshotFormat
		Quality values.Int
	}
)

func IsScreenshotFormatValid(format string) bool {
	value := ScreenshotFormat(format)

	return value == ScreenshotFormatPNG || value == ScreenshotFormatJPEG
}

func NewDefaultHTMLPDFParams() PDFParams {
	return PDFParams{
		Landscape:           values.False,
		DisplayHeaderFooter: values.False,
		PrintBackground:     values.False,
		Scale:               values.Float(1),
		PaperWidth:          values.Float(8.5),
		PaperHeight:         values.Float(11),
		MarginTop:           values.Float(0.4),
		MarginBottom:        values.Float(0.4),
		MarginLeft:          values.Float(0.4),
		MarginRight:         values.Float(0.4),
		PageRanges:          values.EmptyString,
		HeaderTemplate:      values.EmptyString,
		FooterTemplate:      values.EmptyString,
		PreferCSSPageSize:   values.False,
	}
}
