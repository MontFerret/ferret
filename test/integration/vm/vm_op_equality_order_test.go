package vm_test

import (
	"testing"
)

func TestEqualityOrderOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN true > NONE", true),

		// null < ...
		Case("RETURN NONE < false", true),
		Case("RETURN NONE < true", true),
		Case("RETURN NONE < 0", true),
		Case("RETURN NONE < ''", true),
		Case("RETURN NONE < ' '", true),
		Case("RETURN NONE < '0'", true),
		Case("RETURN NONE < 'abc'", true),
		Case("RETURN NONE < []", true),
		Case("RETURN NONE < {}", true),

		// false < ...
		Case("RETURN false < true", true),
		Case("RETURN false < 0", true),
		Case("RETURN false < ''", true),
		Case("RETURN false < ' '", true),
		Case("RETURN false < '0'", true),
		Case("RETURN false < 'abc'", true),
		Case("RETURN false < []", true),
		Case("RETURN false < {}", true),

		// true < ...
		Case("RETURN true < 0", true),
		Case("RETURN true < ''", true),
		Case("RETURN true < ' '", true),
		Case("RETURN true < '0'", true),
		Case("RETURN true < 'abc'", true),
		Case("RETURN true < []", true),
		Case("RETURN true < {}", true),

		// 0 < ...
		Case("RETURN 0 < ''", true),
		Case("RETURN 0 < ' '", true),
		Case("RETURN 0 < '0'", true),
		Case("RETURN 0 < 'abc'", true),
		Case("RETURN 0 < []", true),
		Case("RETURN 0 < {}", true),

		// '' < ...
		Case("RETURN '' < ' '", true),
		Case("RETURN '' < '0'", true),
		Case("RETURN '' < 'abc'", true),
		Case("RETURN '' < []", true),
		Case("RETURN '' < {}", true),

		// [] < {}
		Case("RETURN [] < {}", true),
	})
}
