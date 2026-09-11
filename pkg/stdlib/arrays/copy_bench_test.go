package arrays_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

var copyBenchmarkResult runtime.Value

func BenchmarkListCopyOperations(b *testing.B) {
	for _, operation := range []struct {
		run  func(context.Context, runtime.Value) (runtime.Value, error)
		name string
	}{
		{name: "ToList", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return runtime.ToList(ctx, value)
		}},
		{name: "Append", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.Append(ctx, value, runtime.Int(1))
		}},
		{name: "RemoveAt", run: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return arrays.RemoveAt(ctx, value, runtime.Int(0))
		}},
	} {
		for _, size := range []int{0, 16, 1024} {
			for _, host := range []bool{false, true} {
				kind := "Native"
				if host {
					kind = "Host"
				}

				b.Run(fmt.Sprintf("%s/%s/%d", operation.name, kind, size), func(b *testing.B) {
					array := runtime.NewSizedArray(size)
					var input runtime.Value = array
					if host {
						input = &copyBenchmarkList{List: array}
					}

					b.ReportAllocs()
					b.ResetTimer()
					for b.Loop() {
						result, err := operation.run(b.Context(), input)
						if err != nil {
							b.Fatal(err)
						}

						copyBenchmarkResult = result
					}
				})
			}
		}
	}
}
