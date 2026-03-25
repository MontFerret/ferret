package vm_test

import (
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"testing"
)

func TestInOperator(t *testing.T) {
	RunSpecs(t, []Spec{
		S("RETURN 1 IN [1,2,3]", true),
		S("RETURN 4 IN [1,2,3]", false),
		S("RETURN 1 NOT IN [1,2,3]", false),
	})
}
