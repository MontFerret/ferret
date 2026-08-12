package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type spreadSnapshotList struct {
	runtime.List
	snapshot      *runtime.Array
	snapshotErr   error
	fallbackErr   error
	snapshotCalls int
	fallbackCalls int
}

func (s *spreadSnapshotList) Snapshot(context.Context) (*runtime.Array, error) {
	s.snapshotCalls++

	return s.snapshot, s.snapshotErr
}

func (s *spreadSnapshotList) Iterate(ctx context.Context) (runtime.Iterator, error) {
	s.fallbackCalls++
	if s.fallbackErr != nil {
		return nil, s.fallbackErr
	}

	return s.List.Iterate(ctx)
}
