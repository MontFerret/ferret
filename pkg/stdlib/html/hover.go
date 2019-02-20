package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Hover  fetches an element with selector, scrolls it into view if needed, and then uses page.mouse to hover over the center of the element.
// If there's no element matching selector, the method returns an error.
// @param docOrEl (HTMLDocument|HTMLElement) - Target document or element.
// @param selector (String, options) - If document is passed, this param must represent an element selector.
func Hover(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	// document or element
	err = core.ValidateType(args[0], drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	if len(args) == 2 {
		err = core.ValidateType(args[1], types.String)

		if err != nil {
			return values.None, err
		}

		// Document with a selector
		doc := args[0].(drivers.HTMLDocument)
		selector := args[1].(values.String)

		return values.None, doc.HoverBySelector(ctx, selector)
	}

	// Element
	el := args[0].(drivers.HTMLElement)

	return values.None, el.Hover(ctx)
}
