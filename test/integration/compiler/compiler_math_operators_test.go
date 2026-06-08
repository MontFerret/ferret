package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestMathOperatorsRejectKnownNonNumericOperands(t *testing.T) {
	expected := func(operator string) E {
		return E{
			Kind:    parserd.SemanticError,
			Message: "Operator '" + operator + "' requires numeric operands",
			Hint:    "Use Int or Float values with this operator.",
		}
	}

	RunSpecs(t, []spec.Spec{
		Failure(`RETURN [120, 45, 300] * 2`, expected("*"), "array multiplication"),
		Failure(`RETURN "3" * 2`, expected("*"), "string multiplication"),
		Failure(`RETURN TRUE - 1`, expected("-"), "boolean subtraction"),
		Failure(`RETURN {} / 2`, expected("/"), "object division"),
		Failure(`RETURN -"x"`, expected("-"), "string unary negative"),
		Failure(`RETURN +"3"`, expected("+"), "string unary positive"),
	})
}

func TestMathOperatorDiagnosticSpansMultiline(t *testing.T) {
	src := `RETURN [
    120,
    45,
    300
] * 2`
	diag := compileMathOperatorDiagnostic(t, "arithmetic.fql", src)

	assertMathOperatorDiagnostic(t, diag, "*", 2)
	assertMathOperatorSpan(t, diag.Spans[0], strings.Index(src, "*"), strings.Index(src, "*")+1, "", true)
	assertMathOperatorSpan(t, diag.Spans[1], strings.Index(src, "["), strings.Index(src, "]")+1, "left operand is Array", false)

	expected := `SemanticError: Operator '*' requires numeric operands
 --> arithmetic.fql:1:8
  |
1 | RETURN [
  |        ^ left operand is Array
2 |     120,
 --> arithmetic.fql:5:3
  |
4 |     300
5 | ] * 2
  |   ^
Hint: Use Int or Float values with this operator.
`
	if actual := diag.Format(); actual != expected {
		t.Fatalf("unexpected formatted diagnostic:\n%s\nDiff expected:\n%s", actual, expected)
	}
}

func TestMathOperatorDiagnosticIdentifiesInvalidOperands(t *testing.T) {
	t.Run("right operand", func(t *testing.T) {
		src := `RETURN 2 * "x"`
		diag := compileMathOperatorDiagnostic(t, "right_operand.fql", src)

		assertMathOperatorDiagnostic(t, diag, "*", 2)
		assertMathOperatorSpan(t, diag.Spans[0], strings.Index(src, "*"), strings.Index(src, "*")+1, "", true)
		assertMathOperatorSpan(t, diag.Spans[1], strings.Index(src, `"x"`), strings.Index(src, `"x"`)+3, "right operand is String", false)
	})

	t.Run("both operands", func(t *testing.T) {
		src := `RETURN TRUE * {}`
		diag := compileMathOperatorDiagnostic(t, "both_operands.fql", src)

		assertMathOperatorDiagnostic(t, diag, "*", 3)
		assertMathOperatorSpan(t, diag.Spans[0], strings.Index(src, "*"), strings.Index(src, "*")+1, "", true)
		assertMathOperatorSpan(t, diag.Spans[1], strings.Index(src, "TRUE"), strings.Index(src, "TRUE")+4, "left operand is Bool", false)
		assertMathOperatorSpan(t, diag.Spans[2], strings.Index(src, "{}"), strings.Index(src, "{}")+2, "right operand is Object", false)
	})

	t.Run("unary operand", func(t *testing.T) {
		src := `RETURN -"x"`
		diag := compileMathOperatorDiagnostic(t, "unary_operand.fql", src)

		assertMathOperatorDiagnostic(t, diag, "-", 2)
		assertMathOperatorSpan(t, diag.Spans[0], strings.Index(src, "-"), strings.Index(src, "-")+1, "", true)
		assertMathOperatorSpan(t, diag.Spans[1], strings.Index(src, `"x"`), strings.Index(src, `"x"`)+3, "operand is String", false)
	})

	t.Run("augmented assignment right operand", func(t *testing.T) {
		src := "VAR total = 2\ntotal *= \"x\"\nRETURN total"
		diag := compileMathOperatorDiagnostic(t, "assignment_operand.fql", src)

		assertMathOperatorDiagnostic(t, diag, "*=", 2)
		assertMathOperatorSpan(t, diag.Spans[0], strings.Index(src, "*="), strings.Index(src, "*=")+2, "", true)
		assertMathOperatorSpan(t, diag.Spans[1], strings.Index(src, `"x"`), strings.Index(src, `"x"`)+3, "right operand is String", false)
	})
}

func compileMathOperatorDiagnostic(t *testing.T, name, src string) *diagnostics.Diagnostic {
	t.Helper()

	_, err := compiler.New().Compile(source.New(name, src))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diag := firstCompilationError(err)
	if diag == nil {
		t.Fatalf("expected diagnostic, got %T", err)
	}

	return diag
}

func assertMathOperatorDiagnostic(t *testing.T, diag *diagnostics.Diagnostic, operator string, spanCount int) {
	t.Helper()

	if diag.Kind != parserd.SemanticError {
		t.Fatalf("unexpected diagnostic kind: got %s want %s", diag.Kind, parserd.SemanticError)
	}
	if diag.Message != "Operator '"+operator+"' requires numeric operands" {
		t.Fatalf("unexpected diagnostic message: %q", diag.Message)
	}
	if diag.Hint != "Use Int or Float values with this operator." {
		t.Fatalf("unexpected diagnostic hint: %q", diag.Hint)
	}
	if len(diag.Spans) != spanCount {
		t.Fatalf("unexpected span count: got %d want %d", len(diag.Spans), spanCount)
	}
}

func assertMathOperatorSpan(t *testing.T, actual diagnostics.ErrorSpan, start, end int, label string, main bool) {
	t.Helper()

	if actual.Span.Start != start || actual.Span.End != end {
		t.Fatalf("unexpected span: got [%d,%d) want [%d,%d)", actual.Span.Start, actual.Span.End, start, end)
	}
	if actual.Label != label {
		t.Fatalf("unexpected span label: got %q want %q", actual.Label, label)
	}
	if actual.Main != main {
		t.Fatalf("unexpected main flag: got %t want %t", actual.Main, main)
	}
}
