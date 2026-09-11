package arrays_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

func BenchmarkArrayMutations(b *testing.B) {
	for _, operation := range []struct {
		run  func(context.Context, runtime.Value) (runtime.Value, error)
		name string
	}{
		{name: "Push", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.PushMutable(ctx, value, runtime.Int(0))
		}},
		{name: "Pop", run: arrays.PopMutable},
		{name: "Shift", run: arrays.ShiftMutable},
		{name: "Unshift", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.UnshiftMutable(ctx, value, runtime.Int(0))
		}},
		{name: "Set", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.SetMutable(ctx, value, runtime.Int(0), runtime.Int(0))
		}},
		{name: "Insert", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.InsertMutable(ctx, value, runtime.Int(0), runtime.Int(0))
		}},
		{name: "RemoveAt", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.RemoveAtMutable(ctx, value, runtime.Int(0))
		}},
		{name: "RemoveNone", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.RemoveMutable(ctx, value, runtime.Int(0))
		}},
		{name: "RemoveHalf", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.RemoveMutable(ctx, value, runtime.Int(0))
		}},
		{name: "RemoveAll", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.RemoveMutable(ctx, value, runtime.Int(0))
		}},
		{name: "Clear", run: arrays.ClearMutable},
		{name: "Sorted", run: arrays.SortMutable},
	} {
		for _, size := range []int{16, 1024} {
			for _, host := range []bool{false, true} {
				b.Run(fmt.Sprintf("%s/size=%d/host=%v", operation.name, size, host), func(b *testing.B) {
					b.ReportAllocs()
					// Replenish targets outside the timed region in batches so shrinking or
					// sorting a target repeatedly cannot turn this into an empty/no-op benchmark.
					b.ResetTimer()
					for completed := 0; completed < b.N; {
						b.StopTimer()
						count := min(128, b.N-completed)
						targets := make([]runtime.List, count)
						for index := range targets {
							targets[index] = transformBenchmarkInput(b.Context(), operation.name, size, host)
						}
						b.StartTimer()

						for _, target := range targets {
							var err error
							copyBenchmarkResult, err = operation.run(b.Context(), target)
							if err != nil {
								b.Fatal(err)
							}
						}

						completed += count
					}
				})
			}
		}
	}
}
