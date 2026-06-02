package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	diagpkg "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestForwardDeclarationDiagnosticsTopLevel(t *testing.T) {
	diagnostics := compileDiagnostics(t, `
LET total = price + tax
LET price = 19.99
LET tax = price * 0.08

RETURN total
`)

	if got, want := len(diagnostics), 2; got != want {
		t.Fatalf("expected %d diagnostics, got %d", want, got)
	}

	requireForwardDeclarationDiagnostic(t, diagnostics[0], "price", 3)
	requireForwardDeclarationDiagnostic(t, diagnostics[1], "tax", 4)
}

func TestForwardDeclarationDiagnosticsNestedFunctionScope(t *testing.T) {
	diagnostics := compileDiagnostics(t, `
FUNC outer() (
  FUNC inner() => later
  LET later = 1
  RETURN inner()
)
RETURN outer()
`)

	requireForwardDeclarationDiagnostic(t, diagnostics[0], "later", 4)
}

func TestForwardDeclarationDiagnosticsLoopScope(t *testing.T) {
	diagnostics := compileDiagnostics(t, `
RETURN (
  FOR item IN [1]
    LET value = later
    LET later = item
    RETURN value
)
`)

	requireForwardDeclarationDiagnostic(t, diagnostics[0], "later", 5)
}

func TestForwardDeclarationDiagnosticsAssignmentAndDeleteRoots(t *testing.T) {
	assignment := compileDiagnostics(t, `
target = 1
LET target = 2
RETURN target
`)
	requireForwardDeclarationDiagnostic(t, assignment[0], "target", 3)

	deletion := compileDiagnostics(t, `
DELETE target.value
LET target = { value: 1 }
RETURN target
`)
	requireForwardDeclarationDiagnostic(t, deletion[0], "target", 3)
}

func TestForwardDeclarationDiagnosticsIgnoreDescendantScopeDeclarations(t *testing.T) {
	diagnostics := compileDiagnostics(t, `
LET value = later
LET ignored = (
  FOR item IN [1]
    LET later = item
    RETURN later
)
RETURN value
`)

	if got, want := diagnostics[0].Message, "Variable 'later' is not defined"; got != want {
		t.Fatalf("unexpected diagnostic message: got %q, want %q", got, want)
	}
}

func compileDiagnostics(t *testing.T, src string) []*diagpkg.Diagnostic {
	t.Helper()

	_, err := compiler.New(compiler.WithOptimizationLevel(compiler.O0)).Compile(source.NewAnonymous(src))
	if err == nil {
		t.Fatal("expected compile diagnostics")
	}

	switch e := err.(type) {
	case *diagpkg.Diagnostic:
		return []*diagpkg.Diagnostic{e}
	case *diagpkg.DiagnosticSet:
		out := make([]*diagpkg.Diagnostic, 0, e.Size())
		for idx := 0; idx < e.Size(); idx++ {
			out = append(out, e.Get(idx))
		}
		return out
	default:
		t.Fatalf("unexpected error type: %T", err)
	}

	return nil
}

func requireForwardDeclarationDiagnostic(t *testing.T, actual *diagpkg.Diagnostic, name string, declarationLine int) {
	t.Helper()

	if actual == nil {
		t.Fatal("expected diagnostic")
	}

	if actual.Kind != parserd.NameError {
		t.Fatalf("unexpected diagnostic kind: got %s, want %s", actual.Kind, parserd.NameError)
	}

	if got, want := actual.Message, fmt.Sprintf("Variable '%s' is used before declaration", name); got != want {
		t.Fatalf("unexpected diagnostic message: got %q, want %q", got, want)
	}

	if got, want := actual.Hint, "Move the declaration before this use."; got != want {
		t.Fatalf("unexpected diagnostic hint: got %q, want %q", got, want)
	}

	if got, want := len(actual.Spans), 2; got != want {
		t.Fatalf("expected %d spans, got %d", want, got)
	}

	if got, want := actual.Spans[0].Label, "used before declaration"; got != want {
		t.Fatalf("unexpected main span label: got %q, want %q", got, want)
	}

	if got, want := actual.Spans[1].Label, "declared later"; got != want {
		t.Fatalf("unexpected secondary span label: got %q, want %q", got, want)
	}

	line, _ := actual.Source.LocationAt(actual.Spans[1].Span)
	if line != declarationLine {
		t.Fatalf("unexpected declaration span line: got %d, want %d", line, declarationLine)
	}
}
