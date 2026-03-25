package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestCallArgumentLoweringSemantics(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
FUNC f2(x, y) => x + y
RETURN f2(1, 2)
`, 3, "UDF call with constant arguments should preserve semantics"),
		S(`RETURN CONCAT("1", "2")`, "12", "Host call with constant arguments should preserve semantics"),
	})
}
