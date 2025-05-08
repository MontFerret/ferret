package vm

type Opcode byte

const (
	OpReturn               Opcode = iota
	OpMove                        // Move a value from register A to register B
	OpLoadNone                    // Set None value to a register
	OpLoadBool                    // Set a boolean value to a register
	OpLoadZero                    // Set a zero value to a register
	OpLoadConst                   // Load a constant to a register or a global variable
	OpStoreGlobal                 // Store a value from register A to a global variable
	OpLoadGlobal                  // Load a global variable to a register A
	OpLoadParam                   // Load a parameter to a register A
	OpLoadProperty                // Load a property from an object to a register
	OpLoadPropertyOptional        // Load a property from an object to a register, if it exists

	OpJump
	OpJumpIfFalse
	OpJumpIfTrue
	OpJumpIfEmpty

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

	OpComp
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
	OpKeyValue // Create a key-value pair from a list of registers (KV R2, R3 R4 - creates a key-value pair in R2 with key from R3 and value from R4)

	OpSort     // Sorts a collection of KeyValue pairs. (SORT R2, R3 - sorts a collection in R2 with a sorting direction in R3)
	OpSortMany // Sorts a collection of KeyValue pairs with compound key and multiple directions. (SORT R2, R3, R4 - sorts a collection in R2 with a sorting direction from R3 to R4)
	OpLength
	OpType
	OpClose

	OpCall
	OpProtectedCall

	OpGroupPrep
	OpGroupAdd

	OpLoopBegin // Creates a loop result dataset
	OpLoopPush
	OpLoopLimit
	OpLoopSkip
	OpLoopEnd

	OpForLoopPrep // Creates an iterator for a loop
	OpForLoopNext
	OpForLoopValue
	OpForLoopKey

	OpWhileLoopPrep
	OpWhileLoopNext
	OpWhileLoopValue
	OpSleep
)
