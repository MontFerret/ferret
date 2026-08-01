package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestDurationLiteralCompilationErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(`RETURN 1fortnight`, E{Kind: parserd.SyntaxError}, "unknown duration unit"),
		Failure(`RETURN 1h30m`, E{Kind: parserd.SyntaxError}, "compound duration literal"),
		Failure(`RETURN 1 ms`, E{Kind: parserd.SyntaxError}, "whitespace before duration unit"),
		Failure(`RETURN 1e999s`, E{
			Kind:    parserd.SyntaxError,
			Message: "Duration literal is out of range",
			Hint:    "Use a duration value that fits within the signed nanosecond range.",
		}, "source-aware duration overflow"),
	})
}

func TestDurationLiteralOverflowDiagnosticSpan(t *testing.T) {
	query := "RETURN 1e999s"

	_, err := compiler.New().Compile(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diagnostic := firstCompilationError(err)
	if diagnostic == nil || len(diagnostic.Spans) == 0 {
		t.Fatalf("expected duration diagnostic span, got %v", err)
	}

	span := diagnostic.Spans[0].Span
	if got := query[span.Start:span.End]; got != "1e999s" {
		t.Fatalf("duration diagnostic points at %q, want 1e999s", got)
	}
}
