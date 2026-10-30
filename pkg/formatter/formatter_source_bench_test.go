package formatter_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkFormatterSourceOffsets(b *testing.B) {
	f, err := formatter.New()
	if err != nil {
		b.Fatal(err)
	}

	for _, text := range []struct{ name, value string }{{"ASCII", "ascii"}, {"Unicode", "é中😀"}} {
		for _, count := range []int{16, 128, 1024} {
			b.Run(fmt.Sprintf("%s/%d", text.name, count), func(b *testing.B) {
				var query strings.Builder
				for i := 0; i < count; i++ {
					fmt.Fprintf(&query, "// %s\nLET v%d = [\"%s\", /* keep */ 1]\n", text.value, i, text.value)
				}

				query.WriteString("RETURN 1 // end\n")
				src := source.NewAnonymous(query.String())
				var out bytes.Buffer

				// Warm parser caches and output capacity outside the measurements.
				if err := f.Format(&out, src); err != nil {
					b.Fatal(err)
				}

				b.ReportAllocs()
				b.SetBytes(int64(src.Length()))
				b.ResetTimer()

				for b.Loop() {
					out.Reset()
					if err := f.Format(&out, src); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
