package stdlib_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/random"
	stdlibstrings "github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func BenchmarkStdlibCancellation(b *testing.B) {
	for _, size := range []int{1, 16, 4096} {
		values := runtime.NewArray(size)
		texts := runtime.NewArray(size)
		object := runtime.NewObject()
		for index := range size {
			_ = values.Append(b.Context(), runtime.Int((index*37)%size))
			_ = texts.Append(b.Context(), runtime.String("value"))
			_ = object.Set(b.Context(), runtime.String(fmt.Sprint(index)), runtime.Int(index))
		}

		stream := &cancellationBenchmarkIterable{Value: runtime.None, values: values}
		for _, kind := range []string{"Background", "Cancellable", "Layered"} {
			ctx := context.Background()
			cancel := func() {}
			if kind != "Background" {
				ctx, cancel = context.WithCancel(ctx)
			}

			if kind == "Layered" {
				type key int
				for index := range 4 {
					ctx = context.WithValue(ctx, key(index), index)
				}
			}

			ctx = rnd.WithContext(ctx, rnd.NewSeed(42))
			cases := map[string]func() (runtime.Value, error){
				"CountMeasured": func() (runtime.Value, error) { return collections.Count(ctx, values) },
				"CountScan":     func() (runtime.Value, error) { return collections.Count(ctx, stream) },
				"Distinct":      func() (runtime.Value, error) { return collections.CountDistinct(ctx, values) },
				"IncludesScan":  func() (runtime.Value, error) { return collections.Includes(ctx, stream, runtime.Int(-1)) },
				"IncludesString": func() (runtime.Value, error) {
					return collections.Includes(ctx, runtime.String("value"), runtime.String("v"))
				},
				"Reverse":     func() (runtime.Value, error) { return collections.Reverse(ctx, values) },
				"HasKey":      func() (runtime.Value, error) { return objects.HasKey(ctx, object, runtime.String("0")) },
				"Entries":     func() (runtime.Value, error) { return objects.Entries(ctx, object) },
				"ArraySet":    func() (runtime.Value, error) { return arrays.SetMutable(ctx, values, runtime.Int(0), runtime.Int(0)) },
				"ArrayRemove": func() (runtime.Value, error) { return arrays.RemoveMutable(ctx, values.Copy(), runtime.Int(0)) },
				"Sum":         func() (runtime.Value, error) { return math.Sum(ctx, values) },
				"Median":      func() (runtime.Value, error) { return math.Median(ctx, values) },
				"Choice":      func() (runtime.Value, error) { return random.Choice(ctx, values) },
				"Shuffle":     func() (runtime.Value, error) { return random.Shuffle(ctx, values) },
				"Join":        func() (runtime.Value, error) { return stdlibstrings.Join(ctx, texts, runtime.String(",")) },
			}
			for name, call := range cases {
				b.Run(fmt.Sprintf("%s/Size=%d/%s", name, size, kind), func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						if _, err := call(); err != nil {
							b.Fatal(err)
						}
					}
				})
			}

			cancel()
		}
	}
}
