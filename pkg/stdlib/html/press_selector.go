package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// PRESS_SELECTOR presses a keyboard key.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - CSS selector.
// @param {String | String[]} key - Target keyboard key(s).
// @param {Int} [presses=1] - Count of presses.
func PressSelector(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 4)

	if err != nil {
		return core.False, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.False, err
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return core.None, err
	}

	count := core.NewInt(1)

	if len(args) == 4 {
		countArg := internal.ToInt(args[3])

		if countArg > 0 {
			count = countArg
		}
	}

	keysArg := args[2]

	switch keys := keysArg.(type) {
	case core.String:
		return core.True, el.PressBySelector(ctx, selector, []core.String{keys}, count)
	case *internal.Array:
		return core.True, el.PressBySelector(ctx, selector, internal.ToStrings(keys), count)
	default:
		return core.None, core.TypeError(keysArg.Type(), types.String, types.Array)
	}
}
