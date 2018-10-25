package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

<<<<<<< HEAD
// InnerText returns inner text string of a given or matched by CSS selector element
// @param doc (HTMLDocument|HTMLElement) - Parent document or element.
// @param selector (String, optional) - String of CSS selector.
// @returns (String) - Inner text if an element found, otherwise empty string.
func InnerText(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)
=======
// InnerText returns inner text of a matched element
// @param doc (HTMLDocument|HTMLElement) - Parent document or element.
// @param selector (String) - String of CSS selector.
// @returns (String) - Inner text if an element found, otherwise empty string.
func InnerText(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)
>>>>>>> 9f24172... rewrite comments

	if err != nil {
		return values.EmptyString, err
	}

	el, err := resolveElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 1 {
		return el.InnerText(ctx), nil
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)

	return el.InnerTextBySelector(ctx, selector), nil
}
