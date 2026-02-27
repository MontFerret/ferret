package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("FIRST", First).
		Add("LAST", Last).
		Add("POP", Pop).
		Add("SHIFT", Shift).
		Add("SORTED", Sorted).
		Add("SORTED_UNIQUE", SortedUnique).
		Add("UNIQUE", Unique)

	ns.Function().A2().
		Add("NTH", Nth).
		Add("REMOVE_NTH", RemoveNth).
		Add("REMOVE_VALUES", RemoveValues)

	ns.Function().Var().
		Add("APPEND", Append).
		Add("FLATTEN", Flatten).
		Add("INTERSECTION", Intersection).
		Add("MINUS", Minus).
		Add("OUTERSECTION", Outersection).
		Add("POSITION", Position).
		Add("PUSH", Push).
		Add("REMOVE_VALUE", RemoveValue).
		Add("SLICE", Slice).
		Add("UNION", Union).
		Add("UNION_DISTINCT", UnionDistinct).
		Add("UNSHIFT", Unshift)
}

func ToUniqueList(ctx context.Context, list runtime.List) (runtime.List, error) {
	hashTable := make(map[uint64]bool)

	return list.Filter(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		hash := value.Hash()

		if _, exists := hashTable[hash]; !exists {
			hashTable[hash] = true

			return true, nil
		}

		return false, nil
	})
}
