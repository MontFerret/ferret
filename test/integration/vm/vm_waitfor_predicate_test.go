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
		Case(`
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
		Case(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? { foo: 1 } : {}) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty object with EXISTS"),
		CaseObject(`
			LET start = NOW()
			LET obj = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? { foo: 1 } : {}) TIMEOUT 0.5s EVERY 10ms
			RETURN obj
		`, map[string]any{"foo": 1}, "Should return object once it exists"),
		CaseArray(`
			LET start = NOW()
			LET arr = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? [1, 2] : []) TIMEOUT 0.5s EVERY 10ms
			RETURN arr
		`, []any{1, 2}, "Should return array once it exists"),
		Case(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? "ok" : "") TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty string with EXISTS"),
		Case(`
			LET start = NOW()
			LET ok = WAITFOR NOT EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? {} : { foo: 1 }) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for empty object with NOT EXISTS"),
		Case(`
			LET obj = {
				list: ["bar", "baz"],
			}
			LET ok = WAITFOR EXISTS obj.list TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non empty-array with EXISTS"),
		Case(`
			LET obj = {
				list: [],
			}
			LET ok = WAITFOR NOT EXISTS obj.list TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non empty-array with EXISTS"),
		DebugCaseNil(`
			LET token = WAITFOR VALUE NONE TIMEOUT 20ms EVERY 5ms
			RETURN token
		`, "Should return NONE on timeout for VALUE"),
	})
}
