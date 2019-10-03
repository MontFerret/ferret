package html

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SCREENSHOT takes a screenshot of a given page.
// @param target (HTMLPage|String) - Target page or url.
// @param params (Object) - Optional, An object containing the following properties :
// 		x (Float|Int) - Optional, X position of the viewport.
// 		x (Float|Int) - Optional,Y position of the viewport.
// 		width (Float|Int) - Optional, Width of the viewport.
// 		height (Float|Int) - Optional, Height of the viewport.
// 		format (String) - Optional, Either "jpeg" or "png".
// 		quality (Int) - Optional, Quality, in [0, 100], only for jpeg format.
// @returns data (Binary) - Returns a base64 encoded string in binary format.
func Screenshot(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	arg1 := args[0]

	err = core.ValidateType(arg1, drivers.HTMLDocumentType, types.String)

	if err != nil {
		return values.None, err
	}

	page, closeAfter, err := OpenOrCastPage(ctx, arg1)

	if err != nil {
		return values.None, err
	}

	defer func() {
		if closeAfter {
			page.Close()
		}
	}()

	screenshotParams := drivers.ScreenshotParams{
		X:       0,
		Y:       0,
		Width:   -1,
		Height:  -1,
		Format:  drivers.ScreenshotFormatJPEG,
		Quality: 100,
	}

	if len(args) == 2 {
		arg2 := args[1]
		err = core.ValidateType(arg2, types.Object)

		if err != nil {
			return values.None, err
		}

		params, ok := arg2.(*values.Object)

		if !ok {
			return values.None, core.Error(core.ErrInvalidType, "expected object")
		}

		format, found := params.Get("format")

		if found {
			err = core.ValidateType(format, types.String)

			if err != nil {
				return values.None, err
			}

			if !drivers.IsScreenshotFormatValid(format.String()) {
				return values.None, core.Error(
					core.ErrInvalidArgument,
					fmt.Sprintf("format is not valid, expected jpeg or png, but got %s", format.String()))
			}

			screenshotParams.Format = drivers.ScreenshotFormat(format.String())
		}

		x, found := params.Get("x")

		if found {
			err = core.ValidateType(x, types.Float, types.Int)

			if err != nil {
				return values.None, err
			}

			if x.Type() == types.Int {
				screenshotParams.X = values.Float(x.(values.Int))
			}
		}

		y, found := params.Get("y")

		if found {
			err = core.ValidateType(y, types.Float, types.Int)

			if err != nil {
				return values.None, err
			}

			if y.Type() == types.Int {
				screenshotParams.Y = values.Float(y.(values.Int))
			}
		}

		width, found := params.Get("width")

		if found {
			err = core.ValidateType(width, types.Float, types.Int)

			if err != nil {
				return values.None, err
			}

			if width.Type() == types.Int {
				screenshotParams.Width = values.Float(width.(values.Int))
			}
		}

		height, found := params.Get("height")

		if found {
			err = core.ValidateType(height, types.Float, types.Int)

			if err != nil {
				return values.None, err
			}

			if height.Type() == types.Int {
				screenshotParams.Height = values.Float(height.(values.Int))
			}
		}

		quality, found := params.Get("quality")

		if found {
			err = core.ValidateType(quality, types.Int)

			if err != nil {
				return values.None, err
			}

			screenshotParams.Quality = quality.(values.Int)
		}
	}

	scr, err := page.CaptureScreenshot(ctx, screenshotParams)

	if err != nil {
		return values.None, err
	}

	return scr, nil
}
