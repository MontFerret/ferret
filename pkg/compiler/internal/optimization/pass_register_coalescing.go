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
	cfg := ctx.CFG

	// Fold trivial move chains into their defining instruction to shorten live ranges.
	folded := foldMovesIntoDefs(ctx.Program, cfg, liveness, unsafeRegs)
	if folded {
		// Recompute liveness after rewriting registers to keep analyses consistent.
		liveness = computeLiveness(cfg)
	}

	// Build interference graph
	interferenceGraph := buildInterferenceGraph(cfg, liveness, ctx.Program.Registers)

	// Count register definitions to keep move coalescing safe.
	defCounts := countRegisterDefs(ctx.Program)

	// First, try move-based coalescing
	coalesceMap := findCoalesceCandidates(ctx.Program, cfg, interferenceGraph, unsafeRegs, defCounts)

	// Apply move-based coalescing first
	moveCoalesced := applyCoalescing(ctx.Program, coalesceMap, unsafeRegs)

	if moveCoalesced {
		// Rebuild CFG and analyses after changing register assignments.
		builder := NewBuilder(ctx.Program)
		newCFG, err := builder.Build()
		if err != nil {
			return nil, ErrCFGBuildFailed
		}
		cfg = newCFG
		liveness = computeLiveness(cfg)
		unsafeRegs = collectRangeSensitiveRegs(ctx.Program)
		interferenceGraph = buildInterferenceGraph(cfg, liveness, ctx.Program.Registers)
	}

	// Then, perform register renumbering based on liveness intervals
	var renumberMap map[int]int
	var renumbered bool
	if isLinearCFG(cfg) {
		renumbered = applyLinearRenumbering(ctx.Program, unsafeRegs)
	} else {
		renumberMap = computeRegisterRenumbering(cfg, liveness, ctx.Program.Registers, unsafeRegs, interferenceGraph)
		renumbered = applyRenumbering(ctx.Program, renumberMap, unsafeRegs)
	}

	modified := folded || moveCoalesced || renumbered

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
func computeRegisterRenumbering(cfg *ControlFlowGraph, liveness map[int]*LivenessInfo, numRegisters int, unsafeRegs map[int]bool, interferenceGraph map[int]map[int]bool) map[int]int {
	intervals := computeLiveIntervals(cfg, liveness, numRegisters)

	if len(intervals) == 0 {
		return nil
	}

	// Greedy graph coloring based on interference graph.
	renumberMap := make(map[int]int)
	assigned := make(map[int]int)

	// Pre-assign unsafe registers to themselves.
	for reg := range unsafeRegs {
		if reg > 0 {
			assigned[reg] = reg
		}
	}

	// Collect registers present in intervals.
	regs := make([]int, 0, len(intervals))
	for _, interval := range intervals {
		if interval.reg == 0 || unsafeRegs[interval.reg] {
			continue
		}
		regs = append(regs, interval.reg)
	}

	// Order by descending degree, then register number for stability.
	sort.Slice(regs, func(i, j int) bool {
		di := len(interferenceGraph[regs[i]])
		dj := len(interferenceGraph[regs[j]])
		if di == dj {
			return regs[i] < regs[j]
		}
		return di > dj
	})

	for _, reg := range regs {
		used := make(map[int]bool)
		for neighbor := range interferenceGraph[reg] {
			if color, ok := assigned[neighbor]; ok {
				used[color] = true
			}
		}

		color := 1
		for {
			if unsafeRegs[color] || used[color] {
				color++
				continue
			}
			break
		}

		assigned[reg] = color
		if reg != color {
			renumberMap[reg] = color
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

	for _, block := range cfg.Blocks {
		for i, inst := range block.Instructions {
			allInstructions = append(allInstructions, instrPos{inst: inst, pos: block.Start + i})
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

	// Extend intervals using block-level liveness to cover loop back-edges.
	for _, block := range cfg.Blocks {
		info := liveness[block.ID]
		if info == nil {
			continue
		}

		for reg := range info.LiveIn {
			if reg == 0 {
				continue
			}
			if start, ok := regDef[reg]; ok {
				if block.Start < start {
					regDef[reg] = block.Start
				}
			} else {
				regDef[reg] = block.Start
			}
		}

		for reg := range info.LiveOut {
			if reg == 0 {
				continue
			}
			if end, ok := regUse[reg]; ok {
				if block.End > end {
					regUse[reg] = block.End
				}
			} else {
				regUse[reg] = block.End
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
	newCount := maxReg + 1
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

			// Registers used/defined in the same instruction interfere,
			// except for moves where dst/src coalescing is desirable.
			if inst.Opcode != vm.OpMove {
				regs := make([]int, 0, len(uses)+len(defs))
				regs = append(regs, uses...)
				regs = append(regs, defs...)
				for x := 0; x < len(regs); x++ {
					for y := x + 1; y < len(regs); y++ {
						rx, ry := regs[x], regs[y]
						if rx == ry {
							continue
						}
						graph[rx][ry] = true
						graph[ry][rx] = true
					}
				}
			}

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
func findCoalesceCandidates(program *vm.Program, cfg *ControlFlowGraph, interferenceGraph map[int]map[int]bool, unsafeRegs map[int]bool, defCounts map[int]int) map[int]int {
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
				if dstReg == srcReg {
					continue
				}

				// Avoid coalescing registers that participate in range-based ops.
				if unsafeRegs[dstReg] || unsafeRegs[srcReg] {
					continue
				}

				// Only coalesce if dst is defined exactly once (this move).
				// This avoids merging registers that carry different values
				// across control-flow paths.
				if defCounts[dstReg] != 1 {
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

func countRegisterDefs(program *vm.Program) map[int]int {
	defCounts := make(map[int]int)
	for _, inst := range program.Bytecode {
		_, defs := instructionUseDef(inst)
		for _, reg := range defs {
			defCounts[reg]++
		}
	}
	return defCounts
}

func foldMovesIntoDefs(program *vm.Program, cfg *ControlFlowGraph, liveness map[int]*LivenessInfo, unsafeRegs map[int]bool) bool {
	modified := false

	copySet := func(src map[int]bool) map[int]bool {
		if len(src) == 0 {
			return map[int]bool{}
		}
		dst := make(map[int]bool, len(src))
		for k := range src {
			dst[k] = true
		}
		return dst
	}

	for _, block := range cfg.Blocks {
		if len(block.Instructions) < 2 {
			continue
		}

		info := liveness[block.ID]
		if info == nil {
			continue
		}

		insts := block.Instructions
		live := copySet(info.LiveOut)
		liveAfter := make([]map[int]bool, len(insts))

		for i := len(insts) - 1; i >= 0; i-- {
			liveAfter[i] = copySet(live)
			uses, defs := instructionUseDef(insts[i])
			for _, def := range defs {
				delete(live, def)
			}
			for _, use := range uses {
				live[use] = true
			}
		}

		for i := 0; i < len(insts)-1; i++ {
			inst := insts[i]
			next := insts[i+1]

			// Pattern A: <op tmp ...> ; MOVE dst tmp  ->  <op dst ...> ; MOVE dst dst
			if next.Opcode == vm.OpMove {
				dst, src := next.Operands[0], next.Operands[1]
				if dst.IsRegister() && src.IsRegister() {
					dstReg, srcReg := dst.Register(), src.Register()
					if dstReg != srcReg && !unsafeRegs[dstReg] && !unsafeRegs[srcReg] {
						if !liveAfter[i+1][srcReg] {
							uses, defs := instructionUseDef(inst)
							if len(defs) == 1 && defs[0] == srcReg &&
								inst.Operands[0].IsRegister() && inst.Operands[0].Register() == srcReg {
								usesDef := false
								for _, use := range uses {
									if use == srcReg {
										usesDef = true
										break
									}
								}
								if !usesDef {
									inst.Operands[0] = vm.NewRegister(dstReg)
									next.Operands[1] = vm.NewRegister(dstReg)
									insts[i] = inst
									insts[i+1] = next
									program.Bytecode[block.Start+i] = inst
									program.Bytecode[block.Start+i+1] = next
									modified = true
								}
							}
						}
					}
				}
			}

			// Pattern B: LOADC tmp c ; OP dst ... tmp ... -> LOADC dst c ; OP dst ... dst ...
			if inst.Opcode == vm.OpLoadConst {
				dstLoad := inst.Operands[0]
				if dstLoad.IsRegister() {
					tmpReg := dstLoad.Register()
					if !unsafeRegs[tmpReg] && !liveAfter[i+1][tmpReg] {
						if next.Operands[0].IsRegister() {
							dstReg := next.Operands[0].Register()
							if dstReg != tmpReg && !unsafeRegs[dstReg] {
								uses, defs := instructionUseDef(next)
								if len(defs) == 1 && defs[0] == dstReg {
									usesDst := false
									usesTmp := false
									for _, use := range uses {
										if use == dstReg {
											usesDst = true
										}
										if use == tmpReg {
											usesTmp = true
										}
									}
									if usesTmp && !usesDst {
										inst.Operands[0] = vm.NewRegister(dstReg)
										for j := 1; j < 3; j++ {
											op := next.Operands[j]
											if op.IsRegister() && op.Register() == tmpReg {
												next.Operands[j] = vm.NewRegister(dstReg)
											}
										}
										insts[i] = inst
										insts[i+1] = next
										program.Bytecode[block.Start+i] = inst
										program.Bytecode[block.Start+i+1] = next
										modified = true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return modified
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
		case vm.OpLoadObject:
			// No range-based operands anymore.
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
