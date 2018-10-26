package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

type (
	SortDirection int

	Comparator func(ctx context.Context, scope *core.Scope, first DataSet, second DataSet) (int, error)

	Sorter struct {
		fn        Comparator
		direction SortDirection
	}

	SortIterator struct {
		values  Iterator
		sorters []*Sorter
		ready   bool
		result  []DataSet
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
	values Iterator,
	comparators ...*Sorter,
) (*SortIterator, error) {
	if values == nil {
		return nil, errors.Wrap(core.ErrMissedArgument, "values")
	}

	if comparators == nil || len(comparators) == 0 {
		return nil, errors.Wrap(core.ErrMissedArgument, "comparator")
	}

	return &SortIterator{
		values,
		comparators,
		false,
		nil,
		0,
	}, nil
}

func (iterator *SortIterator) Next(ctx context.Context, scope *core.Scope) (DataSet, error) {
	// we need to initialize the iterator
	if iterator.ready == false {
		iterator.ready = true
		sorted, err := iterator.sort(ctx, scope)

		if err != nil {
			return nil, err
		}

		iterator.result = sorted
	}

	if len(iterator.result) > iterator.pos {
		idx := iterator.pos
		val := iterator.result[idx]
		iterator.pos++

		return val, nil
	}

	return nil, nil
}

func (iterator *SortIterator) sort(ctx context.Context, scope *core.Scope) ([]DataSet, error) {
	res, err := ToSlice(ctx, scope, iterator.values)

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

			eq, err := comp.fn(ctx, scope, left, right)

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
