package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

type (
	SortDirection int

	Comparator func(first core.Value, second core.Value) (int, error)

	Sorter struct {
		fn        Comparator
		direction SortDirection
	}

	SortIterator struct {
		src     Iterator
		sorters []*Sorter
		ready   bool
		values  Iterator
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

	return &SortIterator{src, comparators, false, nil}, nil
}

func (iterator *SortIterator) HasNext() bool {
	// we need to initialize the iterator
	if iterator.ready == false {
		iterator.ready = true
		values, err := iterator.sort()

		if err != nil {
			// set to true because we do not want to initialize next time anymore
			iterator.values = NoopIterator

			return false
		}

		iterator.values = values
	}

	return iterator.values.HasNext()
}

func (iterator *SortIterator) Next() (core.Value, core.Value, error) {
	return iterator.values.Next()
}

func (iterator *SortIterator) sort() (Iterator, error) {
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

	return NewSliceIterator(res), nil
}
