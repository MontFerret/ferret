package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// SCROLL_BOTTOM scrolls the document's window to its bottom.
// @param {HTMLDocument} document - HTML document.
// @param {Int | Float} x - X coordinate.
// @param {Int | Float} y - Y coordinate.
// @param {Object} [params] - Scroll params.
// @param {String} [params.behavior="instant"] - Scroll behavior
// @param {String} [params.block="center"] - Scroll vertical alignment.
// @param {String} [params.inline="center"] - Scroll horizontal alignment.
func ScrollBottom(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	doc, err := drivers.ToDocument(args[0])

	if err != nil {
		return values.None, err
	}

	var opts drivers.ScrollOptions

	if len(args) > 1 {
		opts, err = toScrollOptions(args[1])

		if err != nil {
			return values.None, err
		}
	}

	return values.True, doc.ScrollBottom(ctx, opts)
}
