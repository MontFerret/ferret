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

func TestSyntaxErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`
			LET
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name"),

		Failure(
			`
			LET
			RETURN 5
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name 2"),

		Failure(
			`
			LET = 1
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name 3"),

		Failure(
			`
			LET i = NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected a RETURN or FOR clause at end of query",
				Hint:    "All queries must return a value. Add a RETURN statement to complete the query.",
			}, "Missing return statement"),

		Failure(
			`
			LET i = NONE
			RETURN
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'RETURN'",
				Hint:    "Did you forget to provide a value to return?",
			}, "Missing return value"),
		Failure(
			"RETURN DISTINCT",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'RETURN DISTINCT'",
				Hint:    "RETURN DISTINCT treats DISTINCT as a modifier. To return an identifier named DISTINCT, wrap it in parentheses, e.g. RETURN (DISTINCT).",
			},
			"Missing RETURN DISTINCT expression",
		),
		Failure(
			`
			FUNC unique() (
				RETURN DISTINCT
			)
			RETURN unique()
			`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'RETURN DISTINCT'",
				Hint:    "RETURN DISTINCT treats DISTINCT as a modifier. To return an identifier named DISTINCT, wrap it in parentheses, e.g. RETURN (DISTINCT).",
			},
			"Missing RETURN DISTINCT expression in UDF block",
		),
		Failure(
			`
			FUNC f(x)
			  RETURN x
			RETURN f(1)
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected '=>' or '(' after function declaration",
				Hint:    "Use 'FUNC f(x) => expr' or 'FUNC f(x) ( ... RETURN expr )'.",
			}, "Undelimited function body"),
		Failure(
			`
				FUNC f() => RETURN 1
				RETURN f()
			`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=>'",
				Hint:    "Provide an expression, e.g. FUNC f() => x + 1",
			}, "Missing arrow expression"),
		Failure(
			`=>`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=>'",
				Hint:    "Provide an expression, e.g. FUNC f() => x + 1",
			}, "Missing arrow expression at start of input"),
		Failure(
			`
			FUNC run() (
			  VAR i = 0
			  WHILE i < 10
			    i = i + 1
			  RETURN i
			)
			RETURN run()
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Standalone WHILE loops are not supported",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "Standalone WHILE loop inside function block"),

		Failure(
			`
				LET a = 1
				LET b = 2
			LET i = (a ||
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected right-hand expression after '||'",
				Hint:    "Provide an expression after the logical operator, e.g. (a || b).",
			}, "Incomplete logical expression"),

		Failure(
			`
			LET a = 1
			LET b = 2
			LET i = (a OR
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected right-hand expression after 'OR'",
				Hint:    "Provide an expression after the logical operator, e.g. (a OR b).",
			}, "Incomplete logical expression 2"),

		Failure(
			`
			LET a = 1
			LET b = 2
			LET i = (a &&
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected right-hand expression after '&&'",
				Hint:    "Provide an expression after the logical operator, e.g. (a && b).",
			}, "Incomplete logical expression 3"),

		Failure(
			`
			LET a = 1
			LET b = 2
			LET i = (a AND
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected right-hand expression after 'AND'",
				Hint:    "Provide an expression after the logical operator, e.g. (a AND b).",
			}, "Incomplete logical expression 4"),

		Failure(
			`
			LET a = 1
			LET b = 2
			LET i = b > 1 ? a :
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after ':' in ternary operator",
				Hint:    "Provide an expression after the colon to complete the ternary operation.",
			}, "Incomplete ternary expression"),

		Failure(
			`
			LET a = 1
			LET b = 2
			LET i = b > 1 ? 1 + 1 + 1 :
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after ':' in ternary operator",
				Hint:    "Provide an expression after the colon to complete the ternary operation.",
			}, "Incomplete ternary expression 2"),

		Failure(
			`
			LET a = 1
			LET b = 2
			LET i = b > 1 ?
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '?' in ternary operator",
				Hint:    "Provide an expression after the question mark to complete the ternary operation.",
			}, "Incomplete ternary expression 3"),

		Failure(
			`
			LET i = NONE
			RETURN i,
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "--",
				Hint:    "--",
			}, "Dangling comma in return").Skip(),

		Failure(
			`
			LET a = 1
			LET b = 2
			LET i = (a AND b
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed parenthesized expression",
				Hint:    "Add a closing ')' to complete the expression.",
			}, "Unclosed grouping 2"),

		Failure(
			`
			LET i =
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value"),

		Failure(
			`
			LET i =
			LET j = 5
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value 2"),

		Failure(
			`
			LET i =
			FOR j IN [1, 2, 3] RETURN j
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value 3"),

		Failure(
			`
			FN(1,
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after ','",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete function call"),

		Failure(
			`
			FN(,)
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected a valid list of arguments",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete function call 2"),

		Failure(
			`
			FN(
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed function call",
				Hint:    "Add a closing ')' to complete the function call.",
			}, "Incomplete function call 3"),

		Failure(
			`
			FN(1
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed function call",
				Hint:    "Add a closing ')' to complete the function call.",
			}, "Incomplete function call 4"),

		Failure(
			`
			LET r = 0..
			RETURN r
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected end value after '..' in range expression",
				Hint:    "Provide an end value to complete the range, e.g. ..10.",
			}, "Incomplete range"),

		Failure(
			`
				LET r = ..0
				RETURN r
			`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected end value after '..' in range expression",
				Hint:    "Provide an end value to complete the range, e.g. ..10.",
			}, "Incomplete range 2"),
	})
}

func TestMixedFunctionBodySyntaxDiagnosticDoesNotCascade(t *testing.T) {
	query := `
FUNC fib(n) => (
    RETURN MATCH n (
        0 => 0,
        1 => 1,
        _ => fib(n - 1) + fib(n - 2)
    )
)

RETURN fib(10)`

	_, err := compiler.New().Compile(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diag := firstCompilationError(err)
	if diag == nil {
		t.Fatalf("expected diagnostic, got %T", err)
	}

	if diag.Kind != parserd.SyntaxError {
		t.Fatalf("unexpected diagnostic kind: %s", diag.Kind)
	}

	if diag.Message != "Cannot combine arrow and block function body syntax" {
		t.Fatalf("unexpected diagnostic message: %q", diag.Message)
	}

	if diag.Hint != "Use either 'FUNC f(x) => expr' or 'FUNC f(x) ( ... RETURN expr )'." {
		t.Fatalf("unexpected diagnostic hint: %q", diag.Hint)
	}

	if diag.Note != "Remove '=>' to use a block body, or remove RETURN and keep a single expression after '=>'." {
		t.Fatalf("unexpected diagnostic note: %q", diag.Note)
	}

	if len(diag.Spans) == 0 {
		t.Fatal("expected diagnostic span")
	}

	if diag.Spans[0].Label != "RETURN is only valid in a block function body" {
		t.Fatalf("unexpected span label: %q", diag.Spans[0].Label)
	}

	formatted := pkgdiagnostics.Format(err)
	for _, unexpected := range []string{
		"Unclosed parenthesized expression",
		"mismatched input ')' expecting <EOF>",
		"Variable 'n' is not defined",
	} {
		if strings.Contains(formatted, unexpected) {
			t.Fatalf("formatted diagnostic contains cascade %q:\n%s", unexpected, formatted)
		}
	}
}

func TestMissingFunctionParamsCloseDiagnosticDoesNotCascade(t *testing.T) {
	query := `FUNC fib (
    RETURN MATCH n (
        0 => 0,
        1 => 1,
        _ => fib(n - 1) + fib(n - 2)
    )
)

RETURN fib(10)`

	_, err := compiler.New().Compile(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diag := firstCompilationError(err)
	if diag == nil {
		t.Fatalf("expected diagnostic, got %T", err)
	}

	if diag.Kind != parserd.SyntaxError {
		t.Fatalf("unexpected diagnostic kind: %s", diag.Kind)
	}

	if diag.Message != "Expected function parameters before function body" {
		t.Fatalf("unexpected diagnostic message: %q", diag.Message)
	}

	if diag.Hint != "Add a parameter list before the block body, e.g. FUNC fib(n) ( ... RETURN expr ). Use FUNC fib() ( ... ) for no parameters." {
		t.Fatalf("unexpected diagnostic hint: %q", diag.Hint)
	}

	if len(diag.Spans) == 0 {
		t.Fatal("expected diagnostic span")
	}

	if diag.Spans[0].Label != "missing parameter list before function body" {
		t.Fatalf("unexpected span label: %q", diag.Spans[0].Label)
	}

	line, col := diag.Source.LocationAt(diag.Spans[0].Span)
	if line != 1 || col != 10 {
		t.Fatalf("unexpected span location: got %d:%d, want 1:10", line, col)
	}

	formatted := pkgdiagnostics.Format(err)
	if got := strings.Count(formatted, "SyntaxError:"); got != 1 {
		t.Fatalf("expected one syntax diagnostic, got %d:\n%s", got, formatted)
	}

	if !strings.Contains(formatted, "1 | FUNC fib (\n  |          ^ missing parameter list before function body\n2 |     RETURN MATCH n (") {
		t.Fatalf("diagnostic should point at the premature function body paren, got:\n%s", formatted)
	}

	for _, unexpected := range []string{
		"mismatched input 'RETURN'",
		"mismatched input ')' expecting <EOF>",
	} {
		if strings.Contains(formatted, unexpected) {
			t.Fatalf("formatted diagnostic contains cascade %q:\n%s", unexpected, formatted)
		}
	}
}
