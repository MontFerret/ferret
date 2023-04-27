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
	OpDefineGlobal
	OpGetGlobal
	OpGetLocal
	OpGetProperty
	OpGetPropertyOptional
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
	OpSafeCall
	OpCall1
	OpSafeCall1
	OpCall2
	OpSafeCall2
	OpCall3
	OpSafeCall3
	OpCall4
	OpSafeCall4
	OpCallN
	OpSafeCallN
	OpPop
	OpPush
	OpJumpIfFalse
	OpJumpIfTrue
	OpJump
	OpReturn
)
