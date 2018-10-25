package html

import (
	"context"

<<<<<<< HEAD
=======
	"github.com/MontFerret/ferret/pkg/html/dynamic"
>>>>>>> 9f24172... rewrite comments
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Click dispatches click event on a given element
// @param source (Document | Element) - Event source.
// @param selector (String, optional) - Optional selector. Only used when a document instance is passed.
<<<<<<< HEAD
func Click(ctx context.Context, args ...core.Value) (core.Value, error) {
=======
func Click(_ context.Context, args ...core.Value) (core.Value, error) {
>>>>>>> 9f24172... rewrite comments
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
