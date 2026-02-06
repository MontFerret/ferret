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

	unsafeRegs := collectRangeSensitiveRegs(ctx.Program)

	// Build interference graph
	interferenceGraph := buildInterferenceGraph(ctx.CFG, liveness, ctx.Program.Registers)

	// First, try move-based coalescing
	coalesceMap := findCoalesceCandidates(ctx.Program, ctx.CFG, interferenceGraph, unsafeRegs)

	// Apply move-based coalescing first
	moveCoalesced := applyCoalescing(ctx.Program, coalesceMap, unsafeRegs)

	// Then, perform register renumbering based on liveness intervals
	var renumberMap map[int]int
	var renumbered bool
	if isLinearCFG(ctx.CFG) {
		renumbered = applyLinearRenumbering(ctx.Program, unsafeRegs)
	} else {
		renumberMap = computeRegisterRenumbering(ctx.CFG, liveness, ctx.Program.Registers, unsafeRegs)
		renumbered = applyRenumbering(ctx.Program, renumberMap, unsafeRegs)
	}

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

func isLinearCFG(cfg *ControlFlowGraph) bool {
	if cfg.Entry == nil {
		return true
	}

	blocks := cfg.Blocks
	for i, block := range blocks {
		if block == cfg.Exit {
			continue
		}

		if len(block.Predecessors) > 1 || len(block.Successors) > 1 {
			return false
		}

		if len(block.Successors) == 1 {
			next := blocks[i+1]
			if block.Successors[0] != next && block.Successors[0] != cfg.Exit {
				return false
			}
		}
	}

	return true
}

func applyLinearRenumbering(program *vm.Program, unsafeRegs map[int]bool) bool {
	if len(program.Bytecode) == 0 {
		return updateRegisterCount(program)
	}

	liveAfter := computeLinearLiveAfter(program.Bytecode)
	mapping := make(map[int]int)
	inUse := make(map[int]bool)
	modified := false

	for reg := range unsafeRegs {
		mapping[reg] = reg
		inUse[reg] = true
	}

	allocate := func() int {
		newReg := 1
		for {
			if unsafeRegs[newReg] {
				newReg++
				continue
			}
			if !inUse[newReg] {
				return newReg
			}
			newReg++
		}
	}

	for i := 0; i < len(program.Bytecode); i++ {
		inst := &program.Bytecode[i]
		uses, defs := instructionUseDef(*inst)
		useSet := make(map[int]bool, len(uses))
		usedSet := make(map[int]bool, len(uses)+len(defs))

		for _, r := range uses {
			useSet[r] = true
			usedSet[r] = true
		}
		for _, r := range defs {
			usedSet[r] = true
		}

		// Ensure mappings for used registers.
		for _, reg := range uses {
			if unsafeRegs[reg] {
				continue
			}
			if _, ok := mapping[reg]; !ok {
				newReg := allocate()
				mapping[reg] = newReg
				inUse[newReg] = true
			}
		}

		liveAfterSet := liveAfter[i]

		// Collect reuse candidates from source regs that die after this instruction.
		reuseCandidates := make([]int, 0)
		for _, reg := range uses {
			if unsafeRegs[reg] {
				continue
			}
			if !liveAfterSet[reg] {
				if newReg, ok := mapping[reg]; ok {
					reuseCandidates = append(reuseCandidates, newReg)
				}
			}
		}
		sort.Ints(reuseCandidates)

		defRenames := make(map[int]int)
		for _, reg := range defs {
			if unsafeRegs[reg] {
				continue
			}
			if useSet[reg] {
				// In-place update: keep existing mapping.
				continue
			}

			if len(reuseCandidates) > 0 {
				defRenames[reg] = reuseCandidates[0]
				reuseCandidates = reuseCandidates[1:]
				continue
			}

			newReg := allocate()
			defRenames[reg] = newReg
			inUse[newReg] = true
		}

		// Rewrite operands.
		for j := 0; j < 3; j++ {
			op := inst.Operands[j]
			if !op.IsRegister() {
				continue
			}
			reg := op.Register()
			if !usedSet[reg] || unsafeRegs[reg] {
				continue
			}

			if newReg, ok := defRenames[reg]; ok {
				if newReg != reg {
					changed := inst.Operands[j]
					inst.Operands[j] = vm.NewRegister(newReg)
					modified = modified || changed != inst.Operands[j]
				}
				continue
			}

			if newReg, ok := mapping[reg]; ok {
				if newReg != reg {
					changed := inst.Operands[j]
					inst.Operands[j] = vm.NewRegister(newReg)
					modified = modified || changed != inst.Operands[j]
				}
			}
		}

		// Free mappings for regs not live after this instruction.
		for reg, newReg := range mapping {
			if unsafeRegs[reg] {
				continue
			}
			if !liveAfterSet[reg] {
				delete(mapping, reg)
				inUse[newReg] = false
			}
		}

		// Apply mappings for definitions that are live after this instruction.
		for reg, newReg := range defRenames {
			if !liveAfterSet[reg] {
				inUse[newReg] = false
				continue
			}
			mapping[reg] = newReg
			inUse[newReg] = true
		}
	}

	return updateRegisterCount(program) || modified
}

func computeLinearLiveAfter(instructions []vm.Instruction) []map[int]bool {
	live := make(map[int]bool)
	liveAfter := make([]map[int]bool, len(instructions))

	for i := len(instructions) - 1; i >= 0; i-- {
		snapshot := make(map[int]bool, len(live))
		for reg := range live {
			snapshot[reg] = true
		}
		liveAfter[i] = snapshot

		uses, defs := instructionUseDef(instructions[i])
		for _, reg := range defs {
			delete(live, reg)
		}
		for _, reg := range uses {
			live[reg] = true
		}
	}

	return liveAfter
}

// computeRegisterRenumbering computes a mapping from old registers to new registers
// using a linear scan approach based on liveness intervals
func computeRegisterRenumbering(cfg *ControlFlowGraph, liveness map[int]*LivenessInfo, numRegisters int, unsafeRegs map[int]bool) map[int]int {
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
		// Avoid renumbering range-sensitive registers
		if unsafeRegs[interval.reg] {
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
			if unsafeRegs[newReg] {
				newReg++
				continue
			}
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
	_ = liveness

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
	regDef := make(map[int]int) // register -> first definition position
	regUse := make(map[int]int) // register -> last use position

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
func applyRenumbering(program *vm.Program, renumberMap map[int]int, unsafeRegs map[int]bool) bool {
	modified := false

	if len(renumberMap) > 0 {
		// Replace registers in all instructions
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
	}

	if updateRegisterCount(program) {
		modified = true
	}

	return modified
}

func updateRegisterCount(program *vm.Program) bool {
	maxReg := 0
	for i := range program.Bytecode {
		uses, defs := instructionUseDef(program.Bytecode[i])
		for _, reg := range uses {
			if reg > maxReg {
				maxReg = reg
			}
		}
		for _, reg := range defs {
			if reg > maxReg {
				maxReg = reg
			}
		}
	}
	newCount := maxReg
	if newCount < 1 {
		newCount = 1
	}
	if program.Registers != newCount {
		program.Registers = newCount
		return true
	}
	return false
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
func findCoalesceCandidates(program *vm.Program, cfg *ControlFlowGraph, interferenceGraph map[int]map[int]bool, unsafeRegs map[int]bool) map[int]int {
	coalesceMap := make(map[int]int)
	coalesced := make(map[int]bool) // Track which registers have been coalesced

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
func applyCoalescing(program *vm.Program, coalesceMap map[int]int, unsafeRegs map[int]bool) bool {
	if len(coalesceMap) == 0 {
		return false
	}

	modified := false

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

	markRange := func(start, end vm.Operand) {
		if !start.IsRegister() || !end.IsRegister() {
			return
		}
		startReg := start.Register()
		endReg := end.Register()
		if startReg <= 0 || endReg < startReg {
			return
		}
		for r := startReg; r <= endReg; r++ {
			unsafeRegs[r] = true
		}
	}

	markFixedRange := func(start vm.Operand, count int) {
		if !start.IsRegister() {
			return
		}
		startReg := start.Register()
		if startReg <= 0 || count <= 0 {
			return
		}
		for r := startReg; r < startReg+count; r++ {
			unsafeRegs[r] = true
		}
	}

	for i := range program.Bytecode {
		inst := program.Bytecode[i]
		switch inst.Opcode {
		case vm.OpLoadArray, vm.OpLoadObject:
			markRange(inst.Operands[1], inst.Operands[2])
		case vm.OpCall, vm.OpProtectedCall:
			markRange(inst.Operands[1], inst.Operands[2])
		case vm.OpCall3, vm.OpProtectedCall3:
			markFixedRange(inst.Operands[1], 3)
		case vm.OpCall4, vm.OpProtectedCall4:
			markFixedRange(inst.Operands[1], 4)
		}
	}

	return unsafeRegs
}
