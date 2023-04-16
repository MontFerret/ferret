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
	OpEqual
	OpNotEqual
	OpIn
	OpNotIn
	OpGreater
	OpLess
	OpGreaterOrEqual
	OpLessOrEqual
	OpLike
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
)
