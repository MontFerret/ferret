package diagnostics

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestSpanRendererInsertionPoints(t *testing.T) {
	for _, test := range []struct {
		name  string
		query string
		want  string
		span  source.Span
	}{
		{
			name: "within line", query: "RETURN value", span: source.Span{Start: 6, End: 6},
			want: " --> query.fql:1:7\n  |\n1 | RETURN value\n  |       ^ missing value\n",
		},
		{
			name: "EOF", query: "RETURN", span: source.Span{Start: 6, End: 6},
			want: " --> query.fql:1:7\n  |\n1 | RETURN\n  |       ^ missing value\n",
		},
		{
			name: "Unicode EOF", query: "é RETURN", span: source.Span{Start: 9, End: 9},
			want: " --> query.fql:1:10\n  |\n1 | é RETURN\n  |         ^ missing value\n",
		},
		{
			name: "start of nonempty source", query: "RETURN", span: source.Span{Start: 0, End: 0},
			want: " --> query.fql:1:1\n  |\n1 | RETURN\n  | ^ missing value\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var output strings.Builder
			if !(SpanRenderer{}).Render(&output, source.New("query.fql", test.query), test.span, "missing value") {
				t.Fatal("expected insertion point to render")
			}

			if got := output.String(); got != test.want {
				t.Fatalf("rendered insertion = %q, want %q", got, test.want)
			}
		})
	}
}

func TestSpanRendererRejectsInvalidSpans(t *testing.T) {
	src := source.New("query.fql", "RETURN")
	for _, span := range []source.Span{
		{Start: -1, End: 0},
		{Start: 2, End: 1},
		{Start: 6, End: 7},
		{Start: 7, End: 7},
	} {
		var output strings.Builder
		if (SpanRenderer{}).Render(&output, src, span, "invalid") {
			t.Errorf("rendered invalid span %+v", span)
		}

		if output.Len() != 0 {
			t.Errorf("invalid span %+v produced output %q", span, output.String())
		}
	}
}
