package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SCROLL scrolls by given coordinates.
// @param {HTMLDocument} document - HTML document.
// @param {Int | Float} x - X coordinate.
// @param {Int | Float} y - Y coordinate.
// @param {Object, optional} params - Scroll params.
// @param {String, optional} params.behavior - Scroll behavior
// @param {String, optional} params.block - Scroll vertical alignment.
// @param {String, optional}  params.inline - Scroll horizontal alignment.
func ScrollXY(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 4)

	if err != nil {
		return values.None, err
	}

	doc, err := drivers.ToDocument(args[0])

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[2], types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	x := values.ToFloat(args[1])
	y := values.ToFloat(args[2])

	var opts drivers.ScrollOptions

	if len(args) > 3 {
		opts, err = toScrollOptions(args[3])

		if err != nil {
			return values.None, err
		}
	}

	return values.None, doc.ScrollByXY(ctx, x, y, opts)
}
