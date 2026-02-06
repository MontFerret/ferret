package optimization

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

const RegisterCoalescingPassName = "register-coalescing"

// RegisterCoalescingPass performs register coalescing to reduce register usage
// by merging registers that don't interfere with each other
type RegisterCoalescingPass struct{}

// NewRegisterCoalescingPass creates a new register coalescing pass
func NewRegisterCoalescingPass() Pass {
	return &RegisterCoalescingPass{}
}

// Name returns the pass name
func (p *RegisterCoalescingPass) Name() string {
	return RegisterCoalescingPassName
}

func (p *RegisterCoalescingPass) Requires() []string {
	return []string{LivenessAnalysisPassName}
}

// Run executes register coalescing on the program
func (p *RegisterCoalescingPass) Run(ctx *PassContext) (*PassResult, error) {
	meta, ok := ctx.Metadata[LivenessAnalysisPassName].(map[string]any)

	if !ok {
		return nil, ErrMissingDependency
	}
	raw, ok := meta[LivenessAnalysis]

	if !ok {
		return nil, ErrMissingDependency
	}

	liveness, ok := raw.(map[int]*LivenessInfo)

	if !ok {
		return nil, ErrMissingDependency
	}

	// Build interference graph
	interferenceGraph := buildInterferenceGraph(ctx.CFG, liveness, ctx.Program.Registers)

	// Find register pairs that can be coalesced (connected by moves and don't interfere)
	coalesceMap := findCoalesceCandidates(ctx.Program, ctx.CFG, interferenceGraph)

	if len(coalesceMap) == 0 {
		return &PassResult{
			Modified: false,
			Metadata: map[string]interface{}{
				"registers_coalesced": 0,
			},
		}, nil
	}

	// Apply register coalescing
	modified := applyCoalescing(ctx.Program, coalesceMap)

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

		// Start with live-out at block end and walk instructions backward.
		live := make(map[int]bool, len(info.LiveOut))
		for reg := range info.LiveOut {
			live[reg] = true
		}

		for i := len(block.Instructions) - 1; i >= 0; i-- {
			inst := block.Instructions[i]
			uses, defs := instructionUseDef(inst)

			// Add interferences: defs interfere with all currently-live regs.
			for _, def := range defs {
				for reg := range live {
					if reg == def {
						continue
					}
					graph[def][reg] = true
					graph[reg][def] = true
				}
			}

			// Update live = (live - defs) U uses
			for _, def := range defs {
				delete(live, def)
			}
			for _, use := range uses {
				live[use] = true
			}
		}
	}

	return graph
}

// findCoalesceCandidates finds register pairs that can be coalesced
func findCoalesceCandidates(program *vm.Program, cfg *ControlFlowGraph, interferenceGraph map[int]map[int]bool) map[int]int {
	coalesceMap := make(map[int]int)
	coalesced := make(map[int]bool) // Track which registers have been coalesced
	unsafeRegs := collectRangeSensitiveRegs(program)

	// Look for move instructions between non-interfering registers
	for _, block := range cfg.Blocks {
		for _, inst := range block.Instructions {
			if inst.Opcode == vm.OpMove {
				dst, src := inst.Operands[0], inst.Operands[1]

				if !dst.IsRegister() || !src.IsRegister() {
					continue
				}

				dstReg, srcReg := dst.Register(), src.Register()

				// Avoid coalescing registers that participate in range-based ops.
				if unsafeRegs[dstReg] || unsafeRegs[srcReg] {
					continue
				}

				// Check if they don't interfere and haven't been coalesced yet
				if !interferenceGraph[dstReg][srcReg] &&
					!coalesced[dstReg] && !coalesced[srcReg] {
					// Coalesce dst into src (replace all uses of dst with src)
					coalesceMap[dstReg] = srcReg
					coalesced[dstReg] = true
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
	unsafeRegs := collectRangeSensitiveRegs(program)

	// Replace coalesced registers in all instructions
	for i := range program.Bytecode {
		inst := &program.Bytecode[i]
		changed := false
		uses, defs := instructionUseDef(*inst)
		usedSet := make(map[int]bool, len(uses)+len(defs))

		for _, r := range uses {
			usedSet[r] = true
		}

		for _, r := range defs {
			usedSet[r] = true
		}

		// Check and replace each operand
		for j := 0; j < 3; j++ {
			op := inst.Operands[j]
			if op.IsRegister() {
				reg := op.Register()
				if !usedSet[reg] {
					continue
				}
				if unsafeRegs[reg] {
					continue
				}
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

func collectRangeSensitiveRegs(program *vm.Program) map[int]bool {
	unsafeRegs := make(map[int]bool)

	mark := func(op vm.Operand) {
		if op.IsRegister() {
			unsafeRegs[op.Register()] = true
		}
	}

	for i := range program.Bytecode {
		inst := program.Bytecode[i]
		switch inst.Opcode {
		case vm.OpLoadArray, vm.OpLoadObject:
			mark(inst.Operands[1])
			mark(inst.Operands[2])
		case vm.OpCall, vm.OpProtectedCall:
			mark(inst.Operands[1])
			mark(inst.Operands[2])
		case vm.OpCall3, vm.OpProtectedCall3, vm.OpCall4, vm.OpProtectedCall4:
			mark(inst.Operands[1])
		}
	}

	return unsafeRegs
}
