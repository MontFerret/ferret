package optimization

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

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
			if inst.Opcode != bytecode.OpMove && inst.Opcode != bytecode.OpMoveTracked {
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

func collectPinnedRegs(program *bytecode.Program) map[int]bool {
	unsafeRegs := make(map[int]bool)

	for _, udf := range program.Functions.UserDefined {
		for reg := 1; reg <= udf.Params; reg++ {
			unsafeRegs[reg] = true
		}
	}

	markFixedRange := func(start bytecode.Operand, count int) {
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

		if inst.Opcode == bytecode.OpMakeCell {
			dst := inst.Operands[0]

			if dst.IsRegister() && dst.Register() > 0 {
				unsafeRegs[dst.Register()] = true
			}
		}

		if bytecode.IsCallOpcode(inst.Opcode) {
			start := inst.Operands[1]
			end := inst.Operands[2]

			if start.IsRegister() && end.IsRegister() {
				startReg := start.Register()
				endReg := end.Register()

				// Keep multi-arg call ranges fixed so [src1..src2] remains valid
				// after any register remapping. Single-arg calls are safe to rename.
				if startReg > 0 && endReg > startReg {
					for r := startReg; r <= endReg; r++ {
						unsafeRegs[r] = true
					}
				}
			}
		}

		if inst.Opcode == bytecode.OpConcat {
			markFixedRange(inst.Operands[1], int(inst.Operands[2]))
		}
	}

	return unsafeRegs
}
