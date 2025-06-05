package vm

type Opcode byte

const (
	OpReturn               Opcode = iota
	OpMove                        // Move a value from register A to register B
	OpLoadNone                    // Set None value to a register
	OpLoadBool                    // Set a boolean value to a register
	OpLoadZero                    // Set a zero value to a register
	OpLoadConst                   // Load a constant to a register or a global variable
	OpLoadGlobal                  // Load a global variable to a register A
	OpStoreGlobal                 // Store a value from register A to a global variable
	OpLoadParam                   // Load a parameter to a register A
	OpLoadList                    // Load an array from a list of registers (ARR R2, R3 R5 - creates an array in R2 with elements from R3 to R5)
	OpLoadMap                     // Load an object from a list of registers (OBJ R2, R3 R5 - creates an object in R2 with elements from R3 to R5)
	OpLoadRange                   // Load a range from a list of registers (RNG R2, R3, R4 - creates a range in R2 with start from R3 and end at R4)
	OpLoadIndex                   // Load a value from a list to a register (INDEX R1, R2, R3 - loads a value from a list in R2 to R1)
	OpLoadIndexOptional           // Load a value from a list to a register, if it exists
	OpLoadKey                     // Load a value from a map to a register (KEY R1, R2, R3 - loads a value from a map in R2 to R1)
	OpLoadKeyOptional             // Load a value from a map to a register, if it exists
	OpLoadProperty                // Load a property (key or index) from an object (map or list) to a register
	OpLoadPropertyOptional        // Load a property (key or index) from an object (map or list) to a register, if it exists
	OpLoadDataSet                 // Load a dataset to a register A

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

	OpLength
	OpType
	OpClose
	OpSleep

	OpCall
	OpProtectedCall

	OpStream     // Subscribes to a stream (SMRCV R2, R3, R4 - subscribes to a stream in R2 with a collection from R3 and optional params from R4)
	OpStreamIter // Consumes a stream (SMRD R2, R3 - consumes a stream in R2 with a collection from R3)

	OpIter      // Creates an iterator (ITER R2, R3 [, R4] - creates an iterator in R2 with a collection from R3 and optional params from R4)
	OpIterNext  // Moves to the next element in the iterator (ITER R2, R3  - moves to the next element in the iterator in R2 with a collection from R3)
	OpIterValue // Returns the current value from the iterator (ITER R2, R3  - returns the current value from the iterator in R2 with a collection from R3)
	OpIterKey   // Returns the current key from the iterator (ITER R2, R3  - returns the current key from the iterator in R2 with a collection from R3)

	OpWhileLoopPrep
	OpWhileLoopNext
	OpWhileLoopValue

	OpPush      // Adds a value to a dataset
	OpPushKV    // Adds a key-value pair to a dataset
	OpCollectK  // Adds a key to a group
	OpCollectKV // Adds a value to a group using key
	OpLimit
	OpSkip

	OpSort     // Sorts a collection of KeyValue pairs. (SORT R2, R3 - sorts a collection in R2 with a sorting direction in R3)
	OpSortMany // Sorts a collection of KeyValue pairs with compound key and multiple directions. (SORT R2, R3, R4 - sorts a collection in R2 with a sorting direction from R3 to R4)
)
