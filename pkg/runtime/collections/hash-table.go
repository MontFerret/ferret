package collections

import "github.com/MontFerret/ferret/pkg/runtime/core"

func ToHashTable(iterator Iterator) (map[uint64]core.Value, error) {
	result := make(map[uint64]core.Value)

	for iterator.HasNext() {
		val, _, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		h := val.Hash()

		_, exists := result[h]

		if !exists {
			result[h] = val
		}
	}

	return result, nil
}
