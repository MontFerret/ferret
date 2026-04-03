package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestWaitforCompilationErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(`
			LET ok = WAITFOR TRUE BACKOFF UNKNOWN
			RETURN ok
		`, E{
			Message: "Unknown BACKOFF strategy",
			Hint:    "Use one of: NONE, LINEAR, EXPONENTIAL.",
		}, "Unknown BACKOFF strategy should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE OR THROW
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "OR THROW should fail as a syntax error"),
		Failure(`
			LET ok = WAITFOR TRUE JITTER 1.5
			RETURN ok
		`, E{
			Message: "JITTER must be between 0 and 1",
			Hint:    "Use a value between 0 and 1, e.g. JITTER 0.2.",
		}, "Out-of-range JITTER should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE TIMEOUT 1e999s
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Duration literal is out of range",
			Hint:    "Use a duration value that stays within the supported range, e.g. 100ms, 2s, or 1.5m.",
		}, "Out-of-range WAITFOR TIMEOUT duration should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY 1e999s
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Duration literal is out of range",
			Hint:    "Use a duration value that stays within the supported range, e.g. 100ms, 2s, or 1.5m.",
		}, "Out-of-range WAITFOR EVERY duration should fail compilation"),
	})
}
