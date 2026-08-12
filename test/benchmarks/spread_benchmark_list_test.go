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
