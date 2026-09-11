package arrays_test

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	removalComparisonValue struct {
		target          runtime.Value
		expectedContext context.Context
		failure         error
		seen            *[]int
		index           int
		equal           bool
	}
	removalFilterList struct {
		runtime.List
		expectedContext context.Context
		filters         int
	}
)

func (v *removalComparisonValue) String() string {
	return "removal comparison probe"
}

func (v *removalComparisonValue) Hash() uint64 {
	return 1
}

func (v *removalComparisonValue) Copy() runtime.Value {
	return v
}

func (v *removalComparisonValue) Equal(ctx context.Context, other runtime.Value) (bool, error) {
	if v.seen == nil || other != v.target || ctx != v.expectedContext {
		return false, errors.New("comparison arguments changed")
	}

	*v.seen = append(*v.seen, v.index)

	return v.equal, v.failure
}

func (l *removalFilterList) Filter(ctx context.Context, predicate runtime.IndexReadablePredicate) (runtime.List, error) {
	l.filters++
	if ctx != l.expectedContext {
		return nil, errors.New("filter context changed")
	}

	return l.List.Filter(ctx, predicate)
}

func (l *removalFilterList) New(ctx context.Context) (runtime.List, error) {
	base, err := l.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *l
	result.List = base

	return &result, nil
}
