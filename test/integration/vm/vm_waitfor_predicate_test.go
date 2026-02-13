package vm_test

import "testing"

func TestWaitforPredicate(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
			LET start = NOW()
			LET ok = WAITFOR (DATE_DIFF(start, NOW(), "f") >= 30) TIMEOUT 0.5s EVERY 0.01s BACKOFF LINEAR
			RETURN ok
		`, true, "Should wait until predicate becomes true"),
		Case(`
			LET ok = WAITFOR FALSE TIMEOUT 50ms EVERY 10ms
			RETURN ok
		`, false, "Should return false on timeout"),
		DebugCase(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? [1] : []) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty value with EXISTS"),
		Case(`
			LET start = NOW()
			LET ok = WAITFOR NOT EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? "" : "x") TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for empty value with NOT EXISTS"),
		Case(`
			LET start = NOW()
			LET token = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? "ok" : NONE) TIMEOUT 0.5s EVERY 10ms
			RETURN token
		`, "ok", "Should return value once it exists"),
		CaseNil(`
			LET token = WAITFOR VALUE NONE TIMEOUT 20ms EVERY 5ms
			RETURN token
		`, "Should return NONE on timeout for VALUE"),
	})
}
