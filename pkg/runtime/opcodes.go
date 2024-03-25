package runtime

type Opcode byte

const (
	OpNone Opcode = iota
	OpConstant
	OpCastBool
	OpTrue
	OpFalse
	OpArray
	OpObject
	OpLoadGlobal
	OpStoreGlobal
	OpLoadLocal
	OpStoreLocal
	OpPopLocal
	OpLoadProperty
	OpLoadPropertyOptional
	OpNegate
	OpFlipPositive
	OpFlipNegative
	OpNot
	OpEq
	OpNeq
	OpIn
	OpNotIn
	OpGt
	OpLt
	OpGte
	OpLte
	OpLike
	OpNotLike
	OpAdd
	OpSub
	OpMulti
	OpDiv
	OpMod
	OpIncr
	OpDecr
	OpRange
	OpRegexpPositive
	OpRegexpNegative
	OpCall
	OpCall1
	OpCall2
	OpCall3
	OpCall4
	OpCallN
	OpPush
	OpPop
	OpPopClose
	OpJumpIfFalse
	OpJumpIfTrue
	OpJump
	OpJumpBackward
	OpLoopSourceInit
	OpLoopHasNext
	OpLoopNext
	OpLoopNextValue
	OpLoopNextCounter
	OpLoopReturn
	OpReturn
)
