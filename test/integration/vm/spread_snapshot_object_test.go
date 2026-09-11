package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type spreadSnapshotObject struct {
	runtime.Map
	snapshot      *runtime.Object
	snapshotErr   error
	fallbackErr   error
	snapshotCalls int
	fallbackCalls int
}

func (s *spreadSnapshotObject) Snapshot(context.Context) (*runtime.Object, error) {
	s.snapshotCalls++

	return s.snapshot, s.snapshotErr
}

func (s *spreadSnapshotObject) ForEach(ctx context.Context, predicate runtime.KeyReadablePredicate) error {
	s.fallbackCalls++
	if s.fallbackErr != nil {
		return s.fallbackErr
	}

	return s.Map.ForEach(ctx, predicate)
}

func (s *spreadSnapshotObject) New(ctx context.Context) (runtime.Map, error) {
	base, err := s.Map.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *s
	result.Map = base
	result.snapshot = base.(*runtime.Object)

	return &result, nil
}
