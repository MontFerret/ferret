package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Assert checks whether value is a boolean value.
// @param value (Boolean|Array) - Expression(s) to test for truthiness.
// @param (String) - Message to display on error.
var Assert = Assertion{
	DefaultMessage: func(_ []core.Value) string {
		return "be true"
	},
	MinArgs: 1,
	MaxArgs: 2,
	Fn: func(ctx context.Context, args []core.Value) (bool, error) {
		if err := core.ValidateType(args[0], types.Boolean, types.Array); err != nil {
			return false, err
		}

		result := true

		switch v := args[0].(type) {
		case values.Boolean:
			result = bool(v)
			break
		case *values.Array:
			v.ForEach(func(value core.Value, idx int) bool {
				if value.Compare(values.False) == 0 {
					result = false

					return false
				}

				return true
			})
		}

		return result, nil
	},
}
