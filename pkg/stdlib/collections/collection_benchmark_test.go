package collections_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
)

func BenchmarkCollections(b *testing.B) {
	ctx := context.Background()
	array := runtime.NewArray(1024)
	for i := range 1024 {
		if err := array.Append(ctx, runtime.Int(i%128)); err != nil {
			b.Fatal(err)
		}
	}

	host := &benchmarkList{Array: array, backend: "memory"}
	stream := &benchmarkIterable{values: array}
	for _, tc := range []struct {
		input runtime.Value
		call  func(context.Context, runtime.Value) (runtime.Value, error)
		name  string
	}{
		{name: "CountNative", input: array, call: collections.Count},
		{name: "CountHost", input: host, call: collections.Count},
		{name: "DistinctNative", input: array, call: collections.CountDistinct},
		{name: "DistinctHost", input: host, call: collections.CountDistinct},
		{name: "ReverseNative", input: array, call: collections.Reverse},
		{name: "ReverseHost", input: host, call: collections.Reverse},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				if _, err := tc.call(ctx, tc.input); err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	for _, tc := range []struct {
		input  runtime.Value
		needle runtime.Value
		name   string
	}{
		{name: "IncludesContainable", input: host, needle: runtime.Int(-1)},
		{name: "IncludesScan", input: stream, needle: runtime.Int(-1)},
		{name: "IncludesEarly", input: stream, needle: runtime.Int(0)},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				if _, err := collections.Includes(ctx, tc.input, tc.needle); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// These inputs were previously rejected, so they have no successful old baseline.
func BenchmarkCollectionIterableCounting(b *testing.B) {
	array := runtime.NewArray(1024)
	for i := range 1024 {
		if err := array.Append(b.Context(), runtime.Int(i%128)); err != nil {
			b.Fatal(err)
		}
	}

	stream := &benchmarkIterable{values: array}
	for name, call := range map[string]func(context.Context, runtime.Value) (runtime.Value, error){"Count": collections.Count, "Distinct": collections.CountDistinct} {
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				if _, err := call(context.Background(), stream); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
