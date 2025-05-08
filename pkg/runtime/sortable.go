package runtime

import (
	"context"
	"sort"
)

// Sortable is an interface that defines methods for sorting a collection of values.
type Sortable interface {
	// SortAsc sorts the collection in ascending order.
	SortAsc(context.Context) error

	// SortDesc sorts the collection in descending order.
	SortDesc(context.Context) error

	// SortWith sorts the collection using a custom comparator.
	SortWith(context.Context, Comparator) error
}

func SortSlice(values []Value, ascending Boolean) {
	var pivot int64 = -1

	if ascending {
		pivot = 1
	}

	SortSliceWith(values, func(first, second Value) int64 {
		comp := CompareValues(first, second)

		return pivot * comp
	})
}

func SortSliceWith(values []Value, comparator Comparator) {
	sort.SliceStable(values, func(i, j int) bool {
		comp := comparator(values[i], values[j])

		return comp == -1
	})
}
