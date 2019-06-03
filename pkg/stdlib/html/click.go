package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Click dispatches click event on a given element
// @param source (Open | GetElement) - Event source.
// @param selector (String, optional) - Optional selector. Only used when a document instance is passed.
func Click(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.False, err
	}

	// CLICK(el)
	if len(args) == 1 {
		el, err := toElement(args[0])

		if err != nil {
			return values.False, err
		}

		return el.Click(ctx)
	}

	// CLICK(doc, selector)
	doc, err := toDocument(args[0])

	if err != nil {
		return values.False, err
	}

	selector := args[1].String()

	return doc.ClickBySelector(ctx, values.NewString(selector))
}
