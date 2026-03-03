package vm_test

import "testing"

func TestCallArgumentLoweringSemantics(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
FUNC f2(x, y) => x + y
RETURN f2(1, 2)
`, 3, "UDF call with constant arguments should preserve semantics"),
		Case(`RETURN CONCAT("1", "2")`, "12", "Host call with constant arguments should preserve semantics"),
	})
}
