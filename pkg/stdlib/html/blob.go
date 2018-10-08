package html

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

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

	var doc *dynamic.HTMLDocument
	var ok bool
	if arg1.Type() == core.StringType {
		buf, err := Document(ctx, arg1, values.NewBoolean(true))
		if err != nil {
			return values.None, err
		}
		doc, ok = buf.(*dynamic.HTMLDocument)
		defer doc.Close()
	} else {
		doc, ok = arg1.(*dynamic.HTMLDocument)
	}

	if !ok {
		return values.None, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

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
