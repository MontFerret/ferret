package cfg

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

// RegisterCoalescingPass performs register coalescing to reduce register usage
// by merging registers that don't interfere with each other
type RegisterCoalescingPass struct{}

// NewRegisterCoalescingPass creates a new register coalescing pass
func NewRegisterCoalescingPass() *RegisterCoalescingPass {
	return &RegisterCoalescingPass{}
}

// Name returns the pass name
func (p *RegisterCoalescingPass) Name() string {
	return "register-coalescing"
}

// Run executes register coalescing on the program
func (p *RegisterCoalescingPass) Run(program *vm.Program, cfg *ControlFlowGraph) (*PassResult, error) {
	// First, run liveness analysis to get interference information
	livenessPass := NewLivenessAnalysisPass()
	livenessResult, err := livenessPass.Run(program, cfg)
	if err != nil {
		return nil, err
	}

	liveness := livenessResult.Metadata["liveness"].(map[int]*LivenessInfo)

	// Build interference graph
	interferenceGraph := buildInterferenceGraph(cfg, liveness, program.Registers)

	// Find register pairs that can be coalesced (connected by moves and don't interfere)
	coalesceMap := findCoalesceCandidates(cfg, interferenceGraph)

	if len(coalesceMap) == 0 {
		return &PassResult{
			Modified: false,
			Metadata: map[string]interface{}{
				"registers_coalesced": 0,
			},
		}, nil
	}

	// Apply register coalescing
	modified := applyCoalescing(program, coalesceMap)

	return &PassResult{
		Modified: modified,
		Metadata: map[string]interface{}{
			"registers_coalesced": len(coalesceMap),
			"coalesce_map":        coalesceMap,
		},
	}, nil
}

// buildInterferenceGraph creates an interference graph from liveness information
func buildInterferenceGraph(cfg *ControlFlowGraph, liveness map[int]*LivenessInfo, numRegisters int) map[int]map[int]bool {
	graph := make(map[int]map[int]bool)

	// Initialize graph
	for i := 0; i < numRegisters; i++ {
		graph[i] = make(map[int]bool)
	}

	// Two registers interfere if they are both live at the same program point
	for _, block := range cfg.Blocks {
		info := liveness[block.ID]

		// Check interference at block entry
		liveRegs := make([]int, 0)
		for reg := range info.LiveIn {
			liveRegs = append(liveRegs, reg)
		}

		// Add interference edges
		for i := 0; i < len(liveRegs); i++ {
			for j := i + 1; j < len(liveRegs); j++ {
				r1, r2 := liveRegs[i], liveRegs[j]
				graph[r1][r2] = true
				graph[r2][r1] = true
			}
		}
	}

	return graph
}

// findCoalesceCandidates finds register pairs that can be coalesced
func findCoalesceCandidates(cfg *ControlFlowGraph, interferenceGraph map[int]map[int]bool) map[int]int {
	coalesceMap := make(map[int]int)

	// Look for move instructions between non-interfering registers
	for _, block := range cfg.Blocks {
		for _, inst := range block.Instructions {
			if inst.Opcode == vm.OpMove {
				dst, src := inst.Operands[0], inst.Operands[1]

				if !dst.IsRegister() || !src.IsRegister() {
					continue
				}

				dstReg, srcReg := dst.Register(), src.Register()

				// Check if they don't interfere and haven't been coalesced yet
				if !interferenceGraph[dstReg][srcReg] &&
					coalesceMap[dstReg] == 0 && coalesceMap[srcReg] == 0 {
					// Coalesce dst into src (replace all uses of dst with src)
					coalesceMap[dstReg] = srcReg
				}
			}
		}
	}

	return coalesceMap
}

// applyCoalescing applies the coalescing map to the program
func applyCoalescing(program *vm.Program, coalesceMap map[int]int) bool {
	if len(coalesceMap) == 0 {
		return false
	}

	modified := false

	// Replace coalesced registers in all instructions
	for i := range program.Bytecode {
		inst := &program.Bytecode[i]
		changed := false

		// Check and replace each operand
		for j := 0; j < 3; j++ {
			op := inst.Operands[j]
			if op.IsRegister() {
				reg := op.Register()
				if newReg, ok := coalesceMap[reg]; ok {
					inst.Operands[j] = vm.NewRegister(newReg)
					changed = true
				}
			}
		}

		if changed {
			modified = true
		}
	}

	return modified
}
