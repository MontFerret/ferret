package objects_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
)

func BenchmarkConfiguredMapMerge(b *testing.B) {
	ctx := context.Background()
	for _, size := range []int{0, 1, 16, 1024} {
		for name, merge := range map[string]runtime.Function{"shallow": objects.Merge, "deep": objects.MergeDeep} {
			b.Run(fmt.Sprintf("%s/size=%d", name, size), func(b *testing.B) {
				values := make(map[string]runtime.Value, size)
				for i := range size {
					values[fmt.Sprint(i)] = runtime.Int(i)
				}

				if size > 0 {
					values["0"] = runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.Int(1)})
				}

				left := &configuredMap{Object: runtime.NewObjectWith(values), backend: "configured"}
				right := &configuredMap{Object: runtime.NewObjectWith(values), backend: "other"}
				b.ReportAllocs()
				for b.Loop() {
					result, err := merge(ctx, left, right)
					if err != nil {
						b.Fatal(err)
					}

					if err := result.(*configuredMap).Close(); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
