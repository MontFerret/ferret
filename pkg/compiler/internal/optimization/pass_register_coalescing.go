package optimization

import (
	"sort"

	"github.com/MontFerret/ferret/pkg/vm"
)

const RegisterCoalescingPassName = "register-coalescing"

// RegisterCoalescingPass performs register coalescing to reduce register usage
// by renumbering registers based on liveness intervals
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

// liveInterval represents the live range of a register
type liveInterval struct {
	reg   int
	start int
	end   int
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

	// First, try move-based coalescing
	coalesceMap := findCoalesceCandidates(ctx.Program, ctx.CFG, interferenceGraph)

	// Apply move-based coalescing first
	moveCoalesced := applyCoalescing(ctx.Program, coalesceMap)

	// Then, perform register renumbering based on liveness intervals
	renumberMap := computeRegisterRenumbering(ctx.CFG, liveness, ctx.Program.Registers)
	renumbered := applyRenumbering(ctx.Program, renumberMap)

	modified := moveCoalesced || renumbered

	return &PassResult{
		Modified: modified,
		Metadata: map[string]interface{}{
			"registers_coalesced": len(coalesceMap),
			"coalesce_map":        coalesceMap,
			"renumber_map":        renumberMap,
		},
	}, nil
}

// computeRegisterRenumbering computes a mapping from old registers to new registers
// using a linear scan approach based on liveness intervals
func computeRegisterRenumbering(cfg *ControlFlowGraph, liveness map[int]*LivenessInfo, numRegisters int) map[int]int {
	// Build a flat list of instructions with their positions
	intervals := computeLiveIntervals(cfg, liveness, numRegisters)

	if len(intervals) == 0 {
		return nil
	}

	// Sort intervals by start position
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].start == intervals[j].start {
			return intervals[i].reg < intervals[j].reg
		}
		return intervals[i].start < intervals[j].start
	})

	// Greedy register allocation: assign the lowest available register to each interval
	renumberMap := make(map[int]int)

	// active keeps track of which new registers are in use and when they become free
	// maps new register -> end position of current interval using it
	active := make(map[int]int)

	for _, interval := range intervals {
		// Skip R0 (reserved)
		if interval.reg == 0 {
			continue
		}

		// Expire old intervals: free registers whose intervals have ended
		for newReg, endPos := range active {
			if endPos < interval.start {
				delete(active, newReg)
			}
		}

		// Find the lowest available register (starting from 1, skip R0)
		newReg := 1
		for {
			if _, inUse := active[newReg]; !inUse {
				break
			}
			newReg++
		}

		// Assign this register
		renumberMap[interval.reg] = newReg
		active[newReg] = interval.end
	}

	// Remove identity mappings
	for oldReg, newReg := range renumberMap {
		if oldReg == newReg {
			delete(renumberMap, oldReg)
		}
	}

	return renumberMap
}

// computeLiveIntervals computes the live interval for each register
func computeLiveIntervals(cfg *ControlFlowGraph, liveness map[int]*LivenessInfo, numRegisters int) []liveInterval {
	// First, flatten all instructions with their global positions
	type instrPos struct {
		inst vm.Instruction
		pos  int
	}

	var allInstructions []instrPos
	pos := 0

	for _, block := range cfg.Blocks {
		for _, inst := range block.Instructions {
			allInstructions = append(allInstructions, instrPos{inst: inst, pos: pos})
			pos++
		}
	}

	// For each register, find its definition and last use
	regDef := make(map[int]int)  // register -> first definition position
	regUse := make(map[int]int)  // register -> last use position

	for _, ip := range allInstructions {
		uses, defs := instructionUseDef(ip.inst)

		for _, reg := range defs {
			if _, exists := regDef[reg]; !exists {
				regDef[reg] = ip.pos
			}
		}

		for _, reg := range uses {
			regUse[reg] = ip.pos
		}

		// Handle return instruction specially - the returned register is used at this position
		if ip.inst.Opcode == vm.OpReturn {
			if ip.inst.Operands[0].IsRegister() {
				regUse[ip.inst.Operands[0].Register()] = ip.pos
			}
		}
	}

	// Build intervals
	var intervals []liveInterval

	for reg := 1; reg < numRegisters; reg++ {
		start, hasDef := regDef[reg]
		end, hasUse := regUse[reg]

		if !hasDef && !hasUse {
			continue
		}

		if !hasDef {
			start = 0
		}
		if !hasUse {
			end = start
		}

		// Ensure end >= start
		if end < start {
			end = start
		}

		intervals = append(intervals, liveInterval{
			reg:   reg,
			start: start,
			end:   end,
		})
	}

	return intervals
}

// applyRenumbering applies the register renumbering map to the program
func applyRenumbering(program *vm.Program, renumberMap map[int]int) bool {
	if len(renumberMap) == 0 {
		return false
	}

	modified := false
	unsafeRegs := collectRangeSensitiveRegs(program)

	// Replace registers in all instructions
	for i := range program.Bytecode {
		inst := &program.Bytecode[i]
		changed := false

		for j := 0; j < 3; j++ {
			op := inst.Operands[j]
			if op.IsRegister() {
				reg := op.Register()
				if unsafeRegs[reg] {
					continue
				}
				if newReg, ok := renumberMap[reg]; ok {
					inst.Operands[j] = vm.NewRegister(newReg)
					changed = true
				}
			}
		}

		if changed {
			modified = true
		}
	}

	// Update program's register count
	if modified {
		maxReg := 0
		for i := range program.Bytecode {
			inst := &program.Bytecode[i]
			for j := 0; j < 3; j++ {
				if inst.Operands[j].IsRegister() {
					reg := inst.Operands[j].Register()
					if reg > maxReg {
						maxReg = reg
					}
				}
			}
		}
		program.Registers = maxReg + 1
	}

	return modified
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
