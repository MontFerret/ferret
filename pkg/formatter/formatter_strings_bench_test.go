package formatter_test

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkFormatterStringLiterals(b *testing.B) {
	for _, tt := range []struct {
		name  string
		input string
	}{
		{"Conversion", `RETURN ['alpha', 'beta', 'gamma', 'delta']`},
		{"Fallback", `RETURN ['a"b', 'c"d', 'e"f', 'g"h']`},
		{"Escapes", `RETURN ['a\b', 'c\\d', 'e\nf\tg', 'h\u0041']`},
		{"Unicode", `RETURN [´é中😀´, 'é\中\n😀', ´a\´]`},
		{"Multiline", "RETURN ['a\nb', 'c\\\nd', 'e\n\tf']"},
	} {
		b.Run(tt.name, func(b *testing.B) {
			f, err := formatter.New()
			if err != nil {
				b.Fatal(err)
			}

			src := source.NewAnonymous(tt.input)
			var out bytes.Buffer
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

			b.ReportMetric(float64(out.Len()), "output-B")
		})
	}
}
