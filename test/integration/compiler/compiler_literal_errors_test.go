package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestNumericLiteralRangeErrors(t *testing.T) {
	const hugeInt = "999999999999999999999999999999999999999999999999"

	RunSpecs(t, []spec.Spec{
		Failure(
			"RETURN "+hugeInt,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Integer literal is out of range",
				Hint:    "Use an integer value that fits within the supported range.",
			},
			"Oversized integer literal should report a syntax diagnostic",
		),
		Failure(
			"RETURN (1 == "+hugeInt+") ? 1 : 0",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Integer literal is out of range",
				Hint:    "Use an integer value that fits within the supported range.",
			},
			"Oversized integer literal in predicate fast path should report a syntax diagnostic",
		),
		Failure(
			"RETURN 1e999",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Float literal is out of range",
				Hint:    "Use a finite float value within the supported range.",
			},
			"Oversized float literal should report a syntax diagnostic",
		),
		Failure(
			"RETURN (1 == 1e999) ? 1 : 0",
			E{
				Kind:    parserd.SyntaxError,
				Message: "Float literal is out of range",
				Hint:    "Use a finite float value within the supported range.",
			},
			"Oversized float literal in predicate fast path should report a syntax diagnostic",
		),
	})
}
