package arrays_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

func BenchmarkArrayTransforms(b *testing.B) {
	for _, name := range []string{"RemoveNone", "RemoveHalf", "RemoveAll", "Sorted"} {
		for _, size := range []int{16, 1024} {
			for _, host := range []bool{false, true} {
				b.Run(fmt.Sprintf("%s/size=%d/host=%v", name, size, host), func(b *testing.B) {
					input := transformBenchmarkInput(b.Context(), name, size, host)
					b.ReportAllocs()
					b.ResetTimer()
					for b.Loop() {
						var err error
						if name == "Sorted" {
							copyBenchmarkResult, err = arrays.Sorted(b.Context(), input)
						} else {
							copyBenchmarkResult, err = arrays.Remove(b.Context(), input, runtime.Int(0))
						}

						if err != nil {
							b.Fatal(err)
						}
					}
				})
			}
		}
	}
}

func BenchmarkLegacyArrayRemoval(b *testing.B) {
	remove, ok := arrayFunctions(b).A3().Get("remove_value")
	if !ok {
		b.Fatal("legacy removal is not registered")
	}

	for _, name := range []string{"RemoveNone", "RemoveHalf", "RemoveAll"} {
		for _, size := range []int{16, 1024} {
			for _, host := range []bool{false, true} {
				for _, limit := range []runtime.Int{-1, 0, 1} {
					b.Run(fmt.Sprintf("%s/size=%d/host=%v/limit=%d", name, size, host, limit), func(b *testing.B) {
						input := transformBenchmarkInput(b.Context(), name, size, host)
						b.ReportAllocs()
						b.ResetTimer()
						for b.Loop() {
							var err error
							copyBenchmarkResult, err = remove(b.Context(), input, runtime.Int(0), limit)
							if err != nil {
								b.Fatal(err)
							}
						}
					})
				}
			}
		}
	}
}

func transformBenchmarkInput(ctx context.Context, name string, size int, host bool) runtime.List {
	array := runtime.NewArray(size)
	for index := range size {
		value := runtime.Int(size - index)
		switch name {
		case "RemoveHalf":
			value = runtime.Int(index % 2)
		case "RemoveAll":
			value = 0
		}

		_ = array.Append(ctx, value)
	}

	if host {
		return &copyBenchmarkList{List: array}
	}

	return array
}
