package engine

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
)

var benchmarkOutputLength int

func BenchmarkBufferedOutput(b *testing.B) {
	for _, size := range []int{32, 64 << 10, 1 << 20} {
		for _, mode := range []string{"Collect", "Discard", "Writer"} {
			b.Run(fmt.Sprintf("%d/%s", size, mode), func(b *testing.B) {
				data := make([]byte, size)
				var writer bytes.Buffer
				b.ReportAllocs()
				b.SetBytes(int64(size))
				for b.Loop() {
					// This synthetic handle owns the supplied bytes; this benchmark never
					// mutates or retains collected bytes across iterations.
					out := newOutput(b.Context(), runOutcome{content: &encoding.Content{Metadata: encoding.Metadata{Length: int64(size), LengthKnown: true}, Data: data}})
					if mode == "Collect" {
						content, err := out.Collect(b.Context())
						if err != nil {
							b.Fatal(err)
						}

						benchmarkOutputLength = len(content.Data)
					} else {
						writer.Reset()
						if err := out.Consume(b.Context(), func(_ context.Context, chunk []byte) error {
							benchmarkOutputLength = len(chunk)
							if mode == "Writer" {
								_, err := writer.Write(chunk)

								return err
							}

							return nil
						}); err != nil {
							b.Fatal(err)
						}
					}

					if err := out.Close(); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
