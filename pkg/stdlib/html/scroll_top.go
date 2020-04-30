package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// SCROLL_TOP scrolls the document's window to its top.
// @param doc (HTMLDocument) - Target document.
// @param options (ScrollOptions) - Scroll options. Optional.
func ScrollTop(ctx context.Context, args ...core.Value) (core.Value, error) {
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

	return values.None, doc.ScrollTop(ctx, opts)
}
