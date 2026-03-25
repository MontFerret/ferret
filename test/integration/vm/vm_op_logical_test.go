package vm_test

import (
	"context"
	"fmt"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestLogicalOperators(t *testing.T) {
	RunSpecs(t, []Spec{
		S("RETURN 1 AND 0", 0),
		S("RETURN 1 AND 1", 1),
		S("RETURN 2 > 1 AND 1 > 0", true),
		S("RETURN NONE && true", nil),
		S("RETURN '' && true", ""),
		S("RETURN true && 23", 23),
		S("RETURN 1 OR 0", 1),
		S("RETURN 0 OR 1", 1),
		S("RETURN 2 OR 1", 2),
		S("RETURN 2 > 1 OR 1 > 0", true),
		S("RETURN 2 < 1 OR 1 > 0", true),
		S("RETURN 1 || 7", 1),
		S("RETURN 0 || 7", 7),
		S("RETURN NONE || 'foo'", "foo"),
		S("RETURN '' || 'foo'", "foo"),
		S(`RETURN ERROR()? || 'boo'`, "boo"),
		S(`RETURN !ERROR()? && TRUE`, true),
		S(`LET u = { valid: false } RETURN u.valid || TRUE`, true),
	}, vm.WithFunction("ERROR", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.None, fmt.Errorf("test")
	}))
}
