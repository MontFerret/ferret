package optimization_test

import (
	"encoding/json"
	"fmt"

	"github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func CastToProgram(prog any) *bytecode.Program {
	if p, ok := prog.(*bytecode.Program); ok {
		return p
	}

	panic("expected *vm.Program")
}

// ShouldEqualBytecode asserts that the bytecode of the actual program matches the expected program's bytecode.
// It compares each instruction's opcode and operands, returning an error message if any mismatch is found.
func ShouldEqualBytecode(a any, e ...any) string {
	actual := CastToProgram(a)
	expected := CastToProgram(e[0])

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

// ShouldCheckOpcode asserts that the actual program's bytecode contains (or does not contain) specific opcodes, or that certain opcodes appear a specific number of times, based on the provided expectation type (OpcodeExistence or OpcodeCount).
func ShouldCheckOpcode(a any, e ...any) string {
	actual := CastToProgram(a)

	switch expectation := e[0].(type) {
	case OpcodeExistence:
		return shouldCheckOpcodeExistence(actual, expectation)
	case OpcodeCount:
		return shouldCheckOpcodeCount(actual, expectation)
	default:
		panic(fmt.Sprintf("unsupported opcode expectation type: %T", e[0]))
	}

	return ""
}

func shouldCheckOpcodeCount(actual *bytecode.Program, expected OpcodeCount) string {
	current := make(map[bytecode.Opcode]int)

	for _, inst := range actual.Bytecode {
		current[inst.Opcode]++
	}

	for opcode, count := range expected.Count {
		if current[opcode] != count {
			return fmt.Sprintf("expected %d occurrences of opcode %s, got %d", count, opcode.String(), current[opcode])
		}
	}

	return ""
}

func shouldCheckOpcodeExistence(actual *bytecode.Program, expected OpcodeExistence) string {
	exists := func(op bytecode.Opcode) bool {
		for _, inst := range actual.Bytecode {
			if inst.Opcode == op {
				return true
			}
		}

		return false
	}

	for _, opcode := range expected.Exists {
		if !exists(opcode) {
			return fmt.Sprintf("expected opcode %s to be present", opcode.String())
		}
	}

	for _, opcode := range expected.NotExists {
		if exists(opcode) {
			return fmt.Sprintf("unexpected opcode %s in bytecode", opcode.String())
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

func maxRegisterIndex(bytecode []bytecode.Instruction) int {
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

func instructionUseDef(inst bytecode.Instruction) (uses []int, defs []int) {
	addUse := func(op bytecode.Operand) {
		if op != bytecode.NoopOperand && op.IsRegister() {
			uses = append(uses, op.Register())
		}
	}
	addDef := func(op bytecode.Operand) {
		if op != bytecode.NoopOperand && op.IsRegister() {
			defs = append(defs, op.Register())
		}
	}
	addRangeUses := func(start, end bytecode.Operand) {
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
	addFixedRangeUses := func(start bytecode.Operand, count int) {
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
	case bytecode.OpReturn:
		addUse(dst)
		return

	// Moves / loads.
	case bytecode.OpMove:
		addUse(src1)
		addDef(dst)
		return
	case bytecode.OpLoadConst, bytecode.OpLoadParam, bytecode.OpLoadNone, bytecode.OpLoadBool, bytecode.OpLoadZero:
		addDef(dst)
		return
	case bytecode.OpLoadArray, bytecode.OpLoadObject:
		addDef(dst)
		return
	case bytecode.OpLoadRange:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Simple arithmetic, comparisons, access.
	case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMulti, bytecode.OpDiv, bytecode.OpMod,
		bytecode.OpCmp,
		bytecode.OpEq, bytecode.OpNe, bytecode.OpGt, bytecode.OpLt, bytecode.OpGte, bytecode.OpLte,
		bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte,
		bytecode.OpAnyIn,
		bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte,
		bytecode.OpNoneIn,
		bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte,
		bytecode.OpAllIn,
		bytecode.OpIn, bytecode.OpLike, bytecode.OpRegexp,
		bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional, bytecode.OpLoadKey, bytecode.OpLoadKeyOptional,
		bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Unary ops.
	case bytecode.OpIncr, bytecode.OpDecr:
		addUse(dst)
		addDef(dst)
		return
	case bytecode.OpCastBool, bytecode.OpNegate, bytecode.OpNot, bytecode.OpFlipPositive, bytecode.OpFlipNegative, bytecode.OpLength, bytecode.OpType:
		addUse(src1)
		addDef(dst)
		return

	// Control flow.
	case bytecode.OpJumpIfFalse, bytecode.OpJumpIfTrue:
		addUse(src1)
		return
	case bytecode.OpJump:
		return

	// Dataset operations.
	case bytecode.OpDataSet, bytecode.OpDataSetCollector, bytecode.OpDataSetSorter, bytecode.OpDataSetMultiSorter:
		addDef(dst)
		return
	case bytecode.OpPush, bytecode.OpArrayPush:
		addUse(dst)
		addUse(src1)
		return
	case bytecode.OpPushKV, bytecode.OpObjectSet, bytecode.OpObjectSetConst:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		return

	// Iterators.
	case bytecode.OpIter:
		addUse(src1)
		addDef(dst)
		return
	case bytecode.OpIterValue, bytecode.OpIterKey:
		addUse(src1)
		addDef(dst)
		return
	case bytecode.OpIterLimit, bytecode.OpIterSkip:
		addUse(src1)
		addUse(src2)
		addDef(src1)
		return
	case bytecode.OpIterNext:
		addUse(src1)
		return

	// Calls.
	case bytecode.OpCall, bytecode.OpProtectedCall:
		addRangeUses(src1, src2)
		addDef(dst)
		return
	case bytecode.OpCall0, bytecode.OpProtectedCall0:
		addDef(dst)
		return
	case bytecode.OpCall1, bytecode.OpProtectedCall1:
		addUse(src1)
		addDef(dst)
		return
	case bytecode.OpCall2, bytecode.OpProtectedCall2:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case bytecode.OpCall3, bytecode.OpProtectedCall3:
		addFixedRangeUses(src1, 3)
		addDef(dst)
		return
	case bytecode.OpCall4, bytecode.OpProtectedCall4:
		addFixedRangeUses(src1, 4)
		addDef(dst)
		return

	// Stream.
	case bytecode.OpStream:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case bytecode.OpStreamIter:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case bytecode.OpDispatch:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Utility.
	case bytecode.OpClose:
		addUse(dst)
		addDef(dst)
		return
	case bytecode.OpSleep:
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
