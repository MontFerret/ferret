package runtime

import (
	"sort"
)

// Sortable is an interface that defines methods for sorting a collection of values.
type Sortable interface {
	// SortAsc sorts the collection in ascending order.
	SortAsc(Context) error

	// SortDesc sorts the collection in descending order.
	SortDesc(Context) error
}

// SortAsc sorts the given values in ascending order. It accepts either a List or value that implements Sortable interface.
func SortAsc(ctx Context, values Value) error {
	return Sort(ctx, values, true)
}

// SortDesc sorts the given values in descending order. It accepts either a List or value that implements Sortable interface.
func SortDesc(ctx Context, values Value) error {
	return Sort(ctx, values, false)
}

// Sort is a generic sorting function that accepts either a List or value that implements Sortable interface.
func Sort(ctx Context, values Value, ascending Boolean) error {
	switch value := values.(type) {
	case Sortable:
		if ascending {
			return value.SortAsc(ctx)
		}

		return value.SortDesc(ctx)
	case List:
		return SortList(ctx, value, ascending)
	default:
		return TypeErrorOf(values, TypeList, TypeSortable)
	}
}

// SortListAsc sorts the given List in ascending order using a stable sorting algorithm.
func SortListAsc(ctx Context, values List) error {
	return SortList(ctx, values, true)
}

// SortListDesc sorts the given List in descending order using a stable sorting algorithm.
func SortListDesc(ctx Context, values List) error {
	return SortList(ctx, values, false)
}

// SortList sorts the given List using the stable Sort algorithm
func SortList(ctx Context, values List, ascending Boolean) error {
	var pivot int64 = -1

	if ascending {
		pivot = 1
	}

	size, err := values.Length(ctx)

	if err != nil {
		return err
	}

	return stableSort(ctx, func(ctx Context, a, b Int) (Boolean, error) {
		aVal, err := values.Get(ctx, a)

		if err != nil {
			return false, err
		}

		bVal, err := values.Get(ctx, b)

		if err != nil {
			return false, err
		}

		comp := CompareValues(ctx, aVal, bVal) * pivot

		return comp == -1, nil
	}, values.Swap, size)
}

// SortListWith sorts the given List using the stable Sort algorithm using a custom comparator
func SortListWith(ctx Context, values List, comparator Comparator) error {
	size, err := values.Length(ctx)

	if err != nil {
		return err
	}

	return stableSort(ctx, func(ctx Context, a, b Int) (Boolean, error) {
		aVal, err := values.Get(ctx, a)

		if err != nil {
			return false, err
		}

		bVal, err := values.Get(ctx, b)

		if err != nil {
			return false, err
		}

		return comparator(ctx, aVal, bVal) == -1, nil
	}, values.Swap, size)
}

// SortSlice sorts a slice of Values in either ascending or descending order using a stable sorting algorithm.
func SortSlice(ctx Context, values []Value, ascending Boolean) {
	var pivot int64 = -1

	if ascending {
		pivot = 1
	}

	SortSliceWith(ctx, values, func(_ Context, first, second Value) int64 {
		comp := CompareValues(nil, first, second)

		return pivot * comp
	})
}

// SortSliceWith sorts a slice of Values using a custom comparator and a stable sorting algorithm.
func SortSliceWith(ctx Context, values []Value, comparator Comparator) {
	sort.SliceStable(values, func(i, j int) bool {
		comp := comparator(ctx, values[i], values[j])

		return comp == -1
	})
}
