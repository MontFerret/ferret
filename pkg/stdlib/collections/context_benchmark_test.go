package collections_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
)

// Reused contexts exclude setup and first Done access. Fresh contexts include
// creation, the collection call, and cancellation, including lazy channel costs.
func BenchmarkCollectionContexts(b *testing.B) {
	for _, size := range []int{0, 1, 16, 1024} {
		array := runtime.NewArray(size)
		for i := range size {
			if err := array.Append(context.Background(), runtime.Int(i%128)); err != nil {
				b.Fatal(err)
			}
		}

		host := &benchmarkList{Array: array, backend: "memory"}
		stream := &benchmarkIterable{values: array}
		cases := []struct {
			input runtime.Value
			call  func(context.Context, runtime.Value) (runtime.Value, error)
			name  string
		}{
			{name: "CountNative", input: array, call: collections.Count},
			{name: "CountHost", input: host, call: collections.Count},
			{name: "CountScan", input: stream, call: collections.Count},
			{name: "DistinctNative", input: array, call: collections.CountDistinct},
			{name: "DistinctHost", input: host, call: collections.CountDistinct},
			{name: "DistinctScan", input: stream, call: collections.CountDistinct},
			{name: "ReverseNative", input: array, call: collections.Reverse},
			{name: "ReverseHost", input: host, call: collections.Reverse},
			{name: "IncludesContainable", input: host, call: benchmarkCollectionMissing},
			{name: "IncludesScan", input: stream, call: benchmarkCollectionMissing},
			{name: "IncludesEarly", input: stream, call: benchmarkCollectionFirst},
		}

		for _, kind := range []string{"Background", "Cancellable", "Deadline", "Layered"} {
			for _, lifetime := range []string{"Reused", "Fresh"} {
				for _, tc := range cases {
					b.Run(fmt.Sprintf("Size=%d/%s/%s/%s", size, kind, lifetime, tc.name), func(b *testing.B) {
						if lifetime == "Fresh" {
							for b.Loop() {
								ctx, cancel := collectionBenchmarkContext(kind)
								_, err := tc.call(ctx, tc.input)
								cancel()
								if err != nil {
									b.Fatal(err)
								}
							}

							return
						}

						ctx, cancel := collectionBenchmarkContext(kind)
						defer cancel()
						_ = ctx.Done()
						for b.Loop() {
							if _, err := tc.call(ctx, tc.input); err != nil {
								b.Fatal(err)
							}
						}
					})
				}
			}
		}
	}
}

func collectionBenchmarkContext(kind string) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	switch kind {
	case "Cancellable":
		return context.WithCancel(ctx)
	case "Deadline":
		return context.WithDeadline(ctx, time.Now().Add(time.Hour))
	case "Layered":
		type key int

		ctx, cancel := context.WithCancel(ctx)
		for i := range 4 {
			ctx = context.WithValue(ctx, key(i), i)
		}

		return ctx, cancel
	default:
		return ctx, func() {}
	}
}

func benchmarkCollectionMissing(ctx context.Context, source runtime.Value) (runtime.Value, error) {
	return collections.Includes(ctx, source, runtime.Int(-1))
}

func benchmarkCollectionFirst(ctx context.Context, source runtime.Value) (runtime.Value, error) {
	return collections.Includes(ctx, source, runtime.Int(0))
}
