package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

var nonTemporalArithmeticHostType = runtime.NewType("vm_test", "ArithmeticHost", func(runtime.Value) bool {
	return true
})

type nonTemporalArithmeticHost struct{}

func (nonTemporalArithmeticHost) String() string {
	return "host"
}

func (nonTemporalArithmeticHost) Hash() uint64 {
	return 1
}

func (nonTemporalArithmeticHost) Copy() runtime.Value {
	panic("arithmetic dispatch inspected Copy")
}

func (nonTemporalArithmeticHost) Type() runtime.Type {
	return nonTemporalArithmeticHostType
}

func TestNormalizedNonTemporalArithmetic(t *testing.T) {
	hostEnv := vm.WithParam("host", nonTemporalArithmeticHost{})
	RunSpecs(t, []spec.Spec{
		S(`RETURN 1 + 2`, 3),
		S(`RETURN 1 + 2.5`, 3.5),
		S(`RETURN 2.5 * 3`, 7.5),
		S(`RETURN 5.5 % 2`, 1.5),
		S(`RETURN TYPENAME(5 % 2)`, "Int"),
		S(`RETURN TYPENAME(5 % 2.0)`, "Float"),
		S(`RETURN "a" + true`, "atrue"),
		S(`RETURN NONE + "a"`, "a"),
		S(`RETURN [1, 2] + "a"`, "[1,2]a"),
		S(`RETURN "a" + @host`, "ahost").Env(hostEnv),
		S(`RETURN @host + "a"`, "hosta").Env(hostEnv),
		S(`RETURN TO_NUMBER("10") - 2`, 8),
		S(`RETURN MATCH 5.5 {5.5 => 5.5 % 2, _ => 0}`, 1.5),
		spec.NewSpec(`RETURN "10" - 2`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '-' cannot be applied to String and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN "10" * 2`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '*' cannot be applied to String and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN "10" / 2`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '/' cannot be applied to String and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN 5 / "0"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '/' cannot be applied to Int and String", ":1:8"},
		}),
		spec.NewSpec(`RETURN "10" % 2`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '%' cannot be applied to String and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN 5 % "0"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '%' cannot be applied to Int and String", ":1:8"},
		}),
		spec.NewSpec(`RETURN 5 / false`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '/' cannot be applied to Int and Boolean", ":1:8"},
		}),
		spec.NewSpec(`RETURN true + 1`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '+' cannot be applied to Boolean and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN true + 1 + "a"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '+' cannot be applied to Boolean and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN [1, 2] - 1`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '-' cannot be applied to Array and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN @host - 1`).Env(hostEnv).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '-' cannot be applied to vm_test.ArithmeticHost and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN MATCH true {true => "10" - 2, _ => 0}`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '-' cannot be applied to String and Int", ":1:28"},
		}),
		Error(`RETURN TO_NUMBER("not-a-number") - 2`),
	})
}
