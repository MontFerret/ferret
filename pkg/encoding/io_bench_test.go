package encoding_test

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkCodecWriter(b *testing.B) {
	values := make([]runtime.Value, 16_384)
	props := make(map[string]runtime.Value, 4096)
	for i := range values {
		values[i] = runtime.NewInt(i)
	}
	for i := 0; i < 4096; i++ {
		props[strconv.Itoa(i)] = runtime.NewInt(i)
	}
	// The payloads are constructed before timing; Discard never accumulates the
	// encoded document, while Buffer measures a reusable caller-owned destination.
	workloads := []struct {
		value runtime.Value
		name  string
	}{
		{name: "array", value: runtime.NewArrayOf(values)},
		{name: "object", value: runtime.NewObjectWith(props)},
		{name: "range", value: runtime.NewRange(0, 65_535)},
		{name: "string_1MiB", value: runtime.NewString(strings.Repeat("a", 1<<20))},
		{name: "binary_1MiB", value: runtime.NewBinary(bytes.Repeat([]byte{255}, 1<<20))},
	}
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		b.Run(codec.ContentType(), func(b *testing.B) {
			for _, workload := range workloads {
				b.Run(workload.name, func(b *testing.B) {
					for _, configured := range []bool{false, true} {
						name := "default"
						enc := encoding.Encoder(codec)
						if configured {
							name = "configured"
							enc = codec.EncodeWith().PreHook(func(context.Context, runtime.Value) error { return nil }).Encoder()
						}
						b.Run(name, func(b *testing.B) {
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
									for _, target := range []string{"Buffer", "Discard"} {
										b.Run(target, func(b *testing.B) {
											var buffer bytes.Buffer
											var dst io.Writer = io.Discard
											if target == "Buffer" {
												dst = &buffer
											}
											if err := enc.Encode(ctx, dst, workload.value); err != nil {
												b.Fatal(err)
											}
											b.ReportAllocs()
											b.ResetTimer()
											for i := 0; i < b.N; i++ {
												buffer.Reset()
												if err := enc.Encode(ctx, dst, workload.value); err != nil {
													b.Fatal(err)
												}
											}
										})
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
