package testing

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	snapshotList struct {
		runtime.List
		snapshot      *runtime.Array
		snapshotErr   error
		snapshotCalls int
	}

	snapshotMap struct {
		runtime.Map
		snapshot      *runtime.Object
		snapshotErr   error
		snapshotCalls int
	}

	unsafeList struct {
		runtime.List
		atCalls int
	}

	unsafeMap struct {
		runtime.Map
		forEachCalls int
	}

	equalityErrorValue struct {
		err error
	}

	diagnosticEqualityProbe struct {
		calls *int
		err   error
	}
)

func (s *snapshotList) Snapshot(context.Context) (*runtime.Array, error) {
	s.snapshotCalls++

	return s.snapshot, s.snapshotErr
}

func (s *snapshotMap) Snapshot(context.Context) (*runtime.Object, error) {
	s.snapshotCalls++

	return s.snapshot, s.snapshotErr
}

func (l *unsafeList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	l.atCalls++

	return l.List.At(ctx, index)
}

func (m *unsafeMap) ForEach(ctx context.Context, predicate runtime.KeyReadablePredicate) error {
	m.forEachCalls++

	return m.Map.ForEach(ctx, predicate)
}

func (v *equalityErrorValue) String() string {
	return "error-value"
}

func (v *equalityErrorValue) Hash() uint64 {
	return 1
}

func (v *equalityErrorValue) Copy() runtime.Value {
	return v
}

func (v *equalityErrorValue) Equal(context.Context, runtime.Value) (bool, error) {
	return false, v.err
}

func (v *diagnosticEqualityProbe) String() string {
	return "diagnostic-equality-probe"
}

func (v *diagnosticEqualityProbe) Hash() uint64 {
	return 1
}

func (v *diagnosticEqualityProbe) Copy() runtime.Value {
	return v
}

func (v *diagnosticEqualityProbe) Equal(context.Context, runtime.Value) (bool, error) {
	if v.calls != nil {
		(*v.calls)++
	}

	return false, v.err
}

func (s *snapshotList) New(ctx context.Context) (runtime.List, error) {
	base, err := s.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *s
	result.List = base
	result.snapshot = base.(*runtime.Array)

	return &result, nil
}

func (s *snapshotMap) New(ctx context.Context) (runtime.Map, error) {
	base, err := s.Map.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *s
	result.Map = base
	result.snapshot = base.(*runtime.Object)

	return &result, nil
}

func (l *unsafeList) New(ctx context.Context) (runtime.List, error) {
	base, err := l.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *l
	result.List = base

	return &result, nil
}

func (m *unsafeMap) New(ctx context.Context) (runtime.Map, error) {
	base, err := m.Map.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *m
	result.Map = base

	return &result, nil
}
