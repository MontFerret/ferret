package session

import (
	"context"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func BenchmarkCodecAdoption(b *testing.B) {
	for _, codec := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		b.Run(codec.ContentType(), func(b *testing.B) {
			for _, depth := range []int{0, 8} {
				name := "resources"
				if depth > 0 {
					name = "containers_depth_8"
				}

				b.Run(name, func(b *testing.B) {
					items := make([]runtime.Value, 1024)
					for i := range items {
						var value runtime.Value = &adoptionBenchResource{Value: runtime.String("value")}
						for range depth {
							value = runtime.NewArrayWith(value)
						}

						items[i] = value
					}

					value := &adoptionBenchIterable{items: items}
					for _, control := range []string{"absent", "noop", "result"} {
						b.Run(control, func(b *testing.B) {
							ctx := context.Background()
							if control == "noop" {
								ctx = encodingownership.WithValueAdopter(ctx, func(runtime.Value, bool) {})
							}

							b.ReportAllocs()
							b.ResetTimer()
							for range b.N {
								var result vm.Result
								opCtx := ctx
								if control == "result" {
									adopter := outputValueAdopter{result: &result}
									opCtx = encodingownership.WithValueAdopter(ctx, adopter.Adopt)
								}

								if err := codec.Encode(opCtx, io.Discard, value); err != nil {
									b.Fatal(err)
								}

								if err := result.Close(); err != nil {
									b.Fatal(err)
								}
							}
						})
					}
				})
			}
		})
	}
}
