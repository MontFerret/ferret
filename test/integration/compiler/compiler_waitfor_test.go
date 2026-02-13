package compiler_test

import "testing"

func TestWaitforCompilationErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(`
			LET ok = WAITFOR TRUE BACKOFF UNKNOWN
			RETURN ok
		`, E{
			Message: "Unknown BACKOFF strategy",
			Hint:    "Use one of: NONE, LINEAR, EXPONENTIAL.",
		}, "Unknown BACKOFF strategy should fail compilation"),
		ErrorCase(`
			LET ok = WAITFOR TRUE OR THROW
			RETURN ok
		`, E{
			Message: "OR THROW is not supported",
			Hint:    "Remove OR THROW and handle timeouts explicitly.",
		}, "OR THROW should fail compilation"),
	})
}
