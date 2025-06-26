package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ZIP returns an object assembled from the separate parameters keys and values.
// Keys and values must be arrays and have the same length.
// @param {String[]} keys - An array of strings, to be used as key names in the result.
// @param {hashMap[]} values - An array of core.Value, to be used as key values.
// @return {hashMap} - An object with the keys and values assembled.
// TODO: REWRITE TO USE LIST & MAP instead
func Zip(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return core.None, err
	}

	for _, arg := range args {
		err = core.ValidateType(arg, types.Array)

		if err != nil {
			return core.None, err
		}
	}

	keys := args[0].(*runtime.Array)
	vals := args[1].(*runtime.Array)

	keysSize, _ := keys.Length(ctx)
	valsSize, _ := vals.Length(ctx)

	if keysSize != valsSize {
		return core.None, core.Error(
			core.ErrInvalidArgument,
			fmt.Sprintf("keys and values must have the same length. got keys: %d, values: %d",
				keysSize, valsSize,
			),
		)
	}

	err = validateArrayOf(ctx, types.String, keys)

	if err != nil {
		return core.None, err
	}

	zipped := runtime.NewObject()

	var k core.String
	var val core.Value
	var exists bool
	keyExists := map[core.String]bool{}

	_ = keys.ForEach(ctx, func(c context.Context, key core.Value, idx core.Int) (core.Boolean, error) {
		k = key.(core.String)

		// this is necessary to implement ArangoDB's behavior.
		// in ArangoDB the first value in values is
		// associated with each key. Ex.:
		// -- query --
		// > RETURN ZIP(
		// >     ["a", "b", "a"], [1, 2, 3]
		// > )
		// -- result --
		// > [{"a": 1,"b": 2}]
		if _, exists = keyExists[k]; exists {
			return true, nil
		}

		keyExists[k] = true

		val, _ = vals.Get(c, idx)

		cloneable, ok := val.(core.Cloneable)

		if ok {
			val, _ = cloneable.Clone(c)
		}

		_ = zipped.Set(c, k, val)

		return true, nil
	})

	return zipped, nil
}
