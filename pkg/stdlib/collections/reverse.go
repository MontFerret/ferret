package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// REVERSE returns the reverse of a given string or array value.
// @param {String | Any[]} value - The string or array to reverse.
// @return {String | Any[]} - A reversed version of a given value.
func Reverse(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[0], types.Array, types.String)

	if err != nil {
		return values.None, err
	}

	switch col := args[0].(type) {
	case values.String:
		runes := []rune(string(col))
		size := len(runes)

		// Reverse
		for i := 0; i < size/2; i++ {
			runes[i], runes[size-1-i] = runes[size-1-i], runes[i]
		}

		return values.NewString(string(runes)), nil
	case *values.Array:
		size := int(col.Length())
		result := values.NewArray(size)

		for i := size - 1; i >= 0; i-- {
			result.Push(col.Get(values.NewInt(i)))
		}

		return result, nil

	default:
		return values.None, nil
	}
}
