package optimization_test

import (
	"encoding/json"
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

// ShouldUseEqRegisters asserts that the maximum register index used by the
// program bytecode is == expected.
func ShouldUseEqRegisters(a any, e ...any) string {
	if len(e) == 0 {
		return "expected max registers value is missing"
	}

	expectedMax, ok := toInt(e[0])
	if !ok {
		return fmt.Sprintf("expected max registers must be an integer type, got %T", e[0])
	}

	actual := CastToProgram(a)
	maxReg := maxRegisterIndex(actual.Bytecode)

	if maxReg != expectedMax {
		return fmt.Sprintf("expected max register == %d, got %d", expectedMax, maxReg)
	}

	return ""
}

// ShouldEqualJSONValue compares actual and expected by JSON-encoding both values.
// This normalizes numeric types (e.g., int vs float64) and map ordering.
func ShouldEqualJSONValue(a any, e ...any) string {
	if len(e) == 0 {
		return "expected value is missing"
	}
	if len(e) > 1 {
		return "expected single comparison value"
	}

	actualBytes, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("actual value is not JSON-marshalable: %v", err)
	}

	expectedBytes, err := json.Marshal(e[0])
	if err != nil {
		return fmt.Sprintf("expected value is not JSON-marshalable: %v", err)
	}

	return convey.ShouldEqual(string(actualBytes), string(expectedBytes))
}

func maxRegisterIndex(bytecode []vm.Instruction) int {
	maxReg := 0
	for _, inst := range bytecode {
		uses, defs := instructionUseDef(inst)
		for _, reg := range uses {
			if reg > maxReg {
				maxReg = reg
			}
		}
		for _, reg := range defs {
			if reg > maxReg {
				maxReg = reg
			}
		}
	}
	return maxReg
}

func instructionUseDef(inst vm.Instruction) (uses []int, defs []int) {
	addUse := func(op vm.Operand) {
		if op != vm.NoopOperand && op.IsRegister() {
			uses = append(uses, op.Register())
		}
	}
	addDef := func(op vm.Operand) {
		if op != vm.NoopOperand && op.IsRegister() {
			defs = append(defs, op.Register())
		}
	}
	addRangeUses := func(start, end vm.Operand) {
		if !start.IsRegister() || !end.IsRegister() {
			return
		}
		startReg := start.Register()
		endReg := end.Register()
		if startReg <= 0 || endReg < startReg {
			return
		}
		for r := startReg; r <= endReg; r++ {
			uses = append(uses, r)
		}
	}
	addFixedRangeUses := func(start vm.Operand, count int) {
		if !start.IsRegister() {
			return
		}
		startReg := start.Register()
		if startReg <= 0 || count <= 0 {
			return
		}
		for r := startReg; r < startReg+count; r++ {
			uses = append(uses, r)
		}
	}

	op := inst.Opcode
	dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

	switch op {
	// No-operand terminator.
	case vm.OpReturn:
		addUse(dst)
		return

	// Moves / loads.
	case vm.OpMove:
		addUse(src1)
		addDef(dst)
		return
	case vm.OpLoadConst, vm.OpLoadParam, vm.OpLoadNone, vm.OpLoadBool, vm.OpLoadZero:
		addDef(dst)
		return
	case vm.OpLoadArray, vm.OpLoadObject:
		addDef(dst)
		return
	case vm.OpLoadRange:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Simple arithmetic, comparisons, access.
	case vm.OpAdd, vm.OpSub, vm.OpMulti, vm.OpDiv, vm.OpMod,
		vm.OpCmp,
		vm.OpEq, vm.OpNe, vm.OpGt, vm.OpLt, vm.OpGte, vm.OpLte,
		vm.OpAnyEq, vm.OpAnyNe, vm.OpAnyGt, vm.OpAnyGte, vm.OpAnyLt, vm.OpAnyLte,
		vm.OpAnyIn,
		vm.OpNoneEq, vm.OpNoneNe, vm.OpNoneGt, vm.OpNoneGte, vm.OpNoneLt, vm.OpNoneLte,
		vm.OpNoneIn,
		vm.OpAllEq, vm.OpAllNe, vm.OpAllGt, vm.OpAllGte, vm.OpAllLt, vm.OpAllLte,
		vm.OpAllIn,
		vm.OpIn, vm.OpLike, vm.OpRegexp,
		vm.OpLoadIndex, vm.OpLoadIndexOptional, vm.OpLoadKey, vm.OpLoadKeyOptional,
		vm.OpLoadProperty, vm.OpLoadPropertyOptional:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Unary ops.
	case vm.OpIncr, vm.OpDecr:
		addUse(dst)
		addDef(dst)
		return
	case vm.OpCastBool, vm.OpNegate, vm.OpNot, vm.OpFlipPositive, vm.OpFlipNegative, vm.OpLength, vm.OpType:
		addUse(src1)
		addDef(dst)
		return

	// Control flow.
	case vm.OpJumpIfFalse, vm.OpJumpIfTrue:
		addUse(src1)
		return
	case vm.OpJump:
		return

	// Dataset operations.
	case vm.OpDataSet, vm.OpDataSetCollector, vm.OpDataSetSorter, vm.OpDataSetMultiSorter:
		addDef(dst)
		return
	case vm.OpPush, vm.OpArrayPush:
		addUse(dst)
		addUse(src1)
		return
	case vm.OpPushKV, vm.OpObjectSet:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		return

	// Iterators.
	case vm.OpIter:
		addUse(src1)
		addDef(dst)
		return
	case vm.OpIterValue, vm.OpIterKey:
		addUse(src1)
		addDef(dst)
		return
	case vm.OpIterLimit, vm.OpIterSkip:
		addUse(src1)
		addUse(src2)
		addDef(src1)
		return
	case vm.OpIterNext:
		addUse(src1)
		return

	// Calls.
	case vm.OpCall, vm.OpProtectedCall:
		addRangeUses(src1, src2)
		addDef(dst)
		return
	case vm.OpCall0, vm.OpProtectedCall0:
		addDef(dst)
		return
	case vm.OpCall1, vm.OpProtectedCall1:
		addUse(src1)
		addDef(dst)
		return
	case vm.OpCall2, vm.OpProtectedCall2:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case vm.OpCall3, vm.OpProtectedCall3:
		addFixedRangeUses(src1, 3)
		addDef(dst)
		return
	case vm.OpCall4, vm.OpProtectedCall4:
		addFixedRangeUses(src1, 4)
		addDef(dst)
		return

	// Stream.
	case vm.OpStream:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case vm.OpStreamIter:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Utility.
	case vm.OpClose:
		addUse(dst)
		addDef(dst)
		return
	case vm.OpSleep:
		addUse(dst)
		return
	}

	return
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
