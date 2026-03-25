package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestEqualityOperators(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN 1 == 1", true),
		S("RETURN 1 == 2", false),
		S("RETURN 1 != 1", false),
		S("RETURN 1 != 2", true),
		S("RETURN 1 > 1", false),
		S("RETURN 1 >= 1", true),
		S("RETURN 1 < 1", false),
		S("RETURN 1 <= 1", true),
	})
}
