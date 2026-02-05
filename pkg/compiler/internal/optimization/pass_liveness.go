package optimization

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

const LivenessAnalysisPassName = "liveness-analysis"

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
func (p *LivenessAnalysisPass) Run(program *vm.Program, cfg *ControlFlowGraph) (*PassResult, error) {
	liveness := computeLiveness(cfg, program.Registers)

	return &PassResult{
		Modified: false,
		Metadata: map[string]interface{}{
			"liveness": liveness,
		},
	}, nil
}

// computeLiveness performs dataflow analysis to compute liveness information
func computeLiveness(cfg *ControlFlowGraph, numRegisters int) map[int]*LivenessInfo {
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

// computeUseDefForInstruction updates Use and Def sets for an instruction
func computeUseDefForInstruction(inst vm.Instruction, info *LivenessInfo) {
	op := inst.Opcode
	dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

	// Determine which operands are used and which are defined
	switch op {
	case vm.OpMove, vm.OpLoadConst, vm.OpLoadParam, vm.OpLoadNone, vm.OpLoadBool, vm.OpLoadZero:
		// dst = src
		if src1.IsRegister() && !info.Def[src1.Register()] {
			info.Use[src1.Register()] = true
		}
		if dst.IsRegister() {
			info.Def[dst.Register()] = true
		}

	case vm.OpAdd, vm.OpSub, vm.OpMulti, vm.OpDiv, vm.OpMod,
		vm.OpEq, vm.OpNe, vm.OpGt, vm.OpLt, vm.OpGte, vm.OpLte,
		vm.OpAnyEq, vm.OpAnyNe, vm.OpAnyGt, vm.OpAnyGte, vm.OpAnyLt, vm.OpAnyLte,
		vm.OpNoneEq, vm.OpNoneNe, vm.OpNoneGt, vm.OpNoneGte, vm.OpNoneLt, vm.OpNoneLte,
		vm.OpAllEq, vm.OpAllNe, vm.OpAllGt, vm.OpAllGte, vm.OpAllLt, vm.OpAllLte,
		vm.OpLoadIndex, vm.OpLoadIndexOptional, vm.OpLoadKey, vm.OpLoadKeyOptional:
		// dst = src1 op src2
		if src1.IsRegister() && !info.Def[src1.Register()] {
			info.Use[src1.Register()] = true
		}
		if src2.IsRegister() && !info.Def[src2.Register()] {
			info.Use[src2.Register()] = true
		}
		if dst.IsRegister() {
			info.Def[dst.Register()] = true
		}

	case vm.OpIncr, vm.OpDecr, vm.OpNegate, vm.OpNot, vm.OpCastBool,
		vm.OpFlipPositive, vm.OpFlipNegative, vm.OpLength, vm.OpType:
		// dst = op dst (read and write same register)
		if dst.IsRegister() {
			if !info.Def[dst.Register()] {
				info.Use[dst.Register()] = true
			}
			info.Def[dst.Register()] = true
		}

	case vm.OpJumpIfFalse, vm.OpJumpIfTrue:
		// Conditional jump uses src1
		if src1.IsRegister() && !info.Def[src1.Register()] {
			info.Use[src1.Register()] = true
		}

	case vm.OpReturn:
		// Return uses dst
		if dst.IsRegister() && !info.Def[dst.Register()] {
			info.Use[dst.Register()] = true
		}

	case vm.OpPush:
		// Push uses both dst (collection) and src1 (value)
		if dst.IsRegister() && !info.Def[dst.Register()] {
			info.Use[dst.Register()] = true
		}
		if src1.IsRegister() && !info.Def[src1.Register()] {
			info.Use[src1.Register()] = true
		}

	case vm.OpIter, vm.OpIterValue, vm.OpIterKey:
		// dst = iter(src1) or dst = iter.value()
		if src1.IsRegister() && !info.Def[src1.Register()] {
			info.Use[src1.Register()] = true
		}
		if dst.IsRegister() {
			info.Def[dst.Register()] = true
		}

	case vm.OpIterNext:
		// Uses src1 (iterator), jumps to dst if done
		if src1.IsRegister() && !info.Def[src1.Register()] {
			info.Use[src1.Register()] = true
		}

	case vm.OpClose:
		// Closes resource in dst
		if dst.IsRegister() && !info.Def[dst.Register()] {
			info.Use[dst.Register()] = true
		}
	default:
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
