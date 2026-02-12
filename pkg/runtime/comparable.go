package runtime

type (
	// Comparable is an interface that can be implemented by values that can be compared to other values.
	Comparable interface {
		Compare(ctx Context, other Value) int64
	}

	// Comparator is a function that compares two values and returns an integer indicating their relative order.
	Comparator = func(ctx Context, first, second Value) int64
)

func CompareValues(ctx Context, a, b Value) int64 {
	aComparable, ok := a.(Comparable)

	if ok {
		return aComparable.Compare(ctx, b)
	}

	bComparable, ok := b.(Comparable)

	if ok {
		res := bComparable.Compare(ctx, a)

		return -res
	}

	return CompareTypes(a, b)
}
