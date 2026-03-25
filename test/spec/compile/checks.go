package compile

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
	"github.com/google/go-cmp/cmp"
)

// EqualBytecode asserts that the bytecode of the actual program matches the expected program's bytecode.
// It compares each instruction's opcode and operands, returning an error message if any mismatch is found.
func EqualBytecode(a, e any) error {
	actual := CastToProgram(a)
	expected := CastToProgram(e)

	for i := 0; i < len(expected.Bytecode); i++ {
		actualIns := actual.Bytecode[i]
		expectedIns := expected.Bytecode[i]

		if err := assert.Equal(actualIns.Opcode, expectedIns.Opcode); err != nil {
			return err
		}

		if !cmp.Equal(len(actualIns.Operands), len(expectedIns.Operands)) {
			return fmt.Errorf("operends length mismatch at index %d: expected %d, got %d", i, len(expectedIns.Operands), len(actualIns.Operands))
		}

		for j := 0; j < len(actualIns.Operands); j++ {
			if !cmp.Equal(actualIns.Operands[j], expectedIns.Operands[j]) {
				return fmt.Errorf("operands mismatch at index %d, operand %d: expected %s, got %s", i, j, expectedIns.Operands[j], actualIns.Operands[j])
			}
		}
	}

	return nil
}

// ContainsOpcode asserts that the actual program's bytecode contains (or does not contain) specific opcodes, or that certain opcodes appear a specific number of times, based on the provided expectation type (OpcodeExistence or OpcodeCount).
func ContainsOpcode(a any, e any) error {
	actual := CastToProgram(a)

	switch expectation := e.(type) {
	case OpcodeExistence:
		return checkOpcodeExistence(actual, expectation)
	case OpcodeCount:
		return checkOpcodeCount(actual, expectation)
	default:
		return fmt.Errorf("unsupported opcode expectation type: %T", e)
	}
}

// EqualRegisters asserts that the maximum register index used by the
// program bytecode is == expected.
func EqualRegisters(a any, e any) error {
	expectedMax, ok := toInt(e)
	if !ok {
		return fmt.Errorf("expected max registers must be an integer type, got %T", e)
	}

	actual := CastToProgram(a)
	maxReg := maxRegisterIndex(actual.Bytecode)

	if maxReg != expectedMax {
		return fmt.Errorf("expected max register == %d, got %d", expectedMax, maxReg)
	}

	return nil
}

func IsCompilationError(actual, expected any) error {
	var e error

	switch ex := expected.(type) {
	case *assert.ExpectedError, assert.ExpectedError:
		err, ok := actual.(*diagnostics.Diagnostic)

		if !ok {
			err2, ok := actual.(*diagnostics.DiagnosticSet)

			if !ok {
				return fmt.Errorf("expected diagnostic set but got %T", actual)
			}

			err = err2.First()
		}

		e = assert.DiagnosticError(err, ex)
	case *assert.ExpectedMultiError, assert.ExpectedMultiError:
		err, ok := actual.(*diagnostics.DiagnosticSet)

		if !ok {
			return fmt.Errorf("expected diagnostic set but got %T", actual)
		}

		e = assert.DiagnosticErrors(err, ex)
	default:
		e = fmt.Errorf("expected a compilation error")
	}

	return e
}

func checkOpcodeCount(actual *bytecode.Program, expected OpcodeCount) error {
	current := make(map[bytecode.Opcode]int)

	for _, inst := range actual.Bytecode {
		current[inst.Opcode]++
	}

	for opcode, count := range expected.Count {
		if current[opcode] != count {
			return fmt.Errorf("expected %d occurrences of opcode %s, got %d", count, opcode.String(), current[opcode])
		}
	}

	return nil
}

func checkOpcodeExistence(actual *bytecode.Program, expected OpcodeExistence) error {
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
			return fmt.Errorf("expected opcode %s to be present", opcode.String())
		}
	}

	for _, opcode := range expected.NotExists {
		if exists(opcode) {
			return fmt.Errorf("unexpected opcode %s in bytecode", opcode.String())
		}
	}

	return nil
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
	op := inst.Opcode
	dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

	switch op {
	// No-operand terminator.
	case bytecode.OpReturn:
		addUse(dst)
		return

	// Moves / loads.
	case bytecode.OpMove, bytecode.OpMoveTracked:
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
	case bytecode.OpLoadAggregateKey:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case bytecode.OpAggregateUpdate:
		addUse(dst)
		addUse(src1)
		return
	case bytecode.OpAggregateGroupUpdate:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		return

	// Simple arithmetic, comparisons, access.
	case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod,
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

	// unary ops.
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
	case bytecode.OpCounterInc:
		addUse(dst)
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
	case bytecode.OpHCall, bytecode.OpProtectedHCall, bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
		addUse(dst)
		bytecode.VisitCallArgumentRegisters(op, src1, src2, func(reg int) {
			uses = append(uses, reg)
		})
		if op != bytecode.OpTailCall {
			addDef(dst)
		}
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
