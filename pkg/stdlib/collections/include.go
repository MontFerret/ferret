package collections

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Includes checks membership, preferring a host's Contains capability over a scan.
// Strings convert the needle to text. Iterable scans compare yielded values,
// including map values, using canonical equality and stop at the first match.
// Only the acquired iterator is closed; traversal, equality, cancellation, and
// close failures are preserved, including a close failure after a match.
// Host capabilities receive the caller's context and must observe cancellation
// while they retain control; an entry check cannot interrupt a blocking host.
// @param haystack {String | Containable | Iterable} The value container.
// @param needle {Any} The target value to assert.
// @return {Boolean} A boolean value that indicates whether a container contains a given value.
func Includes(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.False, err
	}

	var err error
	var result runtime.Boolean
	haystack := arg1
	needle := arg2

	switch v := haystack.(type) {
	case runtime.String:
		result = v.Contains(runtime.NewString(needle.String()))
	case runtime.Containable:
		result, err = v.Contains(ctx, needle)
	case runtime.Iterable:
		err = runtime.ForEach(ctx, v, func(c context.Context, value runtime.Value, key runtime.Value) (runtime.Boolean, error) {
			if err := c.Err(); err != nil {
				return false, err
			}

			equal, err := runtime.EqualValues(c, needle, value)
			if err != nil {
				return false, err
			}

			if err := c.Err(); err != nil {
				return false, err
			}

			if equal {
				result = runtime.True

				return false, nil
			}

			return true, nil
		})
	default:
		return runtime.None, runtime.TypeErrorOf(haystack,
			runtime.TypeString,
			runtime.TypeList,
			runtime.TypeMap,
			runtime.TypeContainable,
			runtime.TypeIterable,
		)
	}

	if canceled := ctx.Err(); canceled != nil && !errors.Is(err, canceled) {
		err = errors.Join(err, canceled)
	}

	if err != nil {
		return runtime.False, err
	}

	return result, nil
}
