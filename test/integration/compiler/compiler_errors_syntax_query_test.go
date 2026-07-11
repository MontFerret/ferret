package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestSyntaxErrorsQueryExpression(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`RETURN QUERY IN doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query expression",
				Hint:    "Provide a query expression, e.g. QUERY `.items` IN doc, QUERY email.body IN doc, or QUERY @q IN doc USING css.",
			},
			"Missing query expression after QUERY",
		),
		Failure(
			"RETURN QUERY ANY `.items` IN doc USING css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected IN after query expression",
				Hint:    "Add IN <expr>, e.g. QUERY `.items` IN doc or QUERY email.body IN doc.",
			},
			"Invalid ANY query payload",
		),
		Failure(
			"RETURN QUERY VALUE `.items` IN doc USING css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected IN after query expression",
				Hint:    "Add IN <expr>, e.g. QUERY `.items` IN doc or QUERY email.body IN doc.",
			},
			"Invalid VALUE query payload",
		),
		Failure(
			`RETURN QUERY EXISTS IN doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query expression",
				Hint:    "Provide a query expression, e.g. QUERY `.items` IN doc, QUERY email.body IN doc, or QUERY @q IN doc USING css.",
			},
			"Missing query expression after QUERY modifier",
		),
		Failure(
			`RETURN QUERY COUNT IN doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query expression",
				Hint:    "Provide a query expression, e.g. QUERY `.items` IN doc, QUERY email.body IN doc, or QUERY @q IN doc USING css.",
			},
			"Missing query expression after QUERY COUNT modifier",
		),
		Failure(
			`RETURN QUERY ONE IN doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query expression",
				Hint:    "Provide a query expression, e.g. QUERY `.items` IN doc, QUERY email.body IN doc, or QUERY @q IN doc USING css.",
			},
			"Missing query expression after QUERY ONE modifier",
		),
		Failure(
			"RETURN QUERY `.x` doc USING css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected IN after query expression",
				Hint:    "Add IN <expr>, e.g. QUERY `.items` IN doc or QUERY email.body IN doc.",
			},
			"Missing IN after query expression",
		),
		Failure(
			"RETURN QUERY email.body model USING summarize",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected IN after query expression",
				Hint:    "Add IN <expr>, e.g. QUERY `.items` IN doc or QUERY email.body IN doc.",
			},
			"Missing IN after member query expression",
		),
		Failure(
			"RETURN QUERY `.x` IN USING css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after IN",
				Hint:    "Provide a source expression, e.g. QUERY `.items` IN doc.",
			},
			"Missing source expression after IN",
		),
		Failure(
			"RETURN QUERY `.x` IN doc css",
			E{
				Kind: parserd.SyntaxError,
			},
			"Unexpected trailing token after query source",
		),
		Failure(
			"RETURN QUERY `.x` IN doc USING",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected dialect identifier after USING",
				Hint:    "Provide a dialect identifier, e.g. USING css.",
			},
			"Missing dialect after USING",
		),
		Failure(
			"RETURN QUERY `.x` IN doc USING \"css\"",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Dialect after USING must be an identifier",
				Hint:    "Provide a dialect identifier such as css or xpath.",
			},
			"Invalid dialect token after USING",
		),
		Failure(
			"RETURN QUERY `.x` IN @doc USING css WITH )",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query params expression after WITH",
				Hint:    "Provide a params expression, e.g. WITH { params: [1] }.",
			},
			"Missing WITH value",
		),
		Failure(
			"RETURN QUERY `.x` IN @doc USING css OPTIONS )",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query options expression after OPTIONS",
				Hint:    "Provide an options expression, e.g. OPTIONS { timeout: 5000 }.",
			},
			"Missing OPTIONS value",
		),
		Failure(
			"RETURN QUERY `.x` IN @doc USING css OPTIONS { timeout: 5000 } WITH { params: [1] }",
			E{
				Kind:    parserd.SyntaxError,
				Message: "WITH must appear before OPTIONS in QUERY",
				Hint:    "Move the WITH clause before OPTIONS.",
			},
			"WITH after OPTIONS",
		),
	})
}

func TestSyntaxErrorsQueryExpressionWithAfterOptionsSpan(t *testing.T) {
	query := "RETURN QUERY `.x` IN @doc USING css OPTIONS { timeout: 5000 } WITH { params: [1] }"

	_, err := compiler.New().Compile(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diag := firstCompilationError(err)
	if diag == nil || len(diag.Spans) == 0 {
		t.Fatalf("expected diagnostic span, got %v", err)
	}

	span := diag.Spans[0].Span
	if got := query[span.Start:span.End]; got != "WITH" {
		t.Fatalf("expected diagnostic to point at WITH, got %q", got)
	}
}
