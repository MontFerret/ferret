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

	// page or document or element
	err = core.ValidateType(args[0], drivers.HTMLPageType, drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	selector := values.EmptyString

	if len(args) > 1 {
		err = core.ValidateType(args[1], types.String)

		if err != nil {
			return values.None, err
		}

		selector = args[1].(values.String)
	}

	switch n := args[0].(type) {
	case drivers.HTMLPage:
		if selector == values.EmptyString {
			return values.None, core.Error(core.ErrMissedArgument, "selector")
		}

		return values.None, n.GetMainFrame().MoveMouseBySelector(ctx, selector)
	case drivers.HTMLDocument:
		if selector == values.EmptyString {
			return values.None, core.Error(core.ErrMissedArgument, "selector")
		}

		return values.None, n.MoveMouseBySelector(ctx, selector)
	case drivers.HTMLElement:
		if selector == values.EmptyString {
			return values.None, n.Hover(ctx)
		}

		found, err := n.QuerySelector(ctx, selector)

		if err != nil {
			return values.None, err
		}

		if found == values.None {
			return values.None, core.Errorf(core.ErrNotFound, "element by selector %s", selector)
		}

		el, ok := found.(drivers.HTMLElement)

		if !ok {
			return values.None, core.Errorf(core.ErrNotFound, "element by selector %s", selector)
		}

		defer el.Close()

		return values.None, el.Hover(ctx)
	default:
		return values.None, core.TypeError(n.Type(), drivers.HTMLDocumentType, drivers.HTMLElementType)
	}
}
