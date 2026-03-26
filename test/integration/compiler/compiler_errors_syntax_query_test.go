package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestSyntaxErrorsQueryExpression(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`RETURN QUERY IN doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query literal",
				Hint:    "Provide a query literal, e.g. QUERY `.items` IN doc USING css or QUERY @q IN doc USING css.",
			},
			"Missing query literal after QUERY",
		),
		Failure(
			`RETURN QUERY EXISTS IN doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query literal",
				Hint:    "Provide a query literal, e.g. QUERY `.items` IN doc USING css or QUERY @q IN doc USING css.",
			},
			"Missing query literal after QUERY modifier",
		),
		Failure(
			`RETURN QUERY COUNT IN doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query literal",
				Hint:    "Provide a query literal, e.g. QUERY `.items` IN doc USING css or QUERY @q IN doc USING css.",
			},
			"Missing query literal after QUERY COUNT modifier",
		),
		Failure(
			`RETURN QUERY ONE IN doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query literal",
				Hint:    "Provide a query literal, e.g. QUERY `.items` IN doc USING css or QUERY @q IN doc USING css.",
			},
			"Missing query literal after QUERY ONE modifier",
		),
		Failure(
			"RETURN QUERY `.x` doc USING css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected IN after query literal",
				Hint:    "Add IN <expr>, e.g. QUERY `.items` IN doc USING css.",
			},
			"Missing IN after query literal",
		),
		Failure(
			"RETURN QUERY `.x` IN USING css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after IN",
				Hint:    "Provide a source expression, e.g. QUERY `.items` IN doc USING css.",
			},
			"Missing source expression after IN",
		),
		Failure(
			"RETURN QUERY `.x` IN doc css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected USING <dialect> after IN expression",
				Hint:    "Add USING <dialect>, e.g. QUERY `.items` IN doc USING css.",
			},
			"Missing USING after IN expression",
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
				Message: "Expected options expression after WITH",
				Hint:    "Provide an options expression, e.g. WITH { limit: 10 }.",
			},
			"Missing WITH value",
		).Skip(),
	})
}
