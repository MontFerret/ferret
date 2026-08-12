package benchmarks_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type benchmarkSnapshotObject struct {
	runtime.Map
	snapshot *runtime.Object
}

func (*benchmarkSnapshotObject) ObjectLike() {}

func (b *benchmarkSnapshotObject) Snapshot(context.Context) (*runtime.Object, error) {
	return b.snapshot, nil
}
