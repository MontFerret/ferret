package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ScrollInto scrolls an element on.
// @param docOrEl (HTMLDocument|HTMLElement) - Target document or element.
// @param selector (String, options) - If document is passed, this param must represent an element selector.
func ScrollInto(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	if len(args) == 2 {
		err = core.ValidateType(args[0], drivers.HTMLDocumentType)

		if err != nil {
			return values.None, err
		}

		err = core.ValidateType(args[1], types.String)

		if err != nil {
			return values.None, err
		}

		// Document with a selector
		doc := args[0].(drivers.HTMLDocument)
		selector := args[1].(values.String)

		return values.None, doc.ScrollBySelector(ctx, selector)
	}

	err = core.ValidateType(args[0], drivers.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	// Element
	el := args[0].(drivers.HTMLElement)

	return values.None, el.ScrollIntoView(ctx)
}
