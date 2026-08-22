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
