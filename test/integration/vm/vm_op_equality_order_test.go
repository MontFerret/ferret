package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestEqualityOrderOperators(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN true > NONE", true),

		// null < ...
		S("RETURN NONE < false", true),
		S("RETURN NONE < true", true),
		S("RETURN NONE < 0", true),
		S("RETURN NONE < ''", true),
		S("RETURN NONE < ' '", true),
		S("RETURN NONE < '0'", true),
		S("RETURN NONE < 'abc'", true),
		S("RETURN NONE < []", true),
		S("RETURN NONE < {}", true),

		// false < ...
		S("RETURN false < true", true),
		S("RETURN false < 0", true),
		S("RETURN false < ''", true),
		S("RETURN false < ' '", true),
		S("RETURN false < '0'", true),
		S("RETURN false < 'abc'", true),
		S("RETURN false < []", true),
		S("RETURN false < {}", true),

		// true < ...
		S("RETURN true < 0", true),
		S("RETURN true < ''", true),
		S("RETURN true < ' '", true),
		S("RETURN true < '0'", true),
		S("RETURN true < 'abc'", true),
		S("RETURN true < []", true),
		S("RETURN true < {}", true),

		// 0 < ...
		S("RETURN 0 < ''", true),
		S("RETURN 0 < ' '", true),
		S("RETURN 0 < '0'", true),
		S("RETURN 0 < 'abc'", true),
		S("RETURN 0 < []", true),
		S("RETURN 0 < {}", true),

		// '' < ...
		S("RETURN '' < ' '", true),
		S("RETURN '' < '0'", true),
		S("RETURN '' < 'abc'", true),
		S("RETURN '' < []", true),
		S("RETURN '' < {}", true),

		// [] < {}
		S("RETURN [] < {}", true),
	})
}
