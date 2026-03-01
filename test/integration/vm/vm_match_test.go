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
		Case(`
LET x = MATCH 2 ( 1 => 10, 2 => 20, _ => 0, )
RETURN x
`, 20),
		CaseArray(`
FOR v IN [1, 2, 3]
	RETURN MATCH v ( 1 => "one", 2 => "two", _ => "other", )
`, []any{"one", "two", "other"}),
		Case(`RETURN (MATCH 1 ( 1 => 1, _ => 0, )) ? "yes" : "no"`, "yes"),
		Case(`RETURN LENGTH(MATCH [1,2,3] ( v => v, _ => [], ))`, 3),
		Case(`RETURN "v=" + TO_STRING(MATCH ( WHEN 1 < 2 => 2, _ => 0, ))`, "v=2"),
	})
}
