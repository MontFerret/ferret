package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// COUNT_DISTINCT computes the number of distinct elements in the given collection and returns the count as an integer.
func CountDistinct(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	collection, err := runtime.CastCollection(arg)

	if err != nil {
		return runtime.ZeroInt, err
	}

	// TODO: Use storage backend
	hashmap := map[uint64]bool{}
	var res runtime.Int

	err = runtime.ForEach(ctx, collection, func(c context.Context, value, idx runtime.Value) (runtime.Boolean, error) {
		hash := value.Hash()

		_, exists := hashmap[hash]

		if !exists {
			hashmap[hash] = true
			res++
		}

		return true, nil
	})

	if err != nil {
		return runtime.ZeroInt, err
	}

	return res, nil
}
