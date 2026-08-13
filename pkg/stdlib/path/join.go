package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// join joins any number of path elements into a single path, separating them with slashes.
// @param elements {String, repeated | String[]} The path elements
// @return {String} Single path from the given elements.
func Join(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	argsCount := len(args)

	if argsCount == 0 {
		return runtime.EmptyString, nil
	}

	var arr *runtime.Array

	switch arg := args[0].(type) {
	case *runtime.Array:
		arr = arg
	default:
		arr = runtime.NewArrayWith(args...)
	}

	size, _ := arr.Length(ctx)
	elems := make([]string, size)

	for idx := runtime.ZeroInt; idx < size; idx++ {
		arrElem, _ := arr.At(ctx, idx)
		err := runtime.ValidateArgType(arrElem, 0, runtime.TypeString)

		if err != nil {
			return runtime.None, err
		}

		elems[idx] = arrElem.String()
	}

	return runtime.NewString(path.Join(elems...)), nil
}
