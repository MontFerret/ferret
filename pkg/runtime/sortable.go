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
}

func SortAsc(ctx context.Context, values Value) error {
	return Sort(ctx, values, true)
}

func SortDesc(ctx context.Context, values Value) error {
	return Sort(ctx, values, false)
}

// Sort is a generic sorting function that accepts either a List or value that implements Sortable interface.
func Sort(ctx context.Context, values Value, ascending Boolean) error {
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

func SortListAsc(ctx context.Context, values List) error {
	return SortList(ctx, values, true)
}

func SortListDesc(ctx context.Context, values List) error {
	return SortList(ctx, values, false)
}

// SortList sorts the given List using the stable Sort algorithm
func SortList(ctx context.Context, values List, ascending Boolean) error {
	var pivot int64 = -1

	if ascending {
		pivot = 1
	}

	size, err := values.Length(ctx)

	if err != nil {
		return err
	}

	return stableSort(ctx, func(ctx context.Context, a, b Int) (Boolean, error) {
		aVal, err := values.Get(ctx, a)

		if err != nil {
			return false, err
		}

		bVal, err := values.Get(ctx, b)

		if err != nil {
			return false, err
		}

		comp := CompareValues(aVal, bVal) * pivot

		return comp == -1, nil
	}, values.Swap, size)
}

// SortListWith sorts the given List using the stable Sort algorithm using a custom comparator
func SortListWith(ctx context.Context, values List, comparator Comparator) error {
	size, err := values.Length(ctx)

	if err != nil {
		return err
	}

	return stableSort(ctx, func(ctx context.Context, a, b Int) (Boolean, error) {
		aVal, err := values.Get(ctx, a)

		if err != nil {
			return false, err
		}

		bVal, err := values.Get(ctx, b)

		if err != nil {
			return false, err
		}

		return comparator(aVal, bVal) == -1, nil
	}, values.Swap, size)
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
