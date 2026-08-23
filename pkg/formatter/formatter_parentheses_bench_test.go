package formatter_test

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const formatterParenthesesQuery = `
FUNC effect(a, b) {
  (((a + b)))
  (FOR value IN [1, 2, 3] {
    RETURN ((value + a)) * (b + 1)
  })
}

RETURN ((1 + 2)) * (3 + 4)
`

func BenchmarkFormatterParentheses(b *testing.B) {
	format, err := formatter.New()
	if err != nil {
		b.Fatalf("formatter.New() error = %v", err)
	}

	src := source.New("parentheses_benchmark", formatterParenthesesQuery)
	var output bytes.Buffer

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		output.Reset()

		if err := format.Format(&output, src); err != nil {
			b.Fatal(err)
		}
	}
}
