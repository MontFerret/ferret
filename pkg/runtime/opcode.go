package runtime

type Opcode byte

const (
	OpReturn      Opcode = iota
	OpMove               // Move a value from register A to register B
	OpLoadNone           // Set None value to a register
	OpLoadBool           // Set a boolean value to a register
	OpLoadConst          // Load a constant to a register or a global variable
	OpStoreGlobal        // Store a value from register A to a global variable
	OpLoadGlobal         // Load a global variable to a register A

	OpJump
	OpJumpIfFalse
	OpJumpIfTrue

	OpAdd
	OpSub
	OpMulti
	OpDiv
	OpMod
	OpIncr
	OpDecr

	OpCastBool
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
	OpRegexpPositive
	OpRegexpNegative

	OpArray  // Create an array from a list of registers (ARR R2, R3 R5 - creates an array in R2 with elements from R3 to R5)
	OpObject // Create an object from a list of registers (OBJ R2, R3 R5 - creates an object in R2 with elements from R3 to R5)
	OpRange

	OpLoadProperty
	OpLoadPropertyOptional

	OpLength

	OpCall
	OpProtectedCall

	OpLoopBegin // Creates a loop result dataset
	OpLoopPush
	OpLoopEnd

	OpForLoopInit // Creates an iterator for a loop
	OpForLoopNext
	OpForLoopValue
	OpForLoopKey

	OpWhileLoopInit
	OpWhileLoopNext
	OpWhileLoopValue
)
