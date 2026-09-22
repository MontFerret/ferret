package vm_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestIntegerLiteralBoundaries(t *testing.T) {
	for _, value := range []int64{math.MinInt32 - 1, math.MaxInt32 + 1, math.MaxUint32 + 1, 1<<53 + 1, math.MaxInt64} {
		literal := strconv.FormatInt(value, 10)
		t.Run(literal, func(t *testing.T) {
			RunSpecs(t, []spec.Spec{
				S("RETURN TO_STRING("+literal+")", literal),
				S("RETURN TYPENAME("+literal+")", "Int"),
				S(`RETURN "" + (`+literal+")", literal),
				S("RETURN @value == ("+literal+")", true).Env(vm.WithParam("value", runtime.NewInt64(value))),
			})
		})
	}

	RunSpecs(t, []spec.Spec{
		S("RETURN TO_STRING((-9223372036854775807) - 1)", "-9223372036854775808"),
		S("RETURN WAITFOR 2147483648 TIMEOUT 1ms", true),
	})
}
