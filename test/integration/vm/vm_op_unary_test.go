package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestUnaryOperators(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN !TRUE", false),
		S("RETURN NOT TRUE", false),
		S("RETURN !FALSE", true),
		S("RETURN -1", -1),
		S("RETURN -1.1", -1.1),
		S("RETURN +1", 1),
		S("RETURN +1.1", 1.1),
		S("LET v = 1 RETURN -v", -1),
		S("LET v = 1.1 RETURN -v", -1.1),
		S("LET v = -1 RETURN -v", 1),
		S("LET v = -1.1 RETURN -v", 1.1),
		S("LET v = -1 RETURN +v", -1),
		S("LET v = -1.1 RETURN +v", -1.1),
		Error(`RETURN !""`),
		Error(`RETURN NOT 1`),
		Error(`RETURN +"1"`),
		Error(`RETURN -"1"`),
	})
}
