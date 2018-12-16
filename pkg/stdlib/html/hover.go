package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Hover  fetches an element with selector, scrolls it into view if needed, and then uses page.mouse to hover over the center of the element.
// If there's no element matching selector, the method returns an error.
// @param docOrEl (HTMLDocument|HTMLElement) - Target document or element.
// @param selector (String, options) - If document is passed, this param must represent an element selector.
func Hover(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	// document or element
	err = core.ValidateType(args[0], core.HTMLDocumentType, core.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	if len(args) == 2 {
		err = core.ValidateType(args[1], core.StringType)

		if err != nil {
			return values.None, err
		}

		// Document with a selector
		doc, ok := args[0].(values.DHTMLDocument)

		if !ok {
			return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
		}

		selector := args[1].(values.String)

		return values.None, doc.HoverBySelector(selector)
	}

	// Element
	el, ok := args[0].(values.DHTMLNode)

	if !ok {
		return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	return values.None, el.Hover()
}
