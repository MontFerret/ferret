package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// JOIN joins any number of path elements into a single path, separating them with slashes.
// @param {String, repeated | String[]} elements - The path elements
// @return {String} - Single path from the given elements.
func Join(ctx context.Context, args ...core.Value) (core.Value, error) {
	argsCount := len(args)

	if argsCount == 0 {
		return core.EmptyString, nil
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
		arrElem, _ := arr.Get(ctx, idx)
		err := core.ValidateType(arrElem, types.String)

		if err != nil {
			return core.None, err
		}

		elems[idx] = arrElem.String()
	}

	return core.NewString(path.Join(elems...)), nil
}
