package arrays

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const (
	iteratorValVar = "value"
	iteratorKeyVar = "key"
)

func NewLib() map[string]core.Function {
	return map[string]core.Function{
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
		"REVERSE":        Reverse,
		"SHIFT":          Shift,
		"SLICE":          Slice,
		"SORTED":         Sorted,
		"SORTED_UNIQUE":  SortedUnique,
		"UNION":          Union,
		"UNION_DISTINCT": UnionDistinct,
		"UNIQUE":         Unique,
		"UNSHIFT":        Unshift,
	}
}

func toArrayIterator(arr *values.Array) collections.Iterator {
	return collections.NewArrayIterator(
		iteratorValVar,
		iteratorKeyVar,
		arr,
	)
}

func toArray(iterator collections.Iterator) (core.Value, error) {
	arr := values.NewArray(10)

	for iterator.HasNext() {
		ds, err := iterator.Next()

		if err != nil {
			return values.None, err
		}

		arr.Push(ds.Get(iteratorValVar))
	}

	return arr, nil
}
