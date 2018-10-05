package arrays

import "github.com/MontFerret/ferret/pkg/runtime/core"

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
