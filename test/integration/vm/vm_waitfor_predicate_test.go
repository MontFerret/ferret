package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestWaitforPredicate(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
			LET start = NOW()
			LET ok = WAITFOR (DATE_DIFF(start, NOW(), "f") >= 30) TIMEOUT 0.5s EVERY 0.01s BACKOFF LINEAR
			RETURN ok
		`, true, "Should wait until predicate becomes true"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR (DATE_DIFF(start, NOW(), "f") >= 30) TIMEOUT 0.5s EVERY 10ms, 30ms BACKOFF LINEAR JITTER 0.2
			RETURN ok
		`, true, "Should support EVERY cap with JITTER"),
		S(`
			LET ok = WAITFOR FALSE TIMEOUT 50ms EVERY 10ms
			RETURN ok
		`, false, "Should return false on timeout"),
		S(`
			LET ok = WAITFOR FALSE TIMEOUT 80ms EVERY 10ms, 10ms BACKOFF EXPONENTIAL JITTER 0.5
			RETURN ok
		`, false, "Should honor timeout with cap and jitter"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? [1] : []) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty value with EXISTS"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR NOT EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? "" : "x") TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for empty value with NOT EXISTS"),
		S(`
			LET start = NOW()
			LET token = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? "ok" : NONE) TIMEOUT 0.5s EVERY 10ms
			RETURN token
		`, "ok", "Should return value once it exists"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? { foo: 1 } : {}) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty object with EXISTS"),
		Object(`
			LET start = NOW()
			LET obj = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? { foo: 1 } : {}) TIMEOUT 0.5s EVERY 10ms
			RETURN obj
		`, map[string]any{"foo": 1}, "Should return object once it exists"),
		Array(`
			LET start = NOW()
			LET arr = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? [1, 2] : []) TIMEOUT 0.5s EVERY 10ms
			RETURN arr
		`, []any{1, 2}, "Should return array once it exists"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? "ok" : "") TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty string with EXISTS"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR NOT EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? {} : { foo: 1 }) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for empty object with NOT EXISTS"),
		S(`
			LET obj = {
				list: ["bar", "baz"],
			}
			LET ok = WAITFOR EXISTS obj.list TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non empty-array with EXISTS"),
		S(`
			LET obj = {
				list: [],
			}
			LET ok = WAITFOR NOT EXISTS obj.list TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non empty-array with EXISTS"),
		Nil(`
			LET token = WAITFOR VALUE NONE TIMEOUT 20ms EVERY 5ms
			RETURN token
		`, "Should return NONE on timeout for VALUE"),
	})
}
