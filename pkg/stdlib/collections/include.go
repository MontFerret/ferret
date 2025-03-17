package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// INCLUDES checks whether a container includes a given value.
// @param {String | Any[] | hashMap | Iterable} haystack - The value container.
// @param {Any} needle - The target value to assert.
// @return {Boolean} - A boolean value that indicates whether a container contains a given value.
func Includes(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return core.None, err
	}

	var result core.Boolean
	haystack := args[0]
	needle := args[1]

	switch v := haystack.(type) {
	case core.String:
		result = v.Contains(core.NewString(needle.String()))

		break
	case *internal.Array:
		_, result = v.FindOne(func(value core.Value, _ int) bool {
			return core.CompareValues(needle, value) == 0
		})

		break
	case *internal.Object:
		_, result = v.Find(func(value core.Value, _ string) bool {
			return core.CompareValues(needle, value) == 0
		})

		break
	case core.Iterable:
		iter, err := v.Iterate(ctx)

		if err != nil {
			return core.False, err
		}

		err = core.ForEach(ctx, iter, func(value core.Value, key core.Value) bool {
			if core.CompareValues(needle, value) == 0 {
				result = core.True

				return false
			}

			return true
		})

		if err != nil {
			return core.False, err
		}
	default:
		return core.None, core.TypeError(haystack,
			types.String,
			types.Array,
			types.Object,
			types.Iterable,
		)
	}

	return result, nil
}
