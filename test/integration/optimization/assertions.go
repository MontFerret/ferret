package optimization_test

import (
	"fmt"

	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func CastToProgram(prog any) *vm.Program {
	if p, ok := prog.(*vm.Program); ok {
		return p
	}

	panic("expected *vm.Program")
}

// ShouldEqualBytecode asserts that the bytecode of the actual program matches the expected program's bytecode.
// It compares each instruction's opcode and operands, returning an error message if any mismatch is found.
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

// ShouldUseAtMostRegisters asserts that the maximum register index used by the
// program bytecode is <= expected.
func ShouldUseAtMostRegisters(a any, e ...any) string {
	if len(e) == 0 {
		return "expected max registers value is missing"
	}

	expectedMax, ok := toInt(e[0])
	if !ok {
		return fmt.Sprintf("expected max registers must be an integer type, got %T", e[0])
	}

	actual := CastToProgram(a)
	maxReg := actual.Registers

	if maxReg > expectedMax {
		return fmt.Sprintf("expected max register <= %d, got %d", expectedMax, maxReg)
	}

	return ""
}

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case int32:
		return int(n), true
	case uint:
		return int(n), true
	case uint64:
		return int(n), true
	case uint32:
		return int(n), true
	case runtime.Int:
		return int(n), true
	default:
		return 0, false
	}
}
