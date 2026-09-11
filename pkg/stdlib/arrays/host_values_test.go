package arrays_test

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	sliceProbe struct {
		runtime.List
		called bool
	}

	failingList struct {
		*runtime.Array
		err error
	}

	distinctCollisionValue struct {
		label string
	}

	failingEqualityValue struct {
		err error
	}

	lengthFailingList struct {
		*runtime.Array
		err error
	}

	copyFailureList struct {
		runtime.List
		copied  runtime.Value
		copies  int
		lookups int
	}

	copyMutationProbe struct {
		runtime.Value
		mutations int
	}
)

func (p *sliceProbe) Slice(ctx context.Context, start, end runtime.Int) (runtime.List, error) {
	p.called = true
	if start < 0 || end < start {
		return nil, errors.New("invalid host slice")
	}

	return p.List.Slice(ctx, start, end)
}

func (l *failingList) ForEach(context.Context, runtime.IndexReadablePredicate) error {
	return l.err
}

func (l *failingList) IndexOf(context.Context, runtime.Value) (runtime.Int, error) {
	return -1, l.err
}

func (l *failingList) Length(context.Context) (runtime.Int, error) {
	return 0, l.err
}

func (l *failingList) Filter(context.Context, runtime.IndexReadablePredicate) (runtime.List, error) {
	return nil, l.err
}

func (l *failingList) Copy() runtime.Value {
	return &failingList{Array: l.Array.CopyWithGrowth(0), err: l.err}
}

func (l *failingList) Append(context.Context, runtime.Value) error {
	return l.err
}

func (v distinctCollisionValue) String() string {
	return v.label
}

func (v distinctCollisionValue) Hash() uint64 {
	return 7
}

func (v distinctCollisionValue) Copy() runtime.Value {
	return v
}

func (v distinctCollisionValue) Equal(_ context.Context, other runtime.Value) (bool, error) {
	o, ok := other.(distinctCollisionValue)
	if !ok {
		return false, nil
	}

	return v.label == o.label, nil
}

func (v failingEqualityValue) String() string {
	return "failing"
}

func (v failingEqualityValue) Hash() uint64 {
	return 11
}

func (v failingEqualityValue) Copy() runtime.Value {
	return v
}

func (v failingEqualityValue) Equal(ctx context.Context, _ runtime.Value) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	return false, v.err
}

func (l *lengthFailingList) Length(context.Context) (runtime.Int, error) {
	return 0, l.err
}

func (l *lengthFailingList) Copy() runtime.Value {
	return &lengthFailingList{Array: l.Array.CopyWithGrowth(0), err: l.err}
}

func (list *copyFailureList) Copy() runtime.Value {
	list.copies++

	return list.copied
}

func (list *copyFailureList) IndexOf(ctx context.Context, value runtime.Value) (runtime.Int, error) {
	list.lookups++

	return list.List.IndexOf(ctx, value)
}

func (probe *copyMutationProbe) Append(context.Context, runtime.Value) error {
	probe.mutations++

	return nil
}

func (probe *copyMutationProbe) RemoveAt(context.Context, runtime.Int) (runtime.Value, error) {
	probe.mutations++

	return runtime.None, nil
}
