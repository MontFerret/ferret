package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func BenchmarkRegexNamedCaptures(b *testing.B) {
	ctx := context.Background()
	for _, fn := range []struct {
		run  func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error)
		name string
	}{
		{name: "find", run: strings.RegexFind},
		{name: "find_all", run: strings.RegexFindAll},
	} {
		for _, pattern := range []struct{ name, expression string }{
			{name: "named", expression: `(?P<word>héllo) (?P<emoji>😀)`},
			{name: "unnamed", expression: `(héllo) (😀)`},
		} {
			for _, input := range []struct{ name, text string }{
				{name: "matches", text: "héllo 😀 world héllo 😀 world"},
				{name: "no_match", text: "goodbye world"},
			} {
				b.Run(fn.name+"/"+pattern.name+"/"+input.name, func(b *testing.B) {
					text, expression := runtime.String(input.text), runtime.String(pattern.expression)
					b.ReportAllocs()

					for b.Loop() {
						result, err := fn.run(ctx, text, expression)
						if err != nil {
							b.Fatal(err)
						}

						cleanupBenchmarkResult = result
					}
				})
			}
		}
	}
}
