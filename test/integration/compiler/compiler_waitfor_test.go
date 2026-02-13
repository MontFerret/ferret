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
		ErrorCase(`
			LET ok = WAITFOR TRUE JITTER 1.5
			RETURN ok
		`, E{
			Message: "JITTER must be between 0 and 1",
			Hint:    "Use a value between 0 and 1, e.g. JITTER 0.2.",
		}, "Out-of-range JITTER should fail compilation"),
	})
}
