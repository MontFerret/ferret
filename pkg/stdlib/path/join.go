package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// JOIN joins any number of path elements into a single path, separating them with slashes.
// @param {String, repeated | String[]} elements - The path elements
// @return {String} - Single path from the given elements.
func Join(_ context.Context, args ...core.Value) (core.Value, error) {
	argsCount := len(args)

	if argsCount == 0 {
		return values.EmptyString, nil
	}

	var arr *values.Array

	switch arg := args[0].(type) {
	case *values.Array:
		arr = arg
	default:
		arr = values.NewArrayWith(args...)
	}

	elems := make([]string, arr.Length())

	for idx := 0; idx < arr.Length(); idx++ {
		arrElem := arr.Get(idx)
		err := core.ValidateType(arrElem, types.String)

		if err != nil {
			return values.None, err
		}

		elems[idx] = arrElem.String()
	}

	return values.NewString(path.Join(elems...)), nil
}
