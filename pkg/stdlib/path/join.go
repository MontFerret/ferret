package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// JOIN joins any number of path elements into a single path, separating them with slashes.
// @param {String, repeated | String[]} elements - The path elements
// @return {String} - Single path from the given elements.
func Join(_ context.Context, args ...core.Value) (core.Value, error) {
	argsCount := len(args)

	if argsCount == 0 {
		return core.EmptyString, nil
	}

	var arr *internal.Array

	switch arg := args[0].(type) {
	case *internal.Array:
		arr = arg
	default:
		arr = internal.NewArrayWith(args...)
	}

	elems := make([]string, arr.Length())

	for idx := 0; idx < arr.Length(); idx++ {
		arrElem := arr.Get(idx)
		err := core.ValidateType(arrElem, types.String)

		if err != nil {
			return core.None, err
		}

		elems[idx] = arrElem.String()
	}

	return core.NewString(path.Join(elems...)), nil
}
