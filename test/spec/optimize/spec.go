package optimize

import (
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
	"github.com/MontFerret/ferret/v2/test/spec/compile"
	"github.com/MontFerret/ferret/v2/test/spec/exec"
)

func Opcode[T compile.OpcodeExpectation](expression string, expectation T, out any, desc ...string) spec.Spec {
	return compile.Opcode(expression, expectation, desc...).Expect().Exec(assert.ShouldEqual, out)
}

func OpcodeErr[T compile.OpcodeExpectation](expression string, expectation T, out any, desc ...string) spec.Spec {
	return compile.Opcode(expression, expectation, desc...).Expect().ExecError(assert.ShouldBeError, out)
}

func Registers(expression string, num int, output any, desc ...string) spec.Spec {
	return compile.Registers(expression, num, desc...).Expect().Exec(assert.ShouldEqual, output)
}

func RegistersArray(expression string, num int, output []any, desc ...string) spec.Spec {
	return spec.Compose(spec.NewSpec(expression, desc...), compile.Registers(expression, num), exec.Array(expression, output))
}

func RegistersObject(expression string, num int, output map[string]any, desc ...string) spec.Spec {
	return spec.Compose(spec.NewSpec(expression, desc...), compile.Registers(expression, num), exec.Object(expression, output))
}
