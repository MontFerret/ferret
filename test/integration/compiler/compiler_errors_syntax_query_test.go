package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestSyntaxErrorsQueryExpression(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`RETURN QUERY FROM doc USING css`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "QUERY requires a query literal",
				Hint:    "Provide a query literal, e.g. QUERY `.items` FROM doc USING css.",
			},
			"Missing query literal after QUERY",
		),
		ErrorCase(
			"RETURN QUERY `.x` doc USING css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FROM after query literal",
				Hint:    "Add FROM <expr>, e.g. QUERY `.items` FROM doc USING css.",
			},
			"Missing FROM after query literal",
		),
		ErrorCase(
			"RETURN QUERY `.x` FROM USING css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after FROM",
				Hint:    "Provide a source expression, e.g. QUERY `.items` FROM doc USING css.",
			},
			"Missing source expression after FROM",
		),
		ErrorCase(
			"RETURN QUERY `.x` FROM doc css",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected USING <dialect> after FROM expression",
				Hint:    "Add USING <dialect>, e.g. QUERY `.items` FROM doc USING css.",
			},
			"Missing USING after FROM expression",
		),
		ErrorCase(
			"RETURN QUERY `.x` FROM doc USING",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected dialect identifier after USING",
				Hint:    "Provide a dialect identifier, e.g. USING css.",
			},
			"Missing dialect after USING",
		),
		ErrorCase(
			"RETURN QUERY `.x` FROM doc USING \"css\"",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Dialect after USING must be an identifier",
				Hint:    "Provide a dialect identifier such as css or xpath.",
			},
			"Invalid dialect token after USING",
		),
		ErrorCase(
			"RETURN QUERY `.x` FROM @doc USING css WITH RETURN",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected options expression after WITH",
				Hint:    "Provide an options expression, e.g. WITH { limit: 10 }.",
			},
			"Missing WITH value",
		),
	})
}
