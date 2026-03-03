package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestSyntaxErrorsArrayOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`RETURN doc[~]`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query literal after '~'",
				Hint:    "Provide a query literal, e.g. doc[~ css`...`].",
			},
			"Missing query literal after '~'",
		),
		ErrorCase(
			`RETURN doc[~ 'x']`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query literal after '~'",
				Hint:    "Provide a query literal, e.g. doc[~ css`...`].",
			},
			"Missing query type before literal",
		),
		ErrorCase(
			`LET doc = {} RETURN doc[~ css()]`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query string after 'css'",
				Hint:    "Provide a query string, e.g. doc[~ css`...`].",
			},
			"Missing query string after type",
		),
		ErrorCase(
			`LET doc = {} RETURN doc[~ css"x"`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query type before query literal",
				Hint:    "Provide a type name before the query string, e.g. doc[~ css`...`].",
			},
			"Missing query type before query literal",
		),
		ErrorCase(
			`RETURN [1, 2][* RETURN]`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'RETURN' in array operator",
				Hint:    "Provide a projection expression, e.g. [* RETURN .].",
			},
			"Missing inline RETURN expression",
		),
		ErrorCase(
			`RETURN [1, 2][? NONE]`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FILTER after quantifier in array filter",
				Hint:    "Add a FILTER expression, e.g. [? NONE FILTER <expr>].",
			},
			"Missing FILTER after quantifier",
		),
	})
}
