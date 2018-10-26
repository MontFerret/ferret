package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

type (
	SortDirection int

	Comparator func(first DataSet, second DataSet) (int, error)

	Sorter struct {
		fn        Comparator
		direction SortDirection
	}

	SortIterator struct {
		src     Iterator
		sorters []*Sorter
		ready   bool
		values  []DataSet
		err     error
		pos     int
	}
)

const (
	SortDirectionAsc  SortDirection = 1
	SortDirectionDesc SortDirection = -1
)

func SortDirectionFromString(str string) SortDirection {
	if strings.ToUpper(str) == "DESC" {
		return SortDirectionDesc
	}

	return SortDirectionAsc
}

func IsValidSortDirection(direction SortDirection) bool {
	switch direction {
	case SortDirectionAsc, SortDirectionDesc:
		return true
	default:
		return false
	}
}

func NewSorter(fn Comparator, direction SortDirection) (*Sorter, error) {
	if fn == nil {
		return nil, core.Error(core.ErrMissedArgument, "fn")
	}

	if IsValidSortDirection(direction) == false {
		return nil, core.Error(core.ErrInvalidArgument, "direction")
	}

	return &Sorter{fn, direction}, nil
}

func NewSortIterator(
	src Iterator,
	comparators ...*Sorter,
) (*SortIterator, error) {
	if core.IsNil(src) {
		return nil, errors.Wrap(core.ErrMissedArgument, "source")
	}

	if comparators == nil || len(comparators) == 0 {
		return nil, errors.Wrap(core.ErrMissedArgument, "comparator")
	}

	return &SortIterator{
		src,
		comparators,
		false,
		nil, nil,
		0,
	}, nil
}

func (iterator *SortIterator) HasNext() bool {
	// we need to initialize the iterator
	if iterator.ready == false {
		iterator.ready = true
		sorted, err := iterator.sort()

		if err != nil {
			// dataSet to true because we do not want to initialize next time anymore
			iterator.values = nil
			iterator.err = err

			// if there is an error, we need to show it during Next()
			return true
		}

		iterator.values = sorted
	}

	return iterator.values != nil && len(iterator.values) > iterator.pos
}

func (iterator *SortIterator) Next() (DataSet, error) {
	if iterator.err != nil {
		return nil, iterator.err
	}

	if len(iterator.values) > iterator.pos {
		idx := iterator.pos
		val := iterator.values[idx]
		iterator.pos++

		return val, nil
	}

	return nil, ErrExhausted
}

func (iterator *SortIterator) sort() ([]DataSet, error) {
	res, err := ToSlice(iterator.src)

	if err != nil {
		return nil, err
	}

	var failure error

	sort.SliceStable(res, func(i, j int) bool {
		// ignore next execution
		if failure != nil {
			return false
		}

		var out bool

		for _, comp := range iterator.sorters {
			left := res[i]
			right := res[j]

			eq, err := comp.fn(left, right)

			if err != nil {
				failure = err
				out = false

				break
			}

			eq = eq * int(comp.direction)

			if eq == -1 {
				out = true
				break
			}

			if eq == 1 {
				out = false
				break
			}
		}

		return out
	})

	if failure != nil {
		return nil, failure
	}

	return res, nil
}
