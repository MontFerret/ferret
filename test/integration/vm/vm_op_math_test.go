package vm_test

import (
	"testing"
)

func TestMathOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`RETURN 1 + 1`, 2),
		Case(`RETURN 1 - 1`, 0),
		Case(`RETURN 2 * 2`, 4),
		Case(`RETURN 4 / 2`, 2),
		Case(`RETURN 5 / 2`, 2.5),
		Case(`RETURN -1 / 2`, -0.5),
		Case(`RETURN 1 / (-2)`, -0.5),
		Case(`RETURN (-1) / (-2)`, 0.5),
		Case(`RETURN 1.0 / 2`, 0.5),
		Case(`RETURN 1 / 2.0`, 0.5),
		Case(`RETURN 1.0 / 2.0`, 0.5),
		Case(`RETURN 5 % 2`, 1),
		Case(`
LET a = 1
LET b = 2
RETURN (a - b) / 2
`, -0.5),
	})
}
