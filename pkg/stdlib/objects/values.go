package objects

import "github.com/MontFerret/ferret/pkg/runtime"

// VALUES return the attribute values of the object as an array.
// @param {hashMap} object - Target object.
// @return {Any[]} - Values of document returned in any order.
// TODO: REWRITE TO USE LIST & MAP instead
func Values(ctx runtime.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeObject)

	if err != nil {
		return runtime.None, err
	}

	obj := args[0].(*runtime.Object)
	vals := ctx.Alloc().Array(0)

	_ = obj.ForEach(ctx, func(c runtime.Context, val, key runtime.Value) (runtime.Boolean, error) {
		val, err := runtime.CloneOrCopy(c, val)

		if err != nil {
			return runtime.False, err
		}

		if err := vals.Append(c, val); err != nil {
			return runtime.False, err
		}

		return true, nil
	})

	return vals, nil
}
