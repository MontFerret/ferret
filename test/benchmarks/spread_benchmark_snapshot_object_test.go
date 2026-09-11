package benchmarks_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type benchmarkSnapshotObject struct {
	runtime.Map
	snapshot *runtime.Object
}

func (b *benchmarkSnapshotObject) Snapshot(context.Context) (*runtime.Object, error) {
	return b.snapshot, nil
}

func (b *benchmarkSnapshotObject) New(ctx context.Context) (runtime.Map, error) {
	base, err := b.Map.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *b
	result.Map = base
	result.snapshot = base.(*runtime.Object)

	return &result, nil
}
