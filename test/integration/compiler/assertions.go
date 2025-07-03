package compiler_test

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/smartystreets/goconvey/convey"
)

func CastToProgram(prog any) *vm.Program {
	if p, ok := prog.(*vm.Program); ok {
		return p
	}

	panic("expected *vm.Program")
}

func ShouldEqualBytecode(e any, a ...any) string {
	expected := CastToProgram(e)
	actual := CastToProgram(a[0])

	for i := 0; i < len(expected.Bytecode); i++ {
		actualIns := actual.Bytecode[i]
		expectedIns := expected.Bytecode[i]

		if err := convey.ShouldEqual(actualIns.Opcode, expectedIns.Opcode); err != "" {
			return err
		}

		if err := convey.ShouldEqual(len(actualIns.Operands), len(expectedIns.Operands)); err != "" {
			return fmt.Sprintf("operends length mismatch at index %d: expected %d, got %d", i, len(expectedIns.Operands), len(actualIns.Operands))
		}

		for j := 0; j < len(actualIns.Operands); j++ {
			if err := convey.ShouldEqual(actualIns.Operands[j], expectedIns.Operands[j]); err != "" {
				return fmt.Sprintf("operands mismatch at index %d, operand %d: expected %s, got %s", i, j, expectedIns.Operands[j], actualIns.Operands[j])
			}
		}
	}

	return ""
}
