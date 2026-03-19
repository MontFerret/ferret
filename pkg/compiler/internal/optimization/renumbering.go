package optimization

import (
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

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
		inst bytecode.Instruction
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
		if ip.inst.Opcode == bytecode.OpReturn {
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
func applyRenumbering(program *bytecode.Program, renumberMap map[int]int, unsafeRegs map[int]bool) bool {
	modified := applyRegisterMap(program, renumberMap, unsafeRegs)

	if updateRegisterCount(program) {
		modified = true
	}

	return modified
}

func updateRegisterCount(program *bytecode.Program) bool {
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

func applyLinearRenumbering(program *bytecode.Program, unsafeRegs map[int]bool) bool {
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
			if !operandIsRegister(inst.Opcode, j) {
				continue
			}

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
					inst.Operands[j] = bytecode.NewRegister(newReg)
					modified = modified || changed != inst.Operands[j]
				}

				continue
			}

			if newReg, ok := mapping[reg]; ok {
				if newReg != reg {
					changed := inst.Operands[j]
					inst.Operands[j] = bytecode.NewRegister(newReg)
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

func computeLinearLiveAfter(instructions []bytecode.Instruction) []map[int]bool {
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
