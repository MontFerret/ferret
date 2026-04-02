package bytecode

type (
	OpcodeClass     uint8
	ControlFlowRole uint8
	CallKind        uint8
	CallArgEncoding uint8
)

const (
	OpcodeClassUnknown OpcodeClass = iota
	OpcodeClassControl
	OpcodeClassLoad
	OpcodeClassAccess
	OpcodeClassArithmetic
	OpcodeClassType
	OpcodeClassComparison
	OpcodeClassArray
	OpcodeClassUtility
	OpcodeClassCall
	OpcodeClassDataset
	OpcodeClassStream
	OpcodeClassIterator
)

const (
	ControlFlowNone ControlFlowRole = iota
	ControlFlowJumpUnconditional
	ControlFlowJumpConditional
	ControlFlowTerminator
)

const (
	CallKindNone CallKind = iota
	CallKindHost
	CallKindUser
	CallKindTail
)

const (
	CallArgEncodingNone CallArgEncoding = iota
	CallArgEncodingRegisterRange
)

type OpcodeInfo struct {
	Class           OpcodeClass
	ControlFlow     ControlFlowRole
	CallKind        CallKind
	CallArgEncoding CallArgEncoding
	ProtectedCall   bool
}

func OpcodeInfoOf(op Opcode) OpcodeInfo {
	return OpcodeInfo{
		Class:           opcodeClass(op),
		ControlFlow:     controlFlowRole(op),
		CallKind:        callKind(op),
		CallArgEncoding: callArgEncoding(op),
		ProtectedCall:   isProtectedCallOpcode(op),
	}
}

func IsTerminatorOpcode(op Opcode) bool {
	return controlFlowRole(op) == ControlFlowTerminator
}

func IsConditionalJumpOpcode(op Opcode) bool {
	return controlFlowRole(op) == ControlFlowJumpConditional
}

func IsUnconditionalJumpOpcode(op Opcode) bool {
	return controlFlowRole(op) == ControlFlowJumpUnconditional
}

func IsCallOpcode(op Opcode) bool {
	return callKind(op) != CallKindNone
}

func IsProtectedCallOpcode(op Opcode) bool {
	return isProtectedCallOpcode(op)
}

func JumpTargetOperandIndex(op Opcode) int {
	if op == OpMatchLoadPropertyConst {
		return -1
	}

	role := controlFlowRole(op)
	if role == ControlFlowJumpConditional || role == ControlFlowJumpUnconditional {
		return 0
	}

	return -1
}

func VisitCallArgumentRegisters(op Opcode, src1, src2 Operand, visit func(reg int)) {
	if callArgEncoding(op) != CallArgEncodingRegisterRange {
		return
	}

	if !src1.IsRegister() || !src2.IsRegister() {
		return
	}

	start := src1.Register()
	end := src2.Register()

	if start <= 0 || end < start {
		return
	}

	for r := start; r <= end; r++ {
		visit(r)
	}
}

func opcodeClass(op Opcode) OpcodeClass {
	switch op {
	case OpReturn, OpJump, OpJumpIfFalse, OpJumpIfTrue, OpJumpIfNone,
		OpJumpIfNe, OpJumpIfNeConst, OpJumpIfEq, OpJumpIfEqConst,
		OpJumpIfMissingProperty, OpJumpIfMissingPropertyConst:
		return OpcodeClassControl
	case OpMove, OpMoveTracked,
		OpLoadNone, OpLoadBool, OpLoadZero, OpLoadConst, OpLoadParam, OpMakeCell, OpLoadCell, OpLoadArray, OpLoadObject, OpLoadRange, OpLoadAggregateKey:
		return OpcodeClassLoad
	case OpLoadIndex, OpLoadIndexOptional, OpLoadKey, OpLoadKeyOptional,
		OpLoadProperty, OpLoadPropertyOptional, OpLoadIndexConst, OpLoadIndexOptionalConst,
		OpLoadKeyConst, OpLoadKeyOptionalConst, OpLoadPropertyConst, OpLoadPropertyOptionalConst, OpMatchLoadPropertyConst, OpQuery:
		return OpcodeClassAccess
	case OpAdd, OpAddConst, OpConcat, OpSub, OpMul, OpDiv, OpMod, OpIncr, OpDecr:
		return OpcodeClassArithmetic
	case OpCastBool, OpNegate, OpFlipPositive, OpFlipNegative, OpType:
		return OpcodeClassType
	case OpCmp, OpNot, OpEq, OpNe, OpGt, OpLt, OpGte, OpLte, OpIn, OpLike, OpRegexp:
		return OpcodeClassComparison
	case OpFlatten, OpAnyEq, OpAnyNe, OpAnyGt, OpAnyGte, OpAnyLt, OpAnyLte, OpAnyIn,
		OpNoneEq, OpNoneNe, OpNoneGt, OpNoneGte, OpNoneLt, OpNoneLte, OpNoneIn,
		OpAllEq, OpAllNe, OpAllGt, OpAllGte, OpAllLt, OpAllLte, OpAllIn:
		return OpcodeClassArray
	case OpLength, OpClose, OpSleep, OpExists, OpRand, OpDispatch, OpFail, OpFailTimeout, OpStoreCell:
		return OpcodeClassUtility
	case OpHCall, OpProtectedHCall, OpCall, OpProtectedCall, OpTailCall:
		return OpcodeClassCall
	case OpDataSet, OpDataSetCollector, OpDataSetSorter, OpDataSetMultiSorter,
		OpPush, OpPushKV, OpCounterInc, OpArrayPush, OpObjectSet, OpObjectSetConst,
		OpAggregateUpdate, OpAggregateGroupUpdate:
		return OpcodeClassDataset
	case OpStream, OpStreamIter:
		return OpcodeClassStream
	case OpIter, OpIterNext, OpIterNextTimeout, OpIterValue, OpIterKey, OpIterLimit, OpIterSkip:
		return OpcodeClassIterator
	default:
		return OpcodeClassUnknown
	}
}

func controlFlowRole(op Opcode) ControlFlowRole {
	switch op {
	case OpJump:
		return ControlFlowJumpUnconditional
	case OpJumpIfFalse, OpJumpIfTrue, OpJumpIfNone,
		OpJumpIfNe, OpJumpIfNeConst, OpJumpIfEq, OpJumpIfEqConst,
		OpJumpIfMissingProperty, OpJumpIfMissingPropertyConst, OpMatchLoadPropertyConst,
		OpIterNext, OpIterNextTimeout:
		return ControlFlowJumpConditional
	case OpReturn, OpTailCall:
		return ControlFlowTerminator
	default:
		return ControlFlowNone
	}
}

func callKind(op Opcode) CallKind {
	switch op {
	case OpHCall, OpProtectedHCall:
		return CallKindHost
	case OpCall, OpProtectedCall:
		return CallKindUser
	case OpTailCall:
		return CallKindTail
	default:
		return CallKindNone
	}
}

func callArgEncoding(op Opcode) CallArgEncoding {
	switch op {
	case OpHCall, OpProtectedHCall, OpCall, OpProtectedCall, OpTailCall:
		return CallArgEncodingRegisterRange
	default:
		return CallArgEncodingNone
	}
}

func isProtectedCallOpcode(op Opcode) bool {
	switch op {
	case OpProtectedHCall, OpProtectedCall:
		return true
	default:
		return false
	}
}
