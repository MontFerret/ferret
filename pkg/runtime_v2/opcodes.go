package runtime_v2

type Opcode byte

const (
	OpNone Opcode = iota
	OpConstant
	OpCastBool
	OpTrue
	OpFalse
	OpPop
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
	OpJumpIfFalse
	OpJumpIfTrue
	OpJump
	OpReturn
)
