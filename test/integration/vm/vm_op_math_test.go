package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestMathOperators(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`RETURN 1 + 1`, 2),
		S(`RETURN 1 - 1`, 0),
		S(`RETURN 2 * 2`, 4),
		S(`RETURN 4 / 2`, 2),
		S(`RETURN 5 / 2`, 2.5),
		S(`RETURN 4.87e103`, 4.87e103),
		S(`RETURN 4.87E103`, 4.87e103),
		S(`RETURN -4.87e103`, -4.87e103),
		S(`RETURN -4.87E103`, -4.87e103),
		S(`RETURN -1 / 2`, -0.5),
		S(`RETURN 1 / (-2)`, -0.5),
		S(`RETURN (-1) / (-2)`, 0.5),
		S(`RETURN 1.0 / 2`, 0.5),
		S(`RETURN 1 / 2.0`, 0.5),
		S(`RETURN 1.0 / 2.0`, 0.5),
		S(`RETURN 5 % 2`, 1),
		S(`RETURN "a" + 1`, "a1"),
		S(`
LET a = 1
LET b = 2
RETURN (a - b) / 2
`, -0.5),
	})
}

func TestMathOperatorsRejectDynamicNonNumericOperands(t *testing.T) {
	tests := []struct {
		value runtime.Value
		name  string
		query string
	}{
		{name: "subtract string", query: `RETURN @value - 1`, value: runtime.NewString("3")},
		{name: "multiply array", query: `RETURN @value * 2`, value: runtime.NewArrayWith(runtime.NewInt(120), runtime.NewInt(45), runtime.NewInt(300))},
		{name: "multiply string", query: `RETURN @value * 2`, value: runtime.NewString("3")},
		{name: "multiply boolean", query: `RETURN @value * 2`, value: runtime.True},
		{name: "divide object", query: `RETURN @value / 2`, value: runtime.NewObject()},
		{name: "modulo boolean", query: `RETURN @value % 2`, value: runtime.True},
		{name: "unary negative string", query: `RETURN -@value`, value: runtime.NewString("3")},
		{name: "unary positive boolean", query: `RETURN +@value`, value: runtime.True},
	}

	specs := make([]spec.Spec, 0, len(tests))
	for _, test := range tests {
		specs = append(specs,
			spec.NewSpec(test.query, test.name).
				Env(vm.WithParam("value", test.value)).
				Expect().
				ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "invalid type"}),
		)
	}

	RunSpecs(t, specs)
}
