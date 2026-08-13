package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestReturnlessForValueDiagnostics(t *testing.T) {
	const (
		message = "A FOR loop used as an expression must return a value."
		hint    = "Add RETURN to the loop body, or use the loop as a statement."
		label   = "returnless FOR expression"
	)

	tests := []struct {
		name   string
		query  string
		marker string
	}{
		{name: "initializer", query: `LET values = (FOR value IN [1] {})`, marker: "FOR value"},
		{
			name:   "initializer with recovery",
			query:  `LET values = (FOR value IN [1] {}) ON ERROR RETURN []`,
			marker: "FOR value",
		},
		{name: "return", query: `RETURN FOR value IN [1] {}`, marker: "FOR value"},
		{name: "grouped return", query: `RETURN (FOR value IN [1] {})`, marker: "FOR value"},
		{name: "return distinct", query: `RETURN DISTINCT FOR value IN [1] {}`, marker: "FOR value"},
		{name: "argument", query: `RETURN LENGTH((FOR value IN [1] {}))`, marker: "FOR value"},
		{name: "array member", query: `RETURN [(FOR value IN [1] {})]`, marker: "FOR value"},
		{name: "object member", query: `RETURN { values: (FOR value IN [1] {}) }`, marker: "FOR value"},
		{name: "member source", query: `RETURN (FOR value IN [1] {}).length`, marker: "FOR value"},
		{name: "operator operand", query: `RETURN (FOR value IN [1] {}) == []`, marker: "FOR value"},
		{
			name:   "nested pass-through operand",
			query:  `RETURN FOR outer IN [1] { FOR inner IN [outer] {} }`,
			marker: "FOR inner",
		},
		{
			name: "nested pass-through operand after outer statements",
			query: `VAR n = 0
RETURN FOR outer IN [1] {
  n += 1
  FOR inner IN [outer] {
    n += inner
  }
}`,
			marker: "FOR inner",
		},
		{name: "while return", query: `RETURN FOR WHILE false {}`, marker: "FOR WHILE"},
		{name: "do while return", query: `RETURN FOR DO WHILE false {}`, marker: "FOR DO"},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("O%d/%s", level, test.name), func(t *testing.T) {
				_, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.NewAnonymous(test.query))
				if err == nil {
					t.Fatal("expected returnless FOR diagnostic")
				}

				diagnostic := onlyDiagnostic(t, err)
				if diagnostic.Kind != parserd.SemanticError {
					t.Fatalf("diagnostic kind = %s, want %s", diagnostic.Kind, parserd.SemanticError)
				}

				if diagnostic.Message != message {
					t.Fatalf("diagnostic message = %q, want %q", diagnostic.Message, message)
				}

				if diagnostic.Hint != hint {
					t.Fatalf("diagnostic hint = %q, want %q", diagnostic.Hint, hint)
				}

				if got, want := len(diagnostic.Spans), 1; got != want {
					t.Fatalf("diagnostic span count = %d, want %d", got, want)
				}

				span := diagnostic.Spans[0]
				if !span.Main {
					t.Fatal("diagnostic span is not primary")
				}

				if span.Label != label {
					t.Fatalf("diagnostic span label = %q, want %q", span.Label, label)
				}

				start := strings.Index(test.query, test.marker)
				if start < 0 {
					t.Fatalf("marker %q not found in query", test.marker)
				}

				closeOffset := strings.Index(test.query[start:], "}")
				if closeOffset < 0 {
					t.Fatal("loop closing brace not found")
				}

				wantSpan := source.Span{Start: start, End: start + closeOffset + 1}
				if span.Span != wantSpan {
					t.Fatalf("diagnostic span = %#v, want %#v", span.Span, wantSpan)
				}
			})
		}
	}
}

func onlyDiagnostic(t *testing.T, err error) *diagnostics.Diagnostic {
	t.Helper()

	switch actual := err.(type) {
	case *diagnostics.Diagnostic:
		return actual
	case *diagnostics.DiagnosticSet:
		if got, want := actual.Size(), 1; got != want {
			t.Fatalf("diagnostic count = %d, want %d:\n%s", got, want, actual.Format())
		}

		return actual.First()
	default:
		t.Fatalf("compile error type = %T, want diagnostic", err)

		return nil
	}
}
