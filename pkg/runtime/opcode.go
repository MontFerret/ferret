package runtime

type Opcode byte

const (
	OpMove        Opcode = iota
	OpLoadConst          // Load a constant to a register or a global variable
	OpStoreGlobal        // Store a value from register A to a global variable
	OpLoadGlobal         // Load a global variable to a register A

	OpJump
	OpJumpIfFalse
	OpJumpIfTrue
	OpJumpBackward

	OpAdd
	OpSub
	OpMulti
	OpDiv
	OpMod
	OpIncr
	OpDecr

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

	OpArray  // Create an array from a list of registers (ARR R2, R3 R5 - creates an array in R2 with elements from R3 to R5)
	OpObject // Create an object from a list of registers (OBJ R2, R3 R5 - creates an object in R2 with elements from R3 to R5)
	OpRange

	OpLoadProperty
	OpLoadPropertyOptional

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
