package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestBuiltinArityErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`RETURN LENGTH()`,
			E{
				Kind:    parserd.NameError,
				Message: "Function 'LENGTH' expects 1 arguments, got 0",
			},
			"LENGTH should reject missing arguments without panicking",
		),
		Failure(
			`RETURN LENGTH(1, 2)`,
			E{
				Kind:    parserd.NameError,
				Message: "Function 'LENGTH' expects 1 arguments, got 2",
			},
			"LENGTH should reject extra arguments without panicking",
		),
		Failure(
			`RETURN TYPENAME()`,
			E{
				Kind:    parserd.NameError,
				Message: "Function 'TYPENAME' expects 1 arguments, got 0",
			},
			"TYPENAME should reject missing arguments without panicking",
		),
		Failure(
			`RETURN TYPENAME(1, 2)`,
			E{
				Kind:    parserd.NameError,
				Message: "Function 'TYPENAME' expects 1 arguments, got 2",
			},
			"TYPENAME should reject extra arguments without panicking",
		),
		Failure(
			"WAIT()\nRETURN 1",
			E{
				Kind:    parserd.NameError,
				Message: "Function 'WAIT' expects 1 arguments, got 0",
			},
			"WAIT should reject missing arguments without panicking",
		),
		Failure(
			"WAIT(1, 2)\nRETURN 1",
			E{
				Kind:    parserd.NameError,
				Message: "Function 'WAIT' expects 1 arguments, got 2",
			},
			"WAIT should reject extra arguments without panicking",
		),
	})
}
