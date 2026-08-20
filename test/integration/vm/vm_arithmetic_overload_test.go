package vm_test

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestHostArithmeticOperatorOverloads(t *testing.T) {
	host := newVMArithmeticOverloadValue("host")
	hostEnv := vm.WithParam("host", host)

	RunSpecs(t, []spec.Spec{
		S(`RETURN @host + 10`, "Add").Env(hostEnv),
		S(`RETURN @host - 10`, "Subtract").Env(hostEnv),
		S(`RETURN @host * 10`, "Multiply").Env(hostEnv),
		S(`RETURN @host / 10`, "Divide").Env(hostEnv),
		S(`RETURN @host % 10`, "Mod").Env(hostEnv),
		S(`RETURN 10 + @host`, "RightAdd").Env(hostEnv),
		S(`RETURN 10 - @host`, "RightSubtract").Env(hostEnv),
		S(`RETURN 10 * @host`, "RightMultiply").Env(hostEnv),
		S(`RETURN 10 / @host`, "RightDivide").Env(hostEnv),
		S(`RETURN 10 % @host`, "RightMod").Env(hostEnv),
	})
}

func TestHostArithmeticOperatorNegotiation(t *testing.T) {
	left := newVMArithmeticOverloadValue("left")
	left.responses["Add"] = vmArithmeticOverloadResponse{
		value: runtime.None,
		err:   runtime.ErrUnsupportedOperands,
	}
	right := newVMArithmeticOverloadValue("right")
	right.responses["RightAdd"] = vmArithmeticOverloadResponse{value: runtime.NewString("fallback")}
	env := vm.WithParams(map[string]runtime.Value{"left": left, "right": right})

	RunSpecs(t, []spec.Spec{
		S(`RETURN @left + @right`, "fallback").Env(env),
	})
}

func TestHostArithmeticOperatorErrors(t *testing.T) {
	hostErr := errors.New("currency mismatch")
	failing := newVMArithmeticOverloadValue("failing")
	failing.responses["Subtract"] = vmArithmeticOverloadResponse{value: runtime.None, err: hostErr}
	failingEnv := vm.WithParams(map[string]runtime.Value{
		"left":  failing,
		"right": newVMArithmeticOverloadValue("right"),
	})

	decliningLeft := newVMArithmeticOverloadValue("left")
	decliningLeft.responses["Divide"] = vmArithmeticOverloadResponse{
		value: runtime.None,
		err:   runtime.ErrUnsupportedOperands,
	}
	decliningRight := newVMArithmeticOverloadValue("right")
	decliningRight.responses["RightDivide"] = vmArithmeticOverloadResponse{
		value: runtime.None,
		err:   runtime.ErrUnsupportedOperands,
	}
	decliningEnv := vm.WithParams(map[string]runtime.Value{
		"left":  decliningLeft,
		"right": decliningRight,
	})

	RunSpecs(t, []spec.Spec{
		spec.NewSpec(`RETURN @left - @right`).Env(failingEnv).Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "currency mismatch", Contains: []string{":1:8"}},
		),
		spec.NewSpec(`RETURN @left / @right`).Env(decliningEnv).Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{
				Message: "invalid operation",
				Contains: []string{
					"operator '/' cannot be applied to vm_test.ArithmeticOverload and vm_test.ArithmeticOverload",
					":1:8",
				},
			},
		),
	})
}

func TestHostArithmeticPrimitiveCapabilitiesRemainIndependent(t *testing.T) {
	type typedRuntimeValue interface {
		runtime.Value
		runtime.Typed
	}

	addImplementation := newVMArithmeticOverloadValue("add")
	addOnly := struct {
		typedRuntimeValue
		runtime.Addable
	}{addImplementation, addImplementation}

	multiplyImplementation := newVMArithmeticOverloadValue("multiply")
	multiplyOnly := struct {
		typedRuntimeValue
		runtime.Multipliable
	}{multiplyImplementation, multiplyImplementation}

	RunSpecs(t, []spec.Spec{
		S(`RETURN @value + 1`, "Add").Env(vm.WithParam("value", addOnly)),
		Error(`RETURN @value - 1`).Env(vm.WithParam("value", addOnly)),
		S(`RETURN @value * 2`, "Multiply").Env(vm.WithParam("value", multiplyOnly)),
		Error(`RETURN @value / 2`).Env(vm.WithParam("value", multiplyOnly)),
		Error(`RETURN @value % 2`).Env(vm.WithParam("value", multiplyOnly)),
	})
}

func TestHostArithmeticPreservesNativePrecedence(t *testing.T) {
	hostEnv := vm.WithParam("host", newVMArithmeticOverloadValue("host"))

	RunSpecs(t, []spec.Spec{
		S(`RETURN "left:" + @host`, "left:host").Env(hostEnv),
		S(`RETURN @host + ":right"`, "host:right").Env(hostEnv),
		S(`RETURN 1 + 2`, 3),
		S(`RETURN 5 / 2`, 2.5),
		S(`RETURN 5 % 2`, 1),
		S(`RETURN 1s + 500ms`, "1.5s"),
	})
}
