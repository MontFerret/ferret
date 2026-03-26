package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestInOperator(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN 1 IN [1,2,3]", true),
		S("RETURN 4 IN [1,2,3]", false),
		S("RETURN 1 NOT IN [1,2,3]", false),
	})
}
