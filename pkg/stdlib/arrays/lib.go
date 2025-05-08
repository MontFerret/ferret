package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) error {
	return ns.RegisterFunctions(
		runtime.NewFunctionsFromMap(map[string]runtime.Function{
			"APPEND":         Append,
			"FIRST":          First,
			"FLATTEN":        Flatten,
			"INTERSECTION":   Intersection,
			"LAST":           Last,
			"MINUS":          Minus,
			"NTH":            Nth,
			"OUTERSECTION":   Outersection,
			"POP":            Pop,
			"POSITION":       Position,
			"PUSH":           Push,
			"REMOVE_NTH":     RemoveNth,
			"REMOVE_VALUE":   RemoveValue,
			"REMOVE_VALUES":  RemoveValues,
			"SHIFT":          Shift,
			"SLICE":          Slice,
			"SORTED":         Sorted,
			"SORTED_UNIQUE":  SortedUnique,
			"UNION":          Union,
			"UNION_DISTINCT": UnionDistinct,
			"UNIQUE":         Unique,
			"UNSHIFT":        Unshift,
		}))
}

func ToUniqueList(ctx context.Context, list runtime.List) (runtime.List, error) {
	hashTable := make(map[uint64]bool)

	return list.Find(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		hash := value.Hash()

		if _, exists := hashTable[hash]; !exists {
			hashTable[hash] = true

			return true, nil
		}

		return false, nil
	})
}
