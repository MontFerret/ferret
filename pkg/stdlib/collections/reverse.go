package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// REVERSE returns the reverse of a given string or array value.
// @param {String | Any[]} value - The string or array to reverse.
// @return {String | Any[]} - A reversed version of a given value.
func Reverse(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	switch col := arg.(type) {
	case runtime.String:
		runes := []rune(string(col))
		size := len(runes)

		// Reverse
		for i := 0; i < size/2; i++ {
			runes[i], runes[size-1-i] = runes[size-1-i], runes[i]
		}

		return runtime.NewString(string(runes)), nil
	case runtime.List:
		size, err := col.Length(ctx)

		if err != nil {
			return runtime.None, err
		}

		result := runtime.NewArray(int(size))

		for i := size - 1; i >= 0; i-- {
			item, err := col.Get(ctx, i)

			if err != nil {
				return runtime.None, err
			}

			_ = result.Add(ctx, item)
		}

		return result, nil
	default:
		return runtime.None, runtime.TypeErrorOf(arg, runtime.TypeList, runtime.TypeString)
	}
}
