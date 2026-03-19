package bytecode

type Opcode byte

const (
	// Control Flow
	OpReturn Opcode = iota
	OpJump
	OpJumpIfFalse
	OpJumpIfTrue
	OpJumpIfNone

	// Register Operations
	OpMove // Plain register copy — no tracking

	// Tracked Register Operations
	OpMoveTracked // Tracked register copy — ownership-aware

	// Loading Operations
	OpLoadNone   // Set None value to a register
	OpLoadBool   // Set a boolean value to a register
	OpLoadZero   // Set a zero value to a register
	OpLoadConst  // Load a constant to a register or a global variable
	OpLoadParam  // Load a parameter slot to a register A
	OpLoadArray  // Create an array
	OpLoadObject // Create an object
	OpLoadRange  // Create a range

	// Collection Access Operations
	OpLoadIndex            // Load a value from a built-in Array to a register
	OpLoadIndexOptional    // Load a value from a built-in Array to a register, if it exists
	OpLoadKey              // Load a value from a built-in Object to a register
	OpLoadKeyOptional      // Load a value from a built-in Object to a register, if it exists
	OpLoadProperty         // Load a property from a map or list to a register
	OpLoadPropertyOptional // Load a property from a map or list to a register, if it exists
	OpLoadIndexConst       // Load a value from a built-in Array to a register using a constant index
	OpLoadIndexOptionalConst
	OpLoadKeyConst // Load a value from a built-in Object to a register using a constant key
	OpLoadKeyOptionalConst
	OpLoadPropertyConst // Load a property from a map or list to a register using a constant key
	OpLoadPropertyOptionalConst
	OpMatchLoadPropertyConst // Load an object-pattern property or jump to the arm fail target

	// Integrated Query Operations
	OpQuery // Apply a query to a value

	// Arithmetic Operations
	OpAdd
	OpAddConst
	OpConcat
	OpSub
	OpMul
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
	OpCmp
	OpNot
	OpEq
	OpNe
	OpGt
	OpLt
	OpGte
	OpLte

	// Membership & Pattern Matching
	OpIn
	OpLike
	OpRegexp

	// Array Operations
	OpFlatten // Flattens nested arrays

	// Array Comparison Operations (do not swap order of EQ, GT, GTE, LT, LTE, IN. it must be equal to operators)
	OpAnyEq
	OpAnyNe
	OpAnyGt
	OpAnyGte
	OpAnyLt
	OpAnyLte
	OpAnyIn
	OpNoneEq
	OpNoneNe
	OpNoneGt
	OpNoneGte
	OpNoneLt
	OpNoneLte
	OpNoneIn
	OpAllEq
	OpAllNe
	OpAllGt
	OpAllGte
	OpAllLt
	OpAllLte
	OpAllIn

	// Utility Operations
	OpLength
	OpType
	OpClose
	OpSleep

	// Host Function Operations
	OpHCall
	OpProtectedHCall

	// UDF Operations
	OpCall
	OpProtectedCall
	OpTailCall

	// Dataset Operations
	OpDataSet
	OpDataSetCollector
	OpDataSetSorter
	OpDataSetMultiSorter
	OpPush           // Adds a value to a generic List
	OpPushKV         // Adds a key-value pair to a dataset
	OpCounterInc     // Increments a counter collector
	OpArrayPush      // Adds a value to a built-in Array instance
	OpObjectSet      // Sets a property on a built-in Object instance
	OpObjectSetConst // Sets a property on a built-in Object instance using a constant key

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

	// Existence Operations
	OpExists

	// Random Operations
	OpRand

	// Dispatch Operations
	OpDispatch

	// Compare-Jump Operations
	OpJumpIfNe
	OpJumpIfNeConst
	OpJumpIfEq
	OpJumpIfEqConst
	OpJumpIfMissingProperty
	OpJumpIfMissingPropertyConst
	OpFail // Raises a runtime failure with a constant message

	// Internal Aggregate Operations
	OpLoadAggregateKey // Creates an internal grouped-aggregate selector key
	OpAggregateUpdate
	OpAggregateGroupUpdate
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
	case OpJumpIfNone:
		return "JMPN"
	case OpJumpIfNe:
		return "JMPNE"
	case OpJumpIfNeConst:
		return "JMPNEC"
	case OpJumpIfEq:
		return "JMPEQ"
	case OpJumpIfEqConst:
		return "JMPEQC"
	case OpJumpIfMissingProperty:
		return "JMPMISSPROP"
	case OpJumpIfMissingPropertyConst:
		return "JMPMISSPROPC"
	case OpFail:
		return "FAIL"
	case OpLoadAggregateKey:
		return "LOADAGGK"

	// Register Operations
	case OpMove:
		return "MOVE"
	case OpMoveTracked:
		return "MOVET"

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

	// Collection Access Operations
	case OpLoadIndex:
		return "LOADI"
	case OpLoadIndexOptional:
		return "LOADIO"
	case OpLoadIndexConst:
		return "LOADIC"
	case OpLoadIndexOptionalConst:
		return "LOADIOC"
	case OpLoadKey:
		return "LOADK"
	case OpLoadKeyOptional:
		return "LOADKO"
	case OpLoadKeyConst:
		return "LOADKC"
	case OpLoadKeyOptionalConst:
		return "LOADKOC"
	case OpLoadPropertyConst:
		return "LOADPRC"
	case OpLoadPropertyOptionalConst:
		return "LOADPROC"
	case OpMatchLoadPropertyConst:
		return "MATCHPRC"
	case OpQuery:
		return "QRY"
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
	case OpAddConst:
		return "ADDC"
	case OpConcat:
		return "CONCAT"
	case OpSub:
		return "SUB"
	case OpMul:
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
		return "FPP"
	case OpFlipNegative:
		return "FPN"

	// Comparison Operations
	case OpCmp:
		return "CMP"
	case OpNot:
		return "NOT"
	case OpEq:
		return "EQ"
	case OpNe:
		return "NE"
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
	case OpLike:
		return "LIKE"
	case OpRegexp:
		return "REGEX"

	// Array Comparison Operations
	case OpAnyEq:
		return "ANYEQ"
	case OpAnyNe:
		return "ANYNE"
	case OpAnyGt:
		return "ANYGT"
	case OpAnyGte:
		return "ANYGE"
	case OpAnyLt:
		return "ANYLT"
	case OpAnyLte:
		return "ANYLE"
	case OpAnyIn:
		return "ANYIN"
	case OpNoneEq:
		return "NONEQ"
	case OpNoneNe:
		return "NONNE"
	case OpNoneGt:
		return "NONGT"
	case OpNoneGte:
		return "NONGE"
	case OpNoneLt:
		return "NONLT"
	case OpNoneLte:
		return "NONLE"
	case OpNoneIn:
		return "NONIN"
	case OpAllEq:
		return "ALLEQ"
	case OpAllNe:
		return "ALLNE"
	case OpAllGt:
		return "ALLGT"
	case OpAllGte:
		return "ALLGE"
	case OpAllLt:
		return "ALLLT"
	case OpAllLte:
		return "ALLLE"
	case OpAllIn:
		return "ALLIN"

	// Utility Operations
	case OpLength:
		return "LEN"
	case OpType:
		return "TYPE"
	case OpClose:
		return "CLOSE"
	case OpSleep:
		return "SLEEP"
	case OpExists:
		return "EXISTS"
	case OpRand:
		return "RAND"
	case OpDispatch:
		return "DISP"

	// Host Function Operations
	case OpHCall:
		return "HCALL"
	case OpProtectedHCall:
		return "PHCALL"

	// UDF Operations
	case OpCall:
		return "CALL"
	case OpProtectedCall:
		return "PCALL"
	case OpTailCall:
		return "TAILCALL"

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
	case OpArrayPush:
		return "PUSHA"
	case OpPushKV:
		return "PUSHKV"
	case OpCounterInc:
		return "CNTINC"
	case OpObjectSet:
		return "OBJSET"
	case OpObjectSetConst:
		return "OBJSETC"
	case OpAggregateUpdate:
		return "AGGUPD"
	case OpAggregateGroupUpdate:
		return "AGGGRPUPD"

	// Stream Operations
	case OpStream:
		return "STRM"
	case OpStreamIter:
		return "STRMI"

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
		return "ITLIM"
	case OpIterSkip:
		return "ITSKP"
	case OpFlatten:
		return "FLATTEN"

	default:
		return "UNKNOWN"
	}
}
