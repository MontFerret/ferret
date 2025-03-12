package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// MOUSE moves mouse by given coordinates.
// @param {HTMLDocument} document - HTML document.
// @param {Int|Float} x - X coordinate.
// @param {Int|Float} y - Y coordinate.
func MouseMoveXY(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 3)

	if err != nil {
		return core.None, err
	}

	doc, err := drivers.ToDocument(args[0])

	if err != nil {
		return core.None, err
	}

	err = core.ValidateType(args[1], types.Int, types.Float)

	if err != nil {
		return core.None, err
	}

	err = core.ValidateType(args[2], types.Int, types.Float)

	if err != nil {
		return core.None, err
	}

	x := internal.ToFloat(args[0])
	y := internal.ToFloat(args[1])

	return core.None, doc.MoveMouseByXY(ctx, x, y)
}
