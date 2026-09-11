package benchmarks_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type benchmarkSnapshotList struct {
	runtime.List
	snapshot *runtime.Array
}

func (b *benchmarkSnapshotList) Snapshot(context.Context) (*runtime.Array, error) {
	return b.snapshot, nil
}

func (b *benchmarkSnapshotList) New(ctx context.Context) (runtime.List, error) {
	base, err := b.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *b
	result.List = base
	result.snapshot = base.(*runtime.Array)

	return &result, nil
}
