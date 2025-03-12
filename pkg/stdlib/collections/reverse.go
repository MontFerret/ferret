package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// REVERSE returns the reverse of a given string or array value.
// @param {String | Any[]} value - The string or array to reverse.
// @return {String | Any[]} - A reversed version of a given value.
func Reverse(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.EmptyString, err
	}

	switch col := args[0].(type) {
	case core.String:
		runes := []rune(string(col))
		size := len(runes)

		// Reverse
		for i := 0; i < size/2; i++ {
			runes[i], runes[size-1-i] = runes[size-1-i], runes[i]
		}

		return core.NewString(string(runes)), nil
	case *internal.Array:
		size := int(col.Length())
		result := internal.NewArray(size)

		for i := size - 1; i >= 0; i-- {
			result.Push(col.Get(i))
		}

		return result, nil

	default:
		return core.None, core.TypeError(args[0], types.Array, types.String)
	}
}
