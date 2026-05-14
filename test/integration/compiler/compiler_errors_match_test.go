package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestMatchErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`RETURN MATCH { a: 1, b: 2 } ( { a: v, b: v } => v, _ => 0, )`,
			E{
				Kind:    parserd.NameError,
				Message: "duplicate binding 'v' in MATCH pattern",
			},
		),
		Failure(
			`
			FUNC fib(n) (
				RETURN MATCH n (
					0 => 0
					1 => 1
					_ => fib(n - 1) + fib(n - 2)
				)
			)

			RETURN fib(10)
		`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected ',' between MATCH arms",
				Hint:    "Separate MATCH arms with commas, e.g. 0 => 0, 1 => 1, _ => 0.",
			},
			"Missing comma between MATCH pattern arms",
		),
		Failure(
			`
			RETURN MATCH (
				WHEN TRUE => "yes"
				WHEN FALSE => "no",
				_ => "fallback",
			)
		`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected ',' between MATCH arms",
				Hint:    "Separate MATCH arms with commas, e.g. 0 => 0, 1 => 1, _ => 0.",
			},
			"Missing comma between MATCH guard arms",
		),
		Failure(
			`RETURN MATCH 0 ( 1 => 1 _ => 0 )`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected ',' between MATCH arms",
				Hint:    "Separate MATCH arms with commas, e.g. 0 => 0, 1 => 1, _ => 0.",
			},
			"Missing comma before MATCH default arm",
		),
	})
}

func TestMatchMissingCommaDiagnosticSpan(t *testing.T) {
	c := compiler.New()
	query := `FUNC fib(n) (
    RETURN MATCH n (
        0 => 0
        1 => 1
        _ => fib(n - 1) + fib(n - 2)
    )
)

RETURN fib(10)`

	_, err := c.Compile(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diag := firstCompilationError(err)
	if diag == nil {
		t.Fatalf("expected diagnostic, got %T", err)
	}

	if diag.Message != "Expected ',' between MATCH arms" {
		t.Fatalf("unexpected diagnostic message: %q", diag.Message)
	}

	if len(diag.Spans) == 0 {
		t.Fatal("expected diagnostic span")
	}

	if diag.Spans[0].Label != "missing comma" {
		t.Fatalf("unexpected span label: %q", diag.Spans[0].Label)
	}

	line, col := diag.Source.LocationAt(diag.Spans[0].Span)
	if line != 3 || col != 15 {
		t.Fatalf("unexpected span location: got %d:%d, want 3:15", line, col)
	}

	formatted := pkgdiagnostics.Format(err)
	if !strings.Contains(formatted, "3 |         0 => 0\n  |               ^ missing comma\n4 |         1 => 1") {
		t.Fatalf("diagnostic should point after previous MATCH arm value, got:\n%s", formatted)
	}
}
