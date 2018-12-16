package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ScrollInto scrolls an element on.
// @param docOrEl (HTMLDocument|HTMLElement) - Target document or element.
// @param selector (String, options) - If document is passed, this param must represent an element selector.
func ScrollInto(_ context.Context, args ...core.Value) (core.Value, error) {
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

		return values.None, doc.ScrollBySelector(selector)
	}

	// Element
	el, ok := args[0].(values.DHTMLNode)

	if !ok {
		return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	return values.None, el.ScrollIntoView()
}
