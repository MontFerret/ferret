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
	})
}
