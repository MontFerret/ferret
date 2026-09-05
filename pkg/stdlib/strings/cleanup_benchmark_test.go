package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

var cleanupBenchmarkResult runtime.Value

func BenchmarkRepeat(b *testing.B) {
	ctx := context.Background()
	for _, input := range []struct{ name, text string }{
		{name: "ASCII", text: "a"},
		{name: "Unicode", text: "😀"},
	} {
		for _, size := range []struct {
			name  string
			bytes int
		}{
			{name: "small", bytes: 128},
			{name: "1MiB", bytes: 1 << 20},
		} {
			b.Run(input.name+"/"+size.name, func(b *testing.B) {
				text := runtime.String(input.text)
				count := runtime.Int(size.bytes / len(input.text))
				b.SetBytes(int64(size.bytes))
				b.ReportAllocs()

				for b.Loop() {
					result, err := strings.Repeat(ctx, text, count)
					if err != nil {
						b.Fatal(err)
					}

					cleanupBenchmarkResult = result
				}
			})
		}
	}
}

func BenchmarkStringsCleanup(b *testing.B) {
	ctx := context.Background()
	text := runtime.NewString("héllo 😀 world héllo 😀 world")
	for _, test := range []struct {
		run  func() (runtime.Value, error)
		name string
	}{
		{name: "left", run: func() (runtime.Value, error) { return strings.Left(ctx, text, runtime.Int(5)) }},
		{name: "right", run: func() (runtime.Value, error) { return strings.Right(ctx, text, runtime.Int(5)) }},
		{name: "substring", run: func() (runtime.Value, error) { return strings.Substring(ctx, text, runtime.Int(2), runtime.Int(8)) }},
		{name: "fmt", run: func() (runtime.Value, error) {
			return strings.Fmt(ctx, runtime.String("{1}: {0} {1}"), text, runtime.Int(42))
		}},
		{name: "regex_find", run: func() (runtime.Value, error) { return strings.RegexFind(ctx, text, runtime.String(`(héllo) (😀)`)) }},
		{name: "regex_replace", run: func() (runtime.Value, error) {
			return strings.RegexReplace(ctx, text, runtime.String(`héllo`), runtime.String("hello"))
		}},
		{name: "split", run: func() (runtime.Value, error) { return strings.Split(ctx, text, runtime.String(" ")) }},
		{name: "join", run: func() (runtime.Value, error) {
			return strings.Join(ctx, runtime.NewArrayWith(text, text, text), runtime.String(","))
		}},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				result, err := test.run()
				if err != nil {
					b.Fatal(err)
				}

				cleanupBenchmarkResult = result
			}
		})
	}
}
