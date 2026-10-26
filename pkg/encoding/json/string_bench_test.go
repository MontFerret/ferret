package json_test

import (
	"strings"
	"testing"
	"time"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkJSONRootStringFallbacks(b *testing.B) {
	codec := json.Default

	for _, tc := range []struct {
		value runtime.Value
		name  string
	}{
		{name: "unicode_duration", value: runtime.NewDuration(1500 * time.Nanosecond)},
		{name: "ascii_duration", value: runtime.NewDuration(1500 * time.Millisecond)},
		{name: "short_escaped", value: runtime.String("a\n\"\\b")},
		{name: "short_unicode", value: runtime.String("aλµb")},
		{name: "ascii_over_limit", value: runtime.String(strings.Repeat("a", 4090))},
	} {
		b.Run(tc.name, func(b *testing.B) {
			out, err := ferretencoding.EncodeBytes(b.Context(), codec, tc.value)
			if err != nil {
				b.Fatalf("setup encode failed: %v", err)
			}

			b.ReportAllocs()
			b.SetBytes(int64(len(out)))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				out, err := ferretencoding.EncodeBytes(b.Context(), codec, tc.value)
				if err != nil {
					b.Fatalf("encode failed: %v", err)
				}

				benchmarkJSONEncoded = out
			}
		})
	}
}
