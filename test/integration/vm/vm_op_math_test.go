package vm_test

import (
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"testing"
)

func TestMathOperators(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`RETURN 1 + 1`, 2),
		S(`RETURN 1 - 1`, 0),
		S(`RETURN 2 * 2`, 4),
		S(`RETURN 4 / 2`, 2),
		S(`RETURN 5 / 2`, 2.5),
		S(`RETURN 4.87e103`, 4.87e103),
		S(`RETURN 4.87E103`, 4.87e103),
		S(`RETURN -4.87e103`, -4.87e103),
		S(`RETURN -4.87E103`, -4.87e103),
		S(`RETURN -1 / 2`, -0.5),
		S(`RETURN 1 / (-2)`, -0.5),
		S(`RETURN (-1) / (-2)`, 0.5),
		S(`RETURN 1.0 / 2`, 0.5),
		S(`RETURN 1 / 2.0`, 0.5),
		S(`RETURN 1.0 / 2.0`, 0.5),
		S(`RETURN 5 % 2`, 1),
		S(`
LET a = 1
LET b = 2
RETURN (a - b) / 2
`, -0.5),
	})
}
