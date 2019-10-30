package arrays

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
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

func ToUniqueArray(arr *values.Array) *values.Array {
	hashTable := make(map[uint64]bool)
	result := values.NewArray(int(arr.Length()))

	arr.ForEach(func(item core.Value, _ int) bool {
		h := item.Hash()

		_, exists := hashTable[h]

		if !exists {
			hashTable[h] = true
			result.Push(item)
		}

		return true
	})

	return result
}
