package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestSyntaxErrorsDispatch(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(`
			LET obj = NONE
			LET ok = DISPATCH IN obj
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH event name"),
		ErrorCase(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH target"),
		ErrorCase(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN obj WITH
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH payload"),
		ErrorCase(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN obj OPTIONS
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH options"),
	})
}
