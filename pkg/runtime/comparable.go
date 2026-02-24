package runtime

type (
	Comparable interface {
		Compare(other Value) int
	}

	Comparator = func(first, second Value) int
)

func CompareValues(a, b Value) int {
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
