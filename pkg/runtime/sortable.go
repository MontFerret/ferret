package runtime

import (
	"context"
	"sort"
)

type Sortable interface {
	SortAsc(context.Context) (List, error)
	SortDesc(context.Context) (List, error)
	SortWith(context.Context, Comparator) (List, error)
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

		return comp == 0
	})
}
