package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
)

func spreadArray(ctx context.Context, destination *runtime.Array, source runtime.Value) error {
	if source == runtime.None {
		return nil
	}

	array, ok := source.(*runtime.Array)
	if ok {
		return destination.Concat(ctx, array)
	}

	list, ok := source.(runtime.List)
	if !ok {
		return diagnostics.SpreadErrorOf(source, runtime.TypeArray, runtime.TypeList)
	}

	snapshotter, ok := source.(runtime.ListSnapshotter)
	if !ok {
		return destination.Concat(ctx, list)
	}

	snapshot, err := snapshotter.Snapshot(ctx)
	if err != nil {
		return err
	}

	if snapshot == nil {
		return runtime.Error(runtime.ErrInvalidOperation, "ListSnapshotter returned nil Array")
	}

	return destination.Concat(ctx, snapshot)
}

func spreadObject(ctx context.Context, destination *data.FastObject, source runtime.Value) error {
	if source == runtime.None {
		return nil
	}

	object, ok := source.(*runtime.Object)
	if ok {
		return destination.Merge(ctx, object)
	}

	mapValue, ok := source.(runtime.Map)
	if !ok {
		return diagnostics.SpreadErrorOf(source, runtime.TypeObject, runtime.TypeMap)
	}

	snapshotter, ok := source.(runtime.MapSnapshotter)
	if !ok {
		return destination.Merge(ctx, mapValue)
	}

	snapshot, err := snapshotter.Snapshot(ctx)
	if err != nil {
		return err
	}

	if snapshot == nil {
		return runtime.Error(runtime.ErrInvalidOperation, "MapSnapshotter returned nil Object")
	}

	return destination.Merge(ctx, snapshot)
}
