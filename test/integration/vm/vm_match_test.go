package vm_test

import "testing"

func TestMatchExpression(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`RETURN MATCH 2 ( 1 => "one", 2 => "two", _ => "other", )`, "two"),
		Case(`RETURN MATCH 5 ( v WHEN v > 3 => v * 2, _ => 0, )`, 10),
		Case(`RETURN MATCH { a: 1, b: 2 } ( { a: 1, b: v } => v, _ => 0, )`, 2),
		Case(`RETURN MATCH { a: 1 } ( { a: 1, b: v } => v, _ => 0, )`, 0),
		Case(`RETURN MATCH { a: NONE } ( { a: NONE } => "ok", _ => "no", )`, "ok"),
		DebugCase(`RETURN MATCH ( WHEN 1 < 2 => "first", WHEN 2 < 3 => "second", _ => "default", )`, "first"),
	})
}
