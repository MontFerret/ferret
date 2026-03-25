package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestMatchExpression(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`RETURN MATCH 2 ( 1 => "one", 2 => "two", _ => "other", )`, "two"),
		S(`RETURN MATCH 5 ( v WHEN v > 3 => v * 2, _ => 0, )`, 10),
		S(`RETURN MATCH { a: 1, b: 2 } ( { a: 1, b: v } => v, _ => 0, )`, 2),
		S(`RETURN MATCH { a: 1 } ( { a: 1, b: v } => v, _ => 0, )`, 0),
		S(`RETURN MATCH { a: NONE } ( { a: NONE } => "ok", _ => "no", )`, "ok"),
		S(`RETURN MATCH ( WHEN 1 < 2 => "first", WHEN 2 < 3 => "second", _ => "default", )`, "first").Debug(),
		S(`
LET x = MATCH 2 ( 1 => 10, 2 => 20, _ => 0, )
RETURN x
`, 20),
		Array(`
FOR v IN [1, 2, 3]
	RETURN MATCH v ( 1 => "one", 2 => "two", _ => "other", )
`, []any{"one", "two", "other"}),
		S(`RETURN (MATCH 1 ( 1 => 1, _ => 0, )) ? "yes" : "no"`, "yes"),
		S(`RETURN LENGTH(MATCH [1,2,3] ( v => v, _ => [], ))`, 3),
		S(`RETURN "v=" + TO_STRING(MATCH ( WHEN 1 < 2 => 2, _ => 0, ))`, "v=2"),
	})
}
