package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// INCLUDES checks whether a container includes a given value.
// @param {String | Any[] | hashMap | Iterable} haystack - The value container.
// @param {Any} needle - The target value to assert.
// @return {Boolean} - A boolean value that indicates whether a container contains a given value.
func Includes(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 2)

	if err != nil {
		return runtime.None, err
	}

	var result runtime.Boolean
	haystack := args[0]
	needle := args[1]

	switch v := haystack.(type) {
	case runtime.String:
		result = v.Contains(runtime.NewString(needle.String()))

		break
	case runtime.List:
		_, result, err = v.FindOne(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			return runtime.CompareValues(needle, value) == 0, nil
		})

		break
	case runtime.Map:
		_, result, err = v.FindOne(ctx, func(c context.Context, value, _ runtime.Value) (runtime.Boolean, error) {
			return runtime.CompareValues(needle, value) == 0, nil
		})

		break
	case runtime.Iterable:
		iter, err := v.Iterate(ctx)

		if err != nil {
			return runtime.False, err
		}

		err = runtime.ForEach(ctx, iter, func(c context.Context, value runtime.Value, key runtime.Value) (runtime.Boolean, error) {
			if runtime.CompareValues(needle, value) == 0 {
				result = runtime.True

				return false, nil
			}

			return true, nil
		})

		if err != nil {
			return runtime.False, err
		}
	default:
		return runtime.None, runtime.TypeError(haystack,
			runtime.TypeString,
			runtime.TypeList,
			runtime.TypeMap,
			runtime.TypeIterable,
		)
	}

	return result, nil
}
