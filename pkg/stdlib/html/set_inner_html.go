package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// INNER_HTML_SET sets inner HTML string to a given or matched by CSS selector element
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} htmlOrSelector - HTML or CSS selector.
// @param {String} [html] - String of inner HTML.
func SetInnerHTML(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 2 {
		err := core.ValidateType(args[1], types.String)

		if err != nil {
			return values.None, err
		}

		return values.None, el.SetInnerHTML(ctx, values.ToString(args[1]))
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[2], types.String)

	if err != nil {
		return values.None, err
	}

	innerHTML := values.ToString(args[2])

	return values.None, el.SetInnerHTMLBySelector(ctx, selector, innerHTML)
}
