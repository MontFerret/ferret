package optimization

import (
	"github.com/MontFerret/ferret/pkg/vm"
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
func instructionUseDef(inst vm.Instruction) (uses []int, defs []int) {
	addUse := func(op vm.Operand) {
		if op != vm.NoopOperand && op.IsRegister() {
			uses = append(uses, op.Register())
		}
	}
	addDef := func(op vm.Operand) {
		if op != vm.NoopOperand && op.IsRegister() {
			defs = append(defs, op.Register())
		}
	}
	addRangeUses := func(start, end vm.Operand) {
		if !start.IsRegister() || !end.IsRegister() {
			return
		}
		startReg := start.Register()
		endReg := end.Register()
		if startReg <= 0 || endReg < startReg {
			return
		}
		for r := startReg; r <= endReg; r++ {
			uses = append(uses, r)
		}
	}
	addFixedRangeUses := func(start vm.Operand, count int) {
		if !start.IsRegister() {
			return
		}
		startReg := start.Register()
		if startReg <= 0 || count <= 0 {
			return
		}
		for r := startReg; r < startReg+count; r++ {
			uses = append(uses, r)
		}
	}

	op := inst.Opcode
	dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

	switch op {
	// No-operand terminator.
	case vm.OpReturn:
		addUse(dst)
		return

	// Moves / loads.
	case vm.OpMove:
		addUse(src1)
		addDef(dst)
		return
	case vm.OpLoadConst, vm.OpLoadParam, vm.OpLoadNone, vm.OpLoadBool, vm.OpLoadZero:
		addDef(dst)
		return
	case vm.OpLoadArray:
		addDef(dst)
		return
	case vm.OpLoadObject:
		addDef(dst)
		return
	case vm.OpLoadRange:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Simple arithmetic, comparisons, access.
	case vm.OpAdd, vm.OpSub, vm.OpMulti, vm.OpDiv, vm.OpMod,
		vm.OpCmp,
		vm.OpEq, vm.OpNe, vm.OpGt, vm.OpLt, vm.OpGte, vm.OpLte,
		vm.OpAnyEq, vm.OpAnyNe, vm.OpAnyGt, vm.OpAnyGte, vm.OpAnyLt, vm.OpAnyLte,
		vm.OpAnyIn,
		vm.OpNoneEq, vm.OpNoneNe, vm.OpNoneGt, vm.OpNoneGte, vm.OpNoneLt, vm.OpNoneLte,
		vm.OpNoneIn,
		vm.OpAllEq, vm.OpAllNe, vm.OpAllGt, vm.OpAllGte, vm.OpAllLt, vm.OpAllLte,
		vm.OpAllIn,
		vm.OpIn, vm.OpLike, vm.OpRegexp,
		vm.OpLoadIndex, vm.OpLoadIndexOptional, vm.OpLoadKey, vm.OpLoadKeyOptional,
		vm.OpLoadProperty, vm.OpLoadPropertyOptional:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Unary ops.
	case vm.OpIncr, vm.OpDecr:
		addUse(dst)
		addDef(dst)
		return
	case vm.OpCastBool, vm.OpNegate, vm.OpNot, vm.OpFlipPositive, vm.OpFlipNegative, vm.OpLength, vm.OpType:
		addUse(src1)
		addDef(dst)
		return

	// Control flow.
	case vm.OpJumpIfFalse, vm.OpJumpIfTrue:
		addUse(src1)
		return
	case vm.OpJump:
		return

	// Dataset operations.
	case vm.OpDataSet, vm.OpDataSetCollector, vm.OpDataSetSorter, vm.OpDataSetMultiSorter:
		addDef(dst)
		return
	case vm.OpPush, vm.OpArrayPush:
		addUse(dst)
		addUse(src1)
		return
	case vm.OpPushKV, vm.OpObjectSet:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		return

	// Iterators.
	case vm.OpIter:
		addUse(src1)
		addDef(dst)
		return
	case vm.OpIterValue, vm.OpIterKey:
		addUse(src1)
		addDef(dst)
		return
	case vm.OpIterLimit, vm.OpIterSkip:
		addUse(src1)
		addUse(src2)
		addDef(src1)
		return
	case vm.OpIterNext:
		addUse(src1)
		return

	// Calls.
	case vm.OpCall, vm.OpProtectedCall:
		addUse(dst)
		addRangeUses(src1, src2)
		addDef(dst)
		return
	case vm.OpCall0, vm.OpProtectedCall0:
		addUse(dst)
		addDef(dst)
		return
	case vm.OpCall1, vm.OpProtectedCall1:
		addUse(dst)
		addUse(src1)
		addDef(dst)
		return
	case vm.OpCall2, vm.OpProtectedCall2:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case vm.OpCall3, vm.OpProtectedCall3:
		addUse(dst)
		addFixedRangeUses(src1, 3)
		addDef(dst)
		return
	case vm.OpCall4, vm.OpProtectedCall4:
		addUse(dst)
		addFixedRangeUses(src1, 4)
		addDef(dst)
		return

	// Stream.
	case vm.OpStream:
		addUse(dst)
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return
	case vm.OpStreamIter:
		addUse(src1)
		addUse(src2)
		addDef(dst)
		return

	// Utility.
	case vm.OpClose:
		addUse(dst)
		addDef(dst)
		return
	case vm.OpSleep:
		addUse(dst)
		return
	}

	return
}

// computeUseDefForInstruction updates Use and Def sets for an instruction
func computeUseDefForInstruction(inst vm.Instruction, info *LivenessInfo) {
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
