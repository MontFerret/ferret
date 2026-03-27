package bytecode

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ValidateProgram checks that a bytecode program is structurally valid for
// persistence and later execution. ISA compatibility is validated separately by
// loaders and the VM.
func ValidateProgram(program *Program) error {
	if program == nil {
		return fmt.Errorf("%w: program is nil", ErrInvalidProgram)
	}

	if program.Registers < 0 {
		return fmt.Errorf("%w: register count must be >= 0, got %d", ErrInvalidProgram, program.Registers)
	}

	if len(program.Bytecode) == 0 {
		return fmt.Errorf("%w: bytecode must not be empty", ErrInvalidProgram)
	}

	if err := validateFunctions(program); err != nil {
		return err
	}

	if err := validateMetadata(program); err != nil {
		return err
	}

	if err := validateCatchTable(program); err != nil {
		return err
	}

	return validateInstructions(program)
}

func validateFunctions(program *Program) error {
	for name, arity := range program.Functions.Host {
		if name == "" {
			return fmt.Errorf("%w: host function name must not be empty", ErrInvalidProgram)
		}

		if arity < 0 {
			return fmt.Errorf("%w: host function %q has negative arity %d", ErrInvalidProgram, name, arity)
		}
	}

	bytecodeLen := len(program.Bytecode)

	for i, udf := range program.Functions.UserDefined {
		if udf.Name == "" {
			return fmt.Errorf("%w: udf %d has empty name", ErrInvalidProgram, i)
		}

		if udf.Entry < 0 || udf.Entry >= bytecodeLen {
			return fmt.Errorf("%w: udf %q entry pc %d out of range", ErrInvalidProgram, udf.Name, udf.Entry)
		}

		if udf.Registers < 0 {
			return fmt.Errorf("%w: udf %q has negative register count %d", ErrInvalidProgram, udf.Name, udf.Registers)
		}

		if udf.Params < 0 {
			return fmt.Errorf("%w: udf %q has negative param count %d", ErrInvalidProgram, udf.Name, udf.Params)
		}
	}

	return nil
}

func validateMetadata(program *Program) error {
	bytecodeLen := len(program.Bytecode)

	if len(program.Metadata.AggregateSelectorSlots) > 0 {
		if len(program.Metadata.AggregateSelectorSlots) != bytecodeLen {
			return fmt.Errorf("%w: aggregate selector slot metadata length %d does not match bytecode length %d", ErrInvalidProgram, len(program.Metadata.AggregateSelectorSlots), bytecodeLen)
		}

		for pc, slot := range program.Metadata.AggregateSelectorSlots {
			if slot < -1 {
				return fmt.Errorf("%w: invalid aggregate selector slot %d at pc %d", ErrInvalidProgram, slot, pc)
			}
		}
	}

	if len(program.Metadata.MatchFailTargets) > 0 {
		if len(program.Metadata.MatchFailTargets) != bytecodeLen {
			return fmt.Errorf("%w: match fail target metadata length %d does not match bytecode length %d", ErrInvalidProgram, len(program.Metadata.MatchFailTargets), bytecodeLen)
		}

		for pc, target := range program.Metadata.MatchFailTargets {
			if target < -1 || target >= bytecodeLen {
				return fmt.Errorf("%w: invalid match fail target %d at pc %d", ErrInvalidProgram, target, pc)
			}
		}
	}

	if len(program.Metadata.DebugSpans) > 0 {
		if len(program.Metadata.DebugSpans) != bytecodeLen {
			return fmt.Errorf("%w: debug span metadata length %d does not match bytecode length %d", ErrInvalidProgram, len(program.Metadata.DebugSpans), bytecodeLen)
		}

		for pc, span := range program.Metadata.DebugSpans {
			if span.Start == -1 && span.End == -1 {
				continue
			}

			if span.Start < 0 || span.End < span.Start {
				return fmt.Errorf("%w: invalid debug span at pc %d", ErrInvalidProgram, pc)
			}
		}
	}

	for pc, name := range program.Metadata.Labels {
		if pc < 0 || pc >= bytecodeLen {
			return fmt.Errorf("%w: label %q references pc %d out of range", ErrInvalidProgram, name, pc)
		}
	}

	for i, plan := range program.Metadata.AggregatePlans {
		if len(plan.Keys) != len(plan.Kinds) {
			return fmt.Errorf("%w: aggregate plan %d keys/kinds length mismatch", ErrInvalidProgram, i)
		}

		for j, kind := range plan.Kinds {
			if kind < AggregateCount || kind > AggregateAverage {
				return fmt.Errorf("%w: aggregate plan %d has invalid kind %d at index %d", ErrInvalidProgram, i, kind, j)
			}
		}
	}

	return nil
}

func validateCatchTable(program *Program) error {
	bytecodeLen := len(program.Bytecode)

	for i, entry := range program.CatchTable {
		start, end, handler := entry[0], entry[1], entry[2]

		if start < 0 || start >= bytecodeLen {
			return fmt.Errorf("%w: catch entry %d start pc %d out of range", ErrInvalidProgram, i, start)
		}

		if end < 0 || end >= bytecodeLen {
			return fmt.Errorf("%w: catch entry %d end pc %d out of range", ErrInvalidProgram, i, end)
		}

		if start > end {
			return fmt.Errorf("%w: catch entry %d start pc %d exceeds end pc %d", ErrInvalidProgram, i, start, end)
		}

		if handler < 0 || handler >= bytecodeLen {
			return fmt.Errorf("%w: catch entry %d handler pc %d out of range", ErrInvalidProgram, i, handler)
		}
	}

	return nil
}

func validateInstructions(program *Program) error {
	bytecodeLen := len(program.Bytecode)
	constantsLen := len(program.Constants)
	paramsLen := len(program.Params)
	registers := program.Registers

	for pc, inst := range program.Bytecode {
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

		if OpcodeInfoOf(op).Class == OpcodeClassUnknown {
			return fmt.Errorf("%w: unknown opcode %d at pc %d", ErrInvalidInstruction, op, pc)
		}

		switch op {
		case OpReturn:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}
		case OpJump:
			if err := validatePCOperand(dst, bytecodeLen, pc, "dst"); err != nil {
				return err
			}
		case OpJumpIfFalse, OpJumpIfTrue, OpJumpIfNone:
			if err := validatePCOperand(dst, bytecodeLen, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}
		case OpJumpIfNe, OpJumpIfEq, OpJumpIfMissingProperty:
			if err := validatePCOperand(dst, bytecodeLen, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src2, registers, pc, "src2"); err != nil {
				return err
			}
		case OpJumpIfNeConst, OpJumpIfEqConst:
			if err := validatePCOperand(dst, bytecodeLen, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateConstantOperand(src2, constantsLen, pc, "src2"); err != nil {
				return err
			}
		case OpJumpIfMissingPropertyConst:
			if err := validatePCOperand(dst, bytecodeLen, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateConstantOperand(src2, constantsLen, pc, "src2"); err != nil {
				return err
			}

			if _, ok := program.Constants[src2.Constant()].(runtime.String); !ok {
				return fmt.Errorf("%w: pc %d expects string constant in src2", ErrInvalidInstruction, pc)
			}
		case OpMove, OpMoveTracked:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}
		case OpLoadNone, OpLoadZero, OpIncr, OpDecr, OpClose, OpSleep, OpRand:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}
		case OpLoadBool, OpDataSet:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if src1 != 0 && src1 != 1 {
				return fmt.Errorf("%w: pc %d src1 must be 0 or 1", ErrInvalidInstruction, pc)
			}
		case OpLoadConst:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateConstantOperand(src1, constantsLen, pc, "src1"); err != nil {
				return err
			}
		case OpLoadParam:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			slot := int(src1)
			if slot <= 0 || slot > paramsLen {
				return fmt.Errorf("%w: pc %d parameter slot %d out of range", ErrInvalidInstruction, pc, slot)
			}
		case OpMakeCell, OpStoreCell, OpQuery:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateValueOperand(src1, registers, constantsLen, pc, "src1"); err != nil {
				return err
			}

			if op == OpQuery {
				if err := validateValueOperand(src2, registers, constantsLen, pc, "src2"); err != nil {
					return err
				}
			}
		case OpLoadCell:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}
		case OpLoadArray, OpLoadObject:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			size := int(src1)
			if size < 0 {
				return fmt.Errorf("%w: pc %d src1 must be non-negative", ErrInvalidInstruction, pc)
			}

			if size > MaxCollectionPreallocation {
				return fmt.Errorf("%w: pc %d src1 preallocation %d exceeds limit %d", ErrInvalidInstruction, pc, size, MaxCollectionPreallocation)
			}
		case OpLoadRange, OpLoadIndex, OpLoadIndexOptional, OpLoadKey, OpLoadKeyOptional,
			OpLoadProperty, OpLoadPropertyOptional, OpPushKV, OpObjectSet,
			OpDispatch, OpStream, OpStreamIter, OpAdd, OpSub,
			OpMul, OpDiv, OpMod, OpCmp, OpEq, OpNe, OpGt, OpLt, OpGte,
			OpLte, OpIn, OpLike, OpRegexp, OpAllEq, OpAllNe, OpAllGt,
			OpAllGte, OpAllLt, OpAllLte, OpAllIn, OpAnyEq, OpAnyNe,
			OpAnyGt, OpAnyGte, OpAnyLt, OpAnyLte, OpAnyIn, OpNoneEq,
			OpNoneNe, OpNoneGt, OpNoneGte, OpNoneLt, OpNoneLte, OpNoneIn,
			OpAggregateGroupUpdate:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src2, registers, pc, "src2"); err != nil {
				return err
			}
		case OpConcat:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateRegisterCountRange(src1, src2, registers, pc); err != nil {
				return err
			}
		case OpLoadAggregateKey:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateConstantOperand(src2, constantsLen, pc, "src2"); err != nil {
				return err
			}

			if _, ok := program.Constants[src2.Constant()].(runtime.Int); !ok {
				return fmt.Errorf("%w: pc %d expects int constant in src2", ErrInvalidInstruction, pc)
			}
		case OpLoadIndexConst, OpLoadIndexOptionalConst, OpLoadKeyConst, OpLoadKeyOptionalConst,
			OpLoadPropertyConst, OpLoadPropertyOptionalConst:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateConstantOperand(src2, constantsLen, pc, "src2"); err != nil {
				return err
			}
		case OpMatchLoadPropertyConst:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateConstantOperand(src2, constantsLen, pc, "src2"); err != nil {
				return err
			}

			if len(program.Metadata.MatchFailTargets) == 0 {
				return fmt.Errorf("%w: pc %d requires match fail target metadata", ErrInvalidInstruction, pc)
			}

			if program.Metadata.MatchFailTargets[pc] < 0 || program.Metadata.MatchFailTargets[pc] >= bytecodeLen {
				return fmt.Errorf("%w: pc %d has invalid match fail target %d", ErrInvalidInstruction, pc, program.Metadata.MatchFailTargets[pc])
			}

			if _, ok := program.Constants[src2.Constant()].(runtime.String); !ok {
				return fmt.Errorf("%w: pc %d expects string constant in src2", ErrInvalidInstruction, pc)
			}
		case OpPush, OpLength, OpType, OpNegate, OpFlipPositive, OpFlipNegative, OpCastBool, OpNot,
			OpExists, OpIter, OpIterValue, OpIterKey:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}
		case OpObjectSetConst:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateConstantOperand(src1, constantsLen, pc, "src1"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src2, registers, pc, "src2"); err != nil {
				return err
			}
		case OpAddConst:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateConstantOperand(src2, constantsLen, pc, "src2"); err != nil {
				return err
			}
		case OpCounterInc:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}
		case OpArrayPush, OpAggregateUpdate:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if op == OpAggregateUpdate {
				if len(program.Metadata.AggregateSelectorSlots) == 0 {
					return fmt.Errorf("%w: pc %d requires aggregate selector slot metadata", ErrInvalidInstruction, pc)
				}

				if program.Metadata.AggregateSelectorSlots[pc] < 0 {
					return fmt.Errorf("%w: pc %d has invalid aggregate selector slot %d", ErrInvalidInstruction, pc, program.Metadata.AggregateSelectorSlots[pc])
				}
			}
		case OpHCall, OpProtectedHCall, OpCall, OpProtectedCall, OpTailCall:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateCallRange(src1, src2, registers, pc); err != nil {
				return err
			}
		case OpDataSetCollector:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			collectorType := CollectorType(src1)
			if collectorType < CollectorTypeCounter || collectorType > CollectorTypeAggregateGroup {
				return fmt.Errorf("%w: pc %d has invalid collector type %d", ErrInvalidInstruction, pc, collectorType)
			}

			if collectorType == CollectorTypeAggregate || collectorType == CollectorTypeAggregateGroup {
				planIdx := int(src2)
				if planIdx < 0 || planIdx >= len(program.Metadata.AggregatePlans) {
					return fmt.Errorf("%w: pc %d has invalid aggregate plan index %d", ErrInvalidInstruction, pc, planIdx)
				}
			}
		case OpDataSetSorter:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if int(src1) < 0 {
				return fmt.Errorf("%w: pc %d sort direction must be non-negative", ErrInvalidInstruction, pc)
			}
		case OpDataSetMultiSorter:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if int(src1) < 0 || int(src2) < 0 {
				return fmt.Errorf("%w: pc %d sorter operands must be non-negative", ErrInvalidInstruction, pc)
			}

			if count := int(src2); count > MaxEncodedSortDirections {
				return fmt.Errorf("%w: pc %d sorter direction count %d exceeds limit %d", ErrInvalidInstruction, pc, count, MaxEncodedSortDirections)
			}
		case OpFlatten:
			if err := validateRegisterOperand(dst, registers, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if int(src2) < 0 {
				return fmt.Errorf("%w: pc %d flatten depth must be non-negative", ErrInvalidInstruction, pc)
			}
		case OpIterNext:
			if err := validatePCOperand(dst, bytecodeLen, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}
		case OpIterSkip, OpIterLimit:
			if err := validatePCOperand(dst, bytecodeLen, pc, "dst"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src1, registers, pc, "src1"); err != nil {
				return err
			}

			if err := validateRegisterOperand(src2, registers, pc, "src2"); err != nil {
				return err
			}
		case OpFail:
			if err := validateConstantOperand(dst, constantsLen, pc, "dst"); err != nil {
				return err
			}

			if _, ok := program.Constants[dst.Constant()].(runtime.String); !ok {
				return fmt.Errorf("%w: pc %d FAIL expects a string constant", ErrInvalidInstruction, pc)
			}
		default:
			return fmt.Errorf("%w: opcode %d at pc %d is not supported by validator", ErrInvalidInstruction, op, pc)
		}
	}

	return nil
}

func validateRegisterOperand(op Operand, registers, pc int, field string) error {
	if !op.IsRegister() {
		return fmt.Errorf("%w: pc %d %s must be a register operand", ErrInvalidInstruction, pc, field)
	}

	if op.Register() < 0 || op.Register() >= registers {
		return fmt.Errorf("%w: pc %d %s register %d out of range", ErrInvalidInstruction, pc, field, op.Register())
	}

	return nil
}

func validateConstantOperand(op Operand, constants, pc int, field string) error {
	if !op.IsConstant() {
		return fmt.Errorf("%w: pc %d %s must be a constant operand", ErrInvalidInstruction, pc, field)
	}

	idx := op.Constant()
	if idx < 0 || idx >= constants {
		return fmt.Errorf("%w: pc %d %s constant index %d out of range", ErrInvalidInstruction, pc, field, idx)
	}

	return nil
}

func validateValueOperand(op Operand, registers, constants, pc int, field string) error {
	if op.IsRegister() {
		return validateRegisterOperand(op, registers, pc, field)
	}

	return validateConstantOperand(op, constants, pc, field)
}

func validatePCOperand(op Operand, limit, pc int, field string) error {
	target := int(op)
	if target < 0 || target >= limit {
		return fmt.Errorf("%w: pc %d %s target %d out of range", ErrInvalidInstruction, pc, field, target)
	}

	return nil
}

func validateCallRange(start, end Operand, registers, pc int) error {
	if start == 0 && end == 0 {
		return nil
	}

	if !start.IsRegister() || !end.IsRegister() {
		return fmt.Errorf("%w: pc %d call arguments must be encoded as registers", ErrInvalidInstruction, pc)
	}

	if start.Register() <= 0 {
		return fmt.Errorf("%w: pc %d call argument range must start at register > 0", ErrInvalidInstruction, pc)
	}

	if end.Register() < start.Register() {
		return fmt.Errorf("%w: pc %d call argument range end %d precedes start %d", ErrInvalidInstruction, pc, end.Register(), start.Register())
	}

	if end.Register() >= registers {
		return fmt.Errorf("%w: pc %d call argument range end %d out of range", ErrInvalidInstruction, pc, end.Register())
	}

	return nil
}

func validateRegisterCountRange(start, count Operand, registers, pc int) error {
	if !start.IsRegister() {
		return fmt.Errorf("%w: pc %d concat source must start at a register", ErrInvalidInstruction, pc)
	}

	if count < 0 {
		return fmt.Errorf("%w: pc %d concat register count must be non-negative", ErrInvalidInstruction, pc)
	}

	startReg := start.Register()
	if startReg < 0 || startReg >= registers {
		return fmt.Errorf("%w: pc %d concat start register %d out of range", ErrInvalidInstruction, pc, startReg)
	}

	countVal := int(count)
	if countVal > registers-startReg {
		return fmt.Errorf("%w: pc %d concat register range [%d,%d) exceeds register file size %d", ErrInvalidInstruction, pc, startReg, startReg+countVal, registers)
	}

	return nil
}
