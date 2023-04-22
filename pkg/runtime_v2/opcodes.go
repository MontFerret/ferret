package runtime_v2

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
	OpPop
	OpPush
	OpJumpIfFalse
	OpJumpIfTrue
	OpJump
	OpReturn
)
