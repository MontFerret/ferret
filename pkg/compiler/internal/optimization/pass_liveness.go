package optimization

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

const (
	LivenessAnalysisPassName = "liveness-analysis"
	LivenessAnalysis         = "liveness"
)

// LivenessAnalysisPass performs liveness analysis to determine which registers
// are live (potentially used) at each program point
type (
	LivenessAnalysisPass struct{}

	// LivenessInfo contains liveness information for a basic block
	LivenessInfo struct {
		LiveIn  map[int]bool // Registers live at block entry
		LiveOut map[int]bool // Registers live at block exit
		Use     map[int]bool // Registers used in block before definition
		Def     map[int]bool // Registers defined in block
	}
)

// NewLivenessAnalysisPass creates a new liveness analysis pass
func NewLivenessAnalysisPass() Pass {
	return &LivenessAnalysisPass{}
}

func (p *LivenessAnalysisPass) Requires() []string {
	return []string{}
}

// Name returns the pass name
func (p *LivenessAnalysisPass) Name() string {
	return LivenessAnalysisPassName
}

// Run executes liveness analysis on the program
func (p *LivenessAnalysisPass) Run(c *PassContext) (*PassResult, error) {
	liveness := computeLiveness(c.CFG)

	return &PassResult{
		Modified: false,
		Metadata: map[string]interface{}{
			"liveness": liveness,
		},
	}, nil
}

// computeLiveness performs dataflow analysis to compute liveness information
func computeLiveness(cfg *ControlFlowGraph) map[int]*LivenessInfo {
	info := make(map[int]*LivenessInfo)

	// Initialize liveness info for each block
	for _, block := range cfg.Blocks {
		info[block.ID] = &LivenessInfo{
			LiveIn:  make(map[int]bool),
			LiveOut: make(map[int]bool),
			Use:     make(map[int]bool),
			Def:     make(map[int]bool),
		}

		// Compute Use and Def sets for this block
		for _, inst := range block.Instructions {
			computeUseDefForInstruction(inst, info[block.ID])
		}
	}

	// Iteratively compute LiveIn and LiveOut until convergence
	changed := true
	for changed {
		changed = false

		// Process blocks in reverse order (better for backward dataflow)
		for i := len(cfg.Blocks) - 1; i >= 0; i-- {
			block := cfg.Blocks[i]
			blockInfo := info[block.ID]

			// LiveOut[B] = Union of LiveIn[S] for all successors S of B
			newLiveOut := make(map[int]bool)

			for _, succ := range block.Successors {
				succInfo := info[succ.ID]

				for reg := range succInfo.LiveIn {
					newLiveOut[reg] = true
				}
			}

			// LiveIn[B] = Use[B] Union (LiveOut[B] - Def[B])
			newLiveIn := make(map[int]bool)

			for reg := range blockInfo.Use {
				newLiveIn[reg] = true
			}

			for reg := range newLiveOut {
				if !blockInfo.Def[reg] {
					newLiveIn[reg] = true
				}
			}

			// Check if anything changed
			if !mapsEqual(blockInfo.LiveIn, newLiveIn) || !mapsEqual(blockInfo.LiveOut, newLiveOut) {
				changed = true
				blockInfo.LiveIn = newLiveIn
				blockInfo.LiveOut = newLiveOut
			}
		}
	}

	return info
}

// instructionUseDef returns registers used and defined by an instruction.
// Only registers (non-constants) are included.
func instructionUseDef(inst bytecode.Instruction) (uses []int, defs []int) {
	addUse := func(op bytecode.Operand) {
		if op != bytecode.NoopOperand && op.IsRegister() {
			uses = append(uses, op.Register())
		}
	}
	addDef := func(op bytecode.Operand) {
		if op != bytecode.NoopOperand && op.IsRegister() {
			defs = append(defs, op.Register())
		}
	}
	op := inst.Opcode
	dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

	switch op {
	// No-operand terminator.
	case bytecode.OpReturn:
		addUse(dst)
		return
	// Moves / loads.
	case bytecode.OpMove:
		addUse(src1)
		addDef(dst)
		return
	case bytecode.OpLoadConst, bytecode.OpLoadParam, bytecode.OpLoadNone, bytecode.OpLoadBool, bytecode.OpLoadZero, bytecode.OpRand:
		addDef(dst)
		return
	case bytecode.OpConcat:
		addDef(dst)

		if !src1.IsRegister() {
			return
		}

		startReg := src1.Register()
		count := int(src2)

		if count <= 0 || startReg <= 0 {
			return
		}

		for r := startReg; r < startReg+count; r++ {
			uses = append(uses, r)
		}

		return
	case bytecode.OpLoadArray:
		addDef(dst)
		return
	case bytecode.OpLoadObject:
		addDef(dst)
		return
	case bytecode.OpLoadRange:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Simple arithmetic, comparisons, access.
	case bytecode.OpAdd, bytecode.OpAddConst, bytecode.OpSub, bytecode.OpMulti, bytecode.OpDiv, bytecode.OpMod,
		bytecode.OpCmp,
		bytecode.OpEq, bytecode.OpNe, bytecode.OpGt, bytecode.OpLt, bytecode.OpGte, bytecode.OpLte,
		bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte,
		bytecode.OpAnyIn,
		bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte,
		bytecode.OpNoneIn,
		bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte,
		bytecode.OpAllIn,
		bytecode.OpIn, bytecode.OpLike, bytecode.OpRegexp,
		bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional, bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst,
		bytecode.OpLoadKey, bytecode.OpLoadKeyOptional, bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst,
		bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional, bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst,
		bytecode.OpApplyQuery:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Unary ops.
	case bytecode.OpIncr, bytecode.OpDecr:
		addUse(dst)
		addDef(dst)
		return
	case bytecode.OpCastBool, bytecode.OpNegate, bytecode.OpNot, bytecode.OpFlipPositive, bytecode.OpFlipNegative, bytecode.OpLength, bytecode.OpType, bytecode.OpFlatten, bytecode.OpExists:
		addUse(src1)
		addDef(dst)
		return

	// Control flow.
	case bytecode.OpJumpIfFalse, bytecode.OpJumpIfTrue, bytecode.OpJumpIfNone:
		addUse(src1)
		return
	case bytecode.OpJumpIfNe:
		addUse(src1)
		addUse(src2)
		return
	case bytecode.OpJumpIfNeConst:
		addUse(src1)
		return
	case bytecode.OpJumpIfEq:
		addUse(src1)
		addUse(src2)
		return
	case bytecode.OpJumpIfEqConst:
		addUse(src1)
		return
	case bytecode.OpJumpIfMissingProperty:
		addUse(src1)
		addUse(src2)
		return
	case bytecode.OpJumpIfMissingPropertyConst:
		addUse(src1)
		return
	case bytecode.OpJump, bytecode.OpFail:
		return

	// Dataset operations.
	case bytecode.OpDataSet, bytecode.OpDataSetCollector, bytecode.OpDataSetSorter, bytecode.OpDataSetMultiSorter:
		addDef(dst)
		return
	case bytecode.OpPush, bytecode.OpArrayPush:
		addUse(dst)
		addUse(src1)
		return
	case bytecode.OpPushKV, bytecode.OpObjectSet, bytecode.OpObjectSetConst:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		return

	// Iterators.
	case bytecode.OpIter:
		addUse(src1)
		addDef(dst)
		return
	case bytecode.OpIterValue, bytecode.OpIterKey:
		addUse(src1)
		addDef(dst)
		return
	case bytecode.OpIterLimit, bytecode.OpIterSkip:
		addUse(src1)
		addUse(src2)
		addDef(src1)
		return
	case bytecode.OpIterNext:
		addUse(src1)
		return

	// Host calls.
	case bytecode.OpHCall, bytecode.OpProtectedHCall, bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
		addUse(dst)
		bytecode.VisitCallArgumentRegisters(op, src1, src2, func(reg int) {
			uses = append(uses, reg)
		})
		if op != bytecode.OpTailCall {
			addDef(dst)
		}
		return

	// Stream.
	case bytecode.OpStream:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case bytecode.OpStreamIter:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case bytecode.OpDispatch:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Utility.
	case bytecode.OpClose:
		addUse(dst)
		addDef(dst)
		return
	case bytecode.OpSleep:
		addUse(dst)
		return
	}

	return
}

// computeUseDefForInstruction updates Use and Def sets for an instruction
func computeUseDefForInstruction(inst bytecode.Instruction, info *LivenessInfo) {
	uses, defs := instructionUseDef(inst)

	// Use-before-def within the block
	for _, reg := range uses {
		if !info.Def[reg] {
			info.Use[reg] = true
		}
	}

	for _, reg := range defs {
		info.Def[reg] = true
	}
}

// mapsEqual checks if two register maps are equal
func mapsEqual(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if !b[k] {
			return false
		}
	}

	return true
}
