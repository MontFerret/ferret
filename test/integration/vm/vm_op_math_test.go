package vm_test

import (
	. "github.com/MontFerret/ferret/test/integration/base"
	"testing"
)

func TestMathOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`RETURN 1 + 1`, 2),
		Case(`RETURN 1 - 1`, 0),
		Case(`RETURN 2 * 2`, 4),
		Case(`RETURN 4 / 2`, 2),
		Case(`RETURN 5 % 2`, 1),
	})
}
