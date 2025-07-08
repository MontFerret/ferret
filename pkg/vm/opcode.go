package vm

type Opcode byte

const (
	// Control Flow
	OpReturn Opcode = iota
	OpJump
	OpJumpIfFalse
	OpJumpIfTrue

	// Register Operations
	OpMove // Move a value from register A to register B

	// Loading Operations
	OpLoadNone   // Set None value to a register
	OpLoadBool   // Set a boolean value to a register
	OpLoadZero   // Set a zero value to a register
	OpLoadConst  // Load a constant to a register or a global variable
	OpLoadParam  // Load a parameter to a register A
	OpLoadArray  // Create an array
	OpLoadObject // Create an object
	OpLoadRange  // Create a range

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
	OpCall0
	OpProtectedCall0
	OpCall1
	OpProtectedCall1
	OpCall2
	OpProtectedCall2
	OpCall3
	OpProtectedCall3
	OpCall4
	OpProtectedCall4

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

func (op Opcode) String() string {
	switch op {
	// Control Flow
	case OpReturn:
		return "RET"
	case OpJump:
		return "JMP"
	case OpJumpIfFalse:
		return "JMPF"
	case OpJumpIfTrue:
		return "JMPT"

	// Register Operations
	case OpMove:
		return "MOVE"

	// Loading Operations
	case OpLoadNone:
		return "LOADN"
	case OpLoadBool:
		return "LOADB"
	case OpLoadZero:
		return "LOADZ"
	case OpLoadConst:
		return "LOADC"
	case OpLoadParam:
		return "LOADP"

	// Global Variable Operations
	case OpLoadGlobal:
		return "LOADG"
	case OpStoreGlobal:
		return "STOREG"

	// Collection Access Operations
	case OpLoadIndex:
		return "LOADI"
	case OpLoadIndexOptional:
		return "LOADIO"
	case OpLoadKey:
		return "LOADK"
	case OpLoadKeyOptional:
		return "LOADKO"
	case OpLoadProperty:
		return "LOADPR"
	case OpLoadPropertyOptional:
		return "LOADPRO"

		// Collection Creation
	case OpLoadArray:
		return "LOADARR"
	case OpLoadObject:
		return "LOADOBJ"
	case OpLoadRange:
		return "LOADRANGE"

	// Arithmetic Operations
	case OpAdd:
		return "ADD"
	case OpSub:
		return "SUB"
	case OpMulti:
		return "MUL"
	case OpDiv:
		return "DIV"
	case OpMod:
		return "MOD"
	case OpIncr:
		return "INCR"
	case OpDecr:
		return "DECR"

	// Type Operations
	case OpCastBool:
		return "CASTB"
	case OpNegate:
		return "NEG"
	case OpFlipPositive:
		return "FLIP+"
	case OpFlipNegative:
		return "FLIP-"

	// Comparison Operations
	case OpComp:
		return "COMP"
	case OpNot:
		return "NOT"
	case OpEq:
		return "EQ"
	case OpNeq:
		return "NEQ"
	case OpGt:
		return "GT"
	case OpLt:
		return "LT"
	case OpGte:
		return "GTE"
	case OpLte:
		return "LTE"

	// Membership & Pattern Matching
	case OpIn:
		return "IN"
	case OpNotIn:
		return "NOTIN"
	case OpLike:
		return "LIKE"
	case OpNotLike:
		return "NOTLIKE"
	case OpRegexpPositive:
		return "REGEX+"
	case OpRegexpNegative:
		return "REGEX-"

	// Utility Operations
	case OpLength:
		return "LEN"
	case OpType:
		return "TYPE"
	case OpClose:
		return "CLOSE"
	case OpSleep:
		return "SLEEP"

	// Function Operations
	case OpCall:
		return "CALL"
	case OpProtectedCall:
		return "PCALL"
	case OpCall0:
		return "CALL0"
	case OpProtectedCall0:
		return "PCALL0"
	case OpCall1:
		return "CALL1"
	case OpProtectedCall1:
		return "PCALL1"
	case OpCall2:
		return "CALL2"
	case OpProtectedCall2:
		return "PCALL2"
	case OpCall3:
		return "CALL3"
	case OpProtectedCall3:
		return "PCALL3"
	case OpCall4:
		return "CALL4"
	case OpProtectedCall4:
		return "PCALL4"

	// Dataset Operations
	case OpDataSet:
		return "DSET"
	case OpDataSetCollector:
		return "DSETC"
	case OpDataSetSorter:
		return "DSETS"
	case OpDataSetMultiSorter:
		return "DSETMS"
	case OpPush:
		return "PUSH"
	case OpPushKV:
		return "PUSHKV"

	// Stream Operations
	case OpStream:
		return "STRM"
	case OpStreamIter:
		return "STRMITER"

	// Iterator Operations
	case OpIter:
		return "ITER"
	case OpIterNext:
		return "ITNEXT"
	case OpIterValue:
		return "ITVAL"
	case OpIterKey:
		return "ITKEY"
	case OpIterLimit:
		return "ITLIMIT"
	case OpIterSkip:
		return "ITSKIP"

	default:
		return "UNKNOWN"
	}
}
