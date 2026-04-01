package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestSyntaxErrorsDispatch(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(`
			LET obj = NONE
			LET ok = DISPATCH IN obj
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH event name"),
		Failure(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH target"),
		Failure(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN obj WITH
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH payload"),
		Failure(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN obj OPTIONS
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH options"),
		Failure(`
			LET obj = NONE
			LET ok = -> obj
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected dispatch event before '->'",
			Hint:    `Provide an event expression, e.g. "click" -> btn.`,
		}, "Missing shorthand dispatch event"),
		Failure(`
			LET obj = NONE
			LET ok = "click" ->
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected dispatch target after '->'",
			Hint:    `Provide a dispatchable target, e.g. "click" -> btn.`,
		}, "Missing shorthand dispatch target"),
		Failure(`
			LET obj = NONE
			LET ok = "input" -> obj WITH { value: "x" }
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Dispatch shorthand does not support WITH",
			Hint:    `Use the long form instead, e.g. DISPATCH "input" IN field WITH { value: "x" }.`,
		}, "Shorthand WITH should fail syntax checks"),
		Failure(`
			LET obj = NONE
			LET ok = "click" -> obj OPTIONS { bubbles: true }
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Dispatch shorthand does not support OPTIONS",
			Hint:    `Use the long form instead, e.g. DISPATCH "click" IN btn OPTIONS { bubbles: true }.`,
		}, "Shorthand OPTIONS should fail syntax checks"),
	})
}
