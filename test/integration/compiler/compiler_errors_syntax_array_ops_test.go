package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestSyntaxErrorsArrayOperators(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`RETURN doc[~]`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query literal after '~'",
				Hint:    "Provide a query literal, e.g. doc[~ css`...`].",
			},
			"Missing query literal after '~'",
		),
		Failure(
			`RETURN doc[~ 'x']`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query literal after '~'",
				Hint:    "Provide a query literal, e.g. doc[~ css`...`].",
			},
			"Missing query type before literal",
		),
		Failure(
			`LET doc = {} RETURN doc[~ css()]`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query string after 'css'",
				Hint:    "Provide a query string, e.g. doc[~ css`...`].",
			},
			"Missing query string after type",
		),
		Failure(
			`LET doc = {} RETURN doc[~ css"x"`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected query type before query literal",
				Hint:    "Provide a type name before the query string, e.g. doc[~ css`...`].",
			},
			"Missing query type before query literal",
		),
		Failure(
			`RETURN [1, 2][* RETURN]`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'RETURN' in array operator",
				Hint:    "Provide a projection expression, e.g. [* RETURN .].",
			},
			"Missing inline RETURN expression",
		),
		Failure(
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
