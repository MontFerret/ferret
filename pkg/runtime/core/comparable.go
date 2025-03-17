package core

type (
	Comparable interface {
		Compare(other Value) int64
	}

	Comparator = func(first, second Value) int64
)

func CompareValues(a, b Value) int64 {
	aComparable, ok := a.(Comparable)

	if ok {
		return aComparable.Compare(b)
	}

	bComparable, ok := b.(Comparable)

	if ok {
		res := bComparable.Compare(a)

		return -res
	}

	return CompareTypes(a, b)
}
