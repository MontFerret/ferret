package compiler

import (
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestCompilePublishesByteSpans(t *testing.T) {
	const query = "LET prefix = \"é\"\nRETURN TEST(prefix)"
	src := source.New("cell:unicode", query)
	for _, level := range []OptimizationLevel{None, Basic, Full} {
		program, err := mustNewCompiler(t, WithOptimizationLevel(level)).Compile(t.Context(), src)
		if err != nil {
			t.Fatal(err)
		}

		start := strings.LastIndex(query, "prefix")
		want := source.Span{Start: start, End: start + len("prefix")}
		found := false
		for _, spans := range program.Metadata.CallArgumentSpans {
			for _, span := range spans {
				if span == want {
					found = true
				}
			}
		}

		if !found {
			t.Fatalf("optimization %v: missing byte argument span %+v in %+v", level, want, program.Metadata.CallArgumentSpans)
		}
	}

	program, err := mustNewCompiler(t, WithDebugInfo()).Compile(t.Context(), src)
	if err != nil {
		t.Fatal(err)
	}

	start := strings.Index(query, "RETURN")
	want := source.Span{Start: start, End: len(query)}
	found := false
	for _, point := range program.Metadata.DebugPoints {
		if point.Span == want {
			found = true
		}
	}

	if !found {
		t.Fatalf("missing byte debug point %+v in %+v", want, program.Metadata.DebugPoints)
	}

	found = false
	for _, span := range program.Metadata.DebugSpans {
		if span == want {
			found = true
		}
	}

	if !found {
		t.Fatalf("missing byte instruction span %+v in %+v", want, program.Metadata.DebugSpans)
	}
}

func TestCompileAndAnalyzeDiagnosticByteSpansAgree(t *testing.T) {
	for _, query := range []string{
		`RETURN "é" + missing`,
		"LET prefix = \"é\"\nRETURN missing",
		`LET prefix = "é" RETURN )`,
	} {
		src := source.NewAnonymous(query)
		c := mustNewCompiler(t)
		_, compileErr := c.Compile(t.Context(), src)
		_, analysisErr := c.Analyze(src)
		var compiled, analyzed *diagnostics.Diagnostic
		if !errors.As(compileErr, &compiled) || !errors.As(analysisErr, &analyzed) {
			t.Fatalf("expected diagnostics: compile=%v analyze=%v", compileErr, analysisErr)
		}

		start := strings.Index(query, "missing")
		if start < 0 {
			start = strings.Index(query, "RETURN") + len("RETURN")
		}

		if len(compiled.Spans) == 0 || compiled.Spans[0].Span.Start != start ||
			compiled.Spans[0].Span != analyzed.Spans[0].Span {
			t.Fatalf("%q: compile=%+v analyze=%+v, want start %d", query, compiled.Spans, analyzed.Spans, start)
		}
	}
}
