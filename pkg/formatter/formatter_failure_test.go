package formatter_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type failingWriter struct {
	failure    error
	panicValue any
	buffer     bytes.Buffer
	limit      int
}

func (w *failingWriter) Write(p []byte) (int, error) {
	remaining := w.limit - w.buffer.Len()
	if len(p) <= remaining {
		return w.buffer.Write(p)
	}

	n, _ := w.buffer.Write(p[:remaining])
	if w.panicValue != nil {
		panic(w.panicValue)
	}

	return n, w.failure
}

func TestFormatterWriterFailures(t *testing.T) {
	f, err := formatter.New()
	if err != nil {
		t.Fatal(err)
	}

	limitErr := errors.New("output limit")
	panicErr := &diagnostics.Diagnostic{
		Kind:    diagnostics.TypeError,
		Message: "external writer diagnostic",
		Spans:   []diagnostics.ErrorSpan{{Span: source.Span{Start: 3, End: 4}}},
	}

	for _, input := range []string{"RETURN 1", "// 😀\nRETURN 1"} {
		t.Run(input, func(t *testing.T) {
			src := source.New("failure.fql", input)

			for _, test := range []struct {
				panicValue any
				name       string
			}{
				{name: "ordinary"},
				{name: "string panic", panicValue: "injected panic"},
				{name: "error panic", panicValue: panicErr},
				{name: "value panic", panicValue: 42},
			} {
				for _, limit := range []int{0, 2, 7} {
					t.Run(fmt.Sprintf("%s/%d", test.name, limit), func(t *testing.T) {
						writer := &failingWriter{limit: limit, failure: limitErr, panicValue: test.panicValue}
						err := f.Format(writer, src)
						if err == nil {
							t.Fatalf("failure returned nil with output %q", writer.buffer.String())
						}

						var diagnostic *diagnostics.Diagnostic
						if test.panicValue == nil {
							if !errors.Is(err, limitErr) || errors.As(err, &diagnostic) {
								t.Fatalf("ordinary writer error changed identity: %v", err)
							}
						} else {
							if !errors.As(err, &diagnostic) || diagnostic.Kind != diagnostics.UnexpectedError || diagnostic.Source.ID() != src.ID() || len(diagnostic.Spans) != 0 {
								t.Fatalf("panic diagnostic = %+v, error = %v", diagnostic, err)
							}

							if test.panicValue == panicErr && (!errors.Is(err, panicErr) || panicErr.Spans[0].Span != (source.Span{Start: 3, End: 4})) {
								t.Fatal("panic cause was replaced or normalized")
							}
						}

						if writer.buffer.Len() != limit {
							t.Fatalf("accepted %d bytes, want %d", writer.buffer.Len(), limit)
						}

						var reused bytes.Buffer
						if err := f.Format(&reused, source.NewAnonymous("RETURN 1")); err != nil || reused.String() != "return 1" {
							t.Fatalf("reuse failed: output %q, error %v", reused.String(), err)
						}
					})
				}
			}
		})
	}
}
