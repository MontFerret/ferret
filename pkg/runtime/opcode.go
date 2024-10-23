package runtime

type OpCode byte

const (
	OpMove        OpCode = iota
	OpLoadConst          // Load a constant to a register A
	OpLoadNone           // Load None to a register A
	OpLoadTrue           // Load True to a register A
	OpLoadFalse          // Load False to a register A
	OpLoadGlobal         // Load a global variable to a register A
	OpStoreGlobal        // Store a value from register A to a global variable

	OpAdd
	OpSub
	OpMulti
	OpDiv
	OpMod
	OpIncr
	OpDecr

	OpArray
	OpObject
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

	OpRange
	OpRegexpPositive
	OpRegexpNegative
	OpCall
	OpCallSafe
	OpCall1
	OpCall1Safe
	OpCall2
	OpCall2Safe
	OpCall3
	OpCall3Safe
	OpCall4
	OpCall4Safe
	OpCallN
	OpCallNSafe
	OpJumpIfFalse
	OpJumpIfTrue
	OpJump
	OpJumpBackward
	OpLoopInitOutput
	OpLoopUnwrapOutput
	OpForLoopInitInput
	OpForLoopHasNext
	OpForLoopNext
	OpForLoopNextValue
	OpForLoopNextCounter
	OpWhileLoopInitCounter
	OpWhileLoopNext
	OpLoopReturn
	OpReturn

	OpCastBool
)
