package drivers

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

const (
	ScreenshotFormatPNG  ScreenshotFormat = "png"
	ScreenshotFormatJPEG ScreenshotFormat = "jpeg"
)

type (
	ScreenshotFormat string

	ScreenshotParams struct {
		X      core.Float
		Y      core.Float
		Width  core.Float
		Height core.Float
		Format  ScreenshotFormat
		Quality core.Int
	}
)

func IsScreenshotFormatValid(format string) bool {
	value := ScreenshotFormat(format)

	return value == ScreenshotFormatPNG || value == ScreenshotFormatJPEG
}

func NewDefaultHTMLPDFParams() PDFParams {
	return PDFParams{
		Landscape:           core.False,
		DisplayHeaderFooter: core.False,
		PrintBackground:     core.False,
		Scale:               core.Float(1),
		PaperWidth:          core.Float(8.5),
		PaperHeight:         core.Float(11),
		MarginTop:           core.Float(0.4),
		MarginBottom:        core.Float(0.4),
		MarginLeft:          core.Float(0.4),
		MarginRight:         core.Float(0.4),
		PageRanges:          core.EmptyString,
		HeaderTemplate:      core.EmptyString,
		FooterTemplate:      core.EmptyString,
		PreferCSSPageSize:   core.False,
	}
}
