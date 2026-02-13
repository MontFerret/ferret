package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func TestSyntaxErrorsWaitfor(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(`
			LET ok = WAITFOR EXISTS
			RETURN ok
		`, E{
			Kind:    compiler.SyntaxError,
			Message: "Expected expression after 'EXISTS' in WAITFOR predicate",
			Hint:    "Provide an expression after EXISTS, e.g. WAITFOR EXISTS x.",
		}, "Missing WAITFOR EXISTS expression"),
		ErrorCase(`
			LET ok = WAITFOR NOT EXISTS
			RETURN ok
		`, E{
			Kind:    compiler.SyntaxError,
			Message: "Expected expression after 'NOT EXISTS' in WAITFOR predicate",
			Hint:    "Provide an expression after NOT EXISTS, e.g. WAITFOR NOT EXISTS x.",
		}, "Missing WAITFOR NOT EXISTS expression"),
		ErrorCase(`
			LET ok = WAITFOR VALUE
			RETURN ok
		`, E{
			Kind:    compiler.SyntaxError,
			Message: "Expected expression after 'VALUE' in WAITFOR predicate",
			Hint:    "Provide an expression after VALUE, e.g. WAITFOR VALUE x.",
		}, "Missing WAITFOR VALUE expression"),
		ErrorCase(`
			LET ok = WAITFOR TRUE TIMEOUT
			RETURN ok
		`, E{
			Kind:    compiler.SyntaxError,
			Message: "Expected value after 'TIMEOUT' in WAITFOR clause",
			Hint:    "Provide a duration, e.g. TIMEOUT 100ms.",
		}, "Missing WAITFOR TIMEOUT value"),
		ErrorCase(`
			LET ok = WAITFOR TRUE EVERY
			RETURN ok
		`, E{
			Kind:    compiler.SyntaxError,
			Message: "Expected value after 'EVERY' in WAITFOR clause",
			Hint:    "Provide a duration, e.g. EVERY 100ms.",
		}, "Missing WAITFOR EVERY value"),
		ErrorCase(`
			LET ok = WAITFOR TRUE BACKOFF
			RETURN ok
		`, E{
			Kind:    compiler.SyntaxError,
			Message: "Expected value after 'BACKOFF' in WAITFOR clause",
			Hint:    "Provide a backoff strategy, e.g. BACKOFF LINEAR.",
		}, "Missing WAITFOR BACKOFF strategy"),
	})
}
