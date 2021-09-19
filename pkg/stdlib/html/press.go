package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// PRESS presses a keyboard key.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String | String[]} key - Target keyboard key(s).
// @param {Int} [presses=1] - Count of presses.
func Press(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.False, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.False, err
	}

	count := values.NewInt(1)

	if len(args) == 3 {
		countArg := values.ToInt(args[2])

		if countArg > 0 {
			count = countArg
		}
	}

	keysArg := args[1]

	switch keys := keysArg.(type) {
	case values.String:
		return values.True, el.Press(ctx, []values.String{keys}, count)
	case *values.Array:
		return values.True, el.Press(ctx, values.ToStrings(keys), count)
	default:
		return values.None, core.TypeError(keysArg.Type(), types.String, types.Array)
	}
}
