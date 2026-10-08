package encoding_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkCodecIterable(b *testing.B) {
	items := make([]runtime.Value, 1024)
	for i := range items {
		items[i] = runtime.NewInt(i)
	}

	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		b.Run(codec.ContentType(), func(b *testing.B) {
			for _, ownership := range []string{"owned", "borrowed"} {
				b.Run(ownership, func(b *testing.B) {
					iterator := &selfCodecIterator{codecIterator: codecIterator{items: items}}
					var value runtime.Value = iterator
					if ownership == "owned" {
						value = &codecIterable{iter: iterator}
					}

					for _, target := range []string{"Buffer", "Discard"} {
						b.Run(target, func(b *testing.B) {
							var buffer bytes.Buffer
							writer := io.Writer(io.Discard)
							if target == "Buffer" {
								writer = &buffer
							}

							ctx := context.Background()
							if err := codec.Encode(ctx, writer, value); err != nil {
								b.Fatal(err)
							}

							b.ReportAllocs()
							b.ResetTimer()

							for i := 0; i < b.N; i++ {
								iterator.index = 0
								buffer.Reset()

								if err := codec.Encode(ctx, writer, value); err != nil {
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
