package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// UNION_DISTINCT returns the union of all passed arrays with unique values.
// @param {Any[], repeated} arrays - List of arrays to combine.
// @return {Any[]} - All unique array elements combined in a single array, in any order.
func UnionDistinct(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	firstArrLen := args[0].(*values.Array).Length()
	result := values.NewArray(len(args) * int(firstArrLen))
	hashes := make(map[uint64]bool)

	for _, arg := range args {
		err := core.ValidateType(arg, types.Array)

		if err != nil {
			return values.None, err
		}

		arr := arg.(*values.Array)

		arr.ForEach(func(value core.Value, _ int) bool {
			h := value.Hash()

			_, exists := hashes[h]

			if !exists {
				hashes[h] = true
				result.Push(value)
			}

			return true
		})
	}

	return result, nil
}
