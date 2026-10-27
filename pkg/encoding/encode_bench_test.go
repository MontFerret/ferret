package encoding_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var benchmarkCodecEncoded []byte

// BenchmarkCodecEncodeBytes separates per-call codec boxing from collecting
// output. All modes use EncodeBytes; only the conversion to Encoder changes.
func BenchmarkCodecEncodeBytes(b *testing.B) {
	booleans, nulls := make([]runtime.Value, 1024), make([]runtime.Value, 1024)
	for i := range booleans {
		booleans[i] = runtime.Boolean(i%2 == 0)
		nulls[i] = runtime.None
	}

	jsonInterface := encoding.Encoder(ferretjson.Default)
	msgpackInterface := encoding.Encoder(ferretmsgpack.Default)
	type mode struct {
		encode func(context.Context, runtime.Value) ([]byte, error)
		name   string
	}

	codecs := []struct {
		name  string
		modes []mode
	}{
		{
			name: "json",
			modes: []mode{
				{name: "concrete_value", encode: func(ctx context.Context, value runtime.Value) ([]byte, error) {
					return encoding.EncodeBytes(ctx, ferretjson.Default, value)
				}},
				{name: "preboxed_value", encode: func(ctx context.Context, value runtime.Value) ([]byte, error) {
					return encoding.EncodeBytes(ctx, jsonInterface, value)
				}},
				{name: "pointer", encode: func(ctx context.Context, value runtime.Value) ([]byte, error) {
					return encoding.EncodeBytes(ctx, &ferretjson.Default, value)
				}},
			},
		},
		{
			name: "msgpack",
			modes: []mode{
				{name: "concrete_value", encode: func(ctx context.Context, value runtime.Value) ([]byte, error) {
					return encoding.EncodeBytes(ctx, ferretmsgpack.Default, value)
				}},
				{name: "preboxed_value", encode: func(ctx context.Context, value runtime.Value) ([]byte, error) {
					return encoding.EncodeBytes(ctx, msgpackInterface, value)
				}},
				{name: "pointer", encode: func(ctx context.Context, value runtime.Value) ([]byte, error) {
					return encoding.EncodeBytes(ctx, &ferretmsgpack.Default, value)
				}},
			},
		},
	}

	for _, codec := range codecs {
		b.Run(codec.name, func(b *testing.B) {
			for _, workload := range []struct {
				value runtime.Value
				name  string
			}{
				{name: "scalar_integer", value: runtime.Int(42)},
				{name: "boolean_array_1024", value: runtime.NewArrayOf(booleans)},
				{name: "null_array_1024", value: runtime.NewArrayOf(nulls)},
			} {
				b.Run(workload.name, func(b *testing.B) {
					for _, mode := range codec.modes {
						b.Run(mode.name, func(b *testing.B) {
							for _, cancellable := range []bool{false, true} {
								name := "background"
								ctx := context.Background()
								if cancellable {
									name = "cancellable"
									var cancel context.CancelFunc
									ctx, cancel = context.WithCancel(ctx)
									b.Cleanup(cancel)
								}

								b.Run(name, func(b *testing.B) {
									data, err := mode.encode(ctx, workload.value)
									if err != nil {
										b.Fatal(err)
									}

									b.ReportAllocs()
									b.SetBytes(int64(len(data)))
									b.ResetTimer()

									for range b.N {
										data, err := mode.encode(ctx, workload.value)
										if err != nil {
											b.Fatal(err)
										}

										benchmarkCodecEncoded = data
									}
								})
							}
						})
					}
				})
			}
		})
	}
}
