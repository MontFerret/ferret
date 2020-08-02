package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// INCLUDES checks whether a container includes a given value.
// @param {String | Any[] | Object | Iterable} haystack - The value container.
// @param {Any} needle - The target value to assert.
// @return {Boolean} - A boolean value that indicates whether a container contains a given value.
func Includes(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	var result values.Boolean
	haystack := args[0]
	needle := args[1]

	switch v := haystack.(type) {
	case values.String:
		result = v.Contains(values.NewString(needle.String()))

		break
	case *values.Array:
		_, result = v.FindOne(func(value core.Value, _ int) bool {
			return needle.Compare(value) == 0
		})

		break
	case *values.Object:
		_, result = v.Find(func(value core.Value, _ string) bool {
			return needle.Compare(value) == 0
		})

		break
	case core.Iterable:
		iter, err := v.Iterate(ctx)

		if err != nil {
			return values.False, err
		}

		err = core.ForEach(ctx, iter, func(value core.Value, key core.Value) bool {
			if needle.Compare(value) == 0 {
				result = values.True

				return false
			}

			return true
		})

		if err != nil {
			return values.False, err
		}
	default:
		return values.None, core.TypeError(haystack.Type(),
			types.String,
			types.Array,
			types.Object,
			core.NewType("Iterable"),
		)
	}

	return result, nil
}
