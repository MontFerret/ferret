package vm

type Opcode byte

const (
	// Control Flow
	OpReturn Opcode = iota
	OpJump
	OpJumpIfFalse
	OpJumpIfTrue
	OpJumpIfEmpty

	// Register Operations
	OpMove // Move a value from register A to register B

	// Loading Operations
	OpLoadNone  // Set None value to a register
	OpLoadBool  // Set a boolean value to a register
	OpLoadZero  // Set a zero value to a register
	OpLoadConst // Load a constant to a register or a global variable
	OpLoadParam // Load a parameter to a register A

	// Global Variable Operations
	OpLoadGlobal  // Load a global variable to register A
	OpStoreGlobal // Store a value from register A to a global variable

	// Collection Access Operations
	OpLoadIndex            // Load a value from a list to a register
	OpLoadIndexOptional    // Load a value from a list to a register, if it exists
	OpLoadKey              // Load a value from a map to a register
	OpLoadKeyOptional      // Load a value from a map to a register, if it exists
	OpLoadProperty         // Load a property from an object to a register
	OpLoadPropertyOptional // Load a property from an object to a register, if it exists

	// Arithmetic Operations
	OpAdd
	OpSub
	OpMulti
	OpDiv
	OpMod
	OpIncr
	OpDecr

	// Type Operations
	OpCastBool
	OpNegate
	OpFlipPositive
	OpFlipNegative

	// Comparison Operations
	OpComp
	OpNot
	OpEq
	OpNeq
	OpGt
	OpLt
	OpGte
	OpLte

	// Membership & Pattern Matching
	OpIn
	OpNotIn
	OpLike
	OpNotLike
	OpRegexpPositive
	OpRegexpNegative

	// Utility Operations
	OpLength
	OpType
	OpClose
	OpSleep

	// Function Operations
	OpCall
	OpProtectedCall

	// Collection Creation
	OpList  // Create an array
	OpMap   // Create an object
	OpRange // Create a range

	// Dataset Operations
	OpDataSet
	OpDataSetCollector
	OpDataSetSorter
	OpDataSetMultiSorter
	OpPush   // Adds a value to a dataset
	OpPushKV // Adds a key-value pair to a dataset

	// Stream Operations
	OpStream     // Subscribes to a stream
	OpStreamIter // Consumes a stream

	// Iterator Operations
	OpIter      // Creates an iterator
	OpIterNext  // Moves to the next element
	OpIterValue // Returns the current value
	OpIterKey   // Returns the current key
	OpIterLimit
	OpIterSkip
)
