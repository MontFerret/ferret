package vm_test

import (
	. "github.com/MontFerret/ferret/test/integration/base"
	"testing"
)

func TestInOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN 1 IN [1,2,3]", true),
		Case("RETURN 4 IN [1,2,3]", false),
		Case("RETURN 1 NOT IN [1,2,3]", false),
	})
}
