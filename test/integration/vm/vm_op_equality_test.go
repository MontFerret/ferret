package vm_test

import (
	. "github.com/MontFerret/ferret/test/integration/base"
	"testing"
)

func TestEqualityOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN 1 == 1", true),
		Case("RETURN 1 == 2", false),
		Case("RETURN 1 != 1", false),
		Case("RETURN 1 != 2", true),
		Case("RETURN 1 > 1", false),
		Case("RETURN 1 >= 1", true),
		Case("RETURN 1 < 1", false),
		Case("RETURN 1 <= 1", true),
	})
}
