package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestSyntaxErrorsWaitfor(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(`
			LET ok = WAITFOR EXISTS
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected expression after 'EXISTS' in WAITFOR predicate",
			Hint:    "Provide an expression after EXISTS, e.g. WAITFOR EXISTS x.",
		}, "Missing WAITFOR EXISTS expression"),
		Failure(`
			LET ok = WAITFOR NOT EXISTS
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected expression after 'NOT EXISTS' in WAITFOR predicate",
			Hint:    "Provide an expression after NOT EXISTS, e.g. WAITFOR NOT EXISTS x.",
		}, "Missing WAITFOR NOT EXISTS expression"),
		Failure(`
			LET ok = WAITFOR VALUE
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected expression after 'VALUE' in WAITFOR predicate",
			Hint:    "Provide an expression after VALUE, e.g. WAITFOR VALUE x.",
		}, "Missing WAITFOR VALUE expression"),
		Failure(`
			LET ok = WAITFOR TRUE TIMEOUT
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected value after 'TIMEOUT' in WAITFOR clause",
			Hint:    "Provide a duration, e.g. TIMEOUT 100ms.",
		}, "Missing WAITFOR TIMEOUT value"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected value after 'EVERY' in WAITFOR clause",
			Hint:    "Provide a duration, e.g. EVERY 100ms.",
		}, "Missing WAITFOR EVERY value"),
		Failure(`
			LET ok = WAITFOR TRUE BACKOFF
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected value after 'BACKOFF' in WAITFOR clause",
			Hint:    "Provide a backoff strategy, e.g. BACKOFF LINEAR.",
		}, "Missing WAITFOR BACKOFF strategy"),
		Failure(`
			LET ok = WAITFOR TRUE JITTER
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected value after 'JITTER' in WAITFOR clause",
			Hint:    "Provide a jitter value between 0 and 1, e.g. JITTER 0.2.",
		}, "Missing WAITFOR JITTER value"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY 50ms,
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing WAITFOR EVERY cap value"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY 50ms 2s
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing comma in WAITFOR EVERY cap clause"),
		Failure(`
			LET obs = {}
			LET ok = WAITFOR EVENT "test" IN obs FILTER .type == "match"
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Legacy FILTER keyword is invalid in WAITFOR EVENT"),
	})
}
