package runtime

import (
	"context"
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

// Sort is a generic sorting function that accepts a Sortable value.
func Sort(ctx context.Context, values Value, ascending Boolean) error {
	switch value := values.(type) {
	case Sortable:
		if ascending {
			return value.SortAsc(ctx)
		}

		return value.SortDesc(ctx)
	default:
		return TypeErrorOf(values, TypeSortable)
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
	size, err := values.Length(ctx)

	if err != nil {
		return err
	}

	return stableSort(ctx, func(ctx context.Context, a, b Int) (Boolean, error) {
		aVal, err := values.At(ctx, a)

		if err != nil {
			return false, err
		}

		bVal, err := values.At(ctx, b)

		if err != nil {
			return false, err
		}

		comparison, err := CompareValues(ctx, aVal, bVal)
		if err != nil {
			return false, err
		}

		if !ascending {
			comparison = reverseOrdering(comparison)
		}

		return comparison < Equal, nil
	}, values.Swap, size)
}

// SortListWith sorts the given List using the stable Sort algorithm using a custom comparator
func SortListWith(ctx context.Context, values List, comparator Comparator) error {
	size, err := values.Length(ctx)

	if err != nil {
		return err
	}

	return stableSort(ctx, func(ctx context.Context, a, b Int) (Boolean, error) {
		aVal, err := values.At(ctx, a)

		if err != nil {
			return false, err
		}

		bVal, err := values.At(ctx, b)

		if err != nil {
			return false, err
		}

		comparison, err := comparator(ctx, aVal, bVal)
		if err != nil {
			return false, err
		}

		return comparison < Equal, nil
	}, values.Swap, size)
}

func SortSlice(ctx context.Context, values []Value, ascending Boolean) error {
	return SortSliceWith(ctx, values, func(ctx context.Context, first, second Value) (Ordering, error) {
		comparison, err := CompareValues(ctx, first, second)
		if err != nil {
			return Equal, err
		}

		if !ascending {
			comparison = reverseOrdering(comparison)
		}

		return comparison, nil
	})
}

func SortSliceWith(ctx context.Context, values []Value, comparator Comparator) error {
	return stableSort(ctx, func(ctx context.Context, first, second Int) (Boolean, error) {
		comparison, err := comparator(ctx, values[first], values[second])
		if err != nil {
			return false, err
		}

		return comparison < Equal, nil
	}, func(_ context.Context, first, second Int) error {
		values[first], values[second] = values[second], values[first]
		return nil
	}, Int(len(values)))
}
