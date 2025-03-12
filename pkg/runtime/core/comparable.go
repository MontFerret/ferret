package core

import "context"

type (
	Comparable interface {
		Compare(ctx context.Context, other Value) (int64, error)
	}

	Comparator = func(ctx context.Context, first, second Value) (int64, error)
)

func CompareValues(ctx context.Context, a, b Value) (int64, error) {
	aComparable, ok := a.(Comparable)

	if ok {
		return aComparable.Compare(ctx, b)
	}

	bComparable, ok := b.(Comparable)

	if ok {
		res, err := bComparable.Compare(ctx, a)

		if err != nil {
			return 0, err
		}

		return -res, nil
	}

	return CompareTypes(a, b), nil
}
