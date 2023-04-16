package runtime_v2

type Opcode byte

const (
	OpNone Opcode = iota
	OpConstant
	OpTrue
	OpFalse
	OpPop
	OpDefineGlobal
	OpGetGlobal
	OpReturn
	OpNegate
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
	OpIncrement
	OpDecrement
	OpRange
	OpRegexpPositive
	OpRegexpNegative
	OpJumpIfFalse
	OpJump
)
