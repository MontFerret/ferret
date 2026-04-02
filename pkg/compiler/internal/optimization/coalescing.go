package optimization

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

// findCoalesceCandidates finds register pairs that can be coalesced
func findCoalesceCandidates(program *bytecode.Program, cfg *ControlFlowGraph, interferenceGraph map[int]map[int]bool, unsafeRegs map[int]bool, defCounts map[int]int) map[int]int {
	coalesceMap := make(map[int]int)
	coalesced := make(map[int]bool) // Track which registers have been coalesced

	// Look for move instructions between non-interfering registers
	for _, block := range cfg.Blocks {
		for _, inst := range block.Instructions {
			if inst.Opcode == bytecode.OpMove || inst.Opcode == bytecode.OpMoveTracked {
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

func countRegisterDefs(program *bytecode.Program) map[int]int {
	defCounts := make(map[int]int)

	for _, inst := range program.Bytecode {
		_, defs := instructionUseDef(inst)

		for _, reg := range defs {
			defCounts[reg]++
		}
	}

	return defCounts
}

func foldMovesIntoDefs(program *bytecode.Program, cfg *ControlFlowGraph, liveness map[int]*LivenessInfo, unsafeRegs map[int]bool) bool {
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
			if next.Opcode == bytecode.OpMove || next.Opcode == bytecode.OpMoveTracked {
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
									inst.Operands[0] = bytecode.NewRegister(dstReg)
									next.Operands[1] = bytecode.NewRegister(dstReg)
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
			if inst.Opcode == bytecode.OpLoadConst {
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
										inst.Operands[0] = bytecode.NewRegister(dstReg)

										for j := 1; j < 3; j++ {
											op := next.Operands[j]

											if op.IsRegister() && op.Register() == tmpReg {
												next.Operands[j] = bytecode.NewRegister(dstReg)
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
func applyCoalescing(program *bytecode.Program, coalesceMap map[int]int, unsafeRegs map[int]bool) bool {
	return applyRegisterMap(program, coalesceMap, unsafeRegs)
}

func applyRegisterMap(program *bytecode.Program, mapping map[int]int, unsafeRegs map[int]bool) bool {
	if len(mapping) == 0 {
		return false
	}

	modified := false

	// Replace mapped registers in all instructions.
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

			if !operandIsRegister(inst.Opcode, j) {
				continue
			}

			if op.IsRegister() {
				reg := op.Register()

				if !usedSet[reg] {
					continue
				}

				if unsafeRegs[reg] {
					continue
				}

				if newReg, ok := mapping[reg]; ok {
					if newReg == reg {
						continue
					}

					inst.Operands[j] = bytecode.NewRegister(newReg)
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

func operandIsRegister(op bytecode.Opcode, idx int) bool {
	if bytecode.IsCallOpcode(op) {
		return idx == 0 || idx == 1 || idx == 2
	}

	switch op {
	case bytecode.OpJump:
		return false
	case bytecode.OpJumpIfFalse, bytecode.OpJumpIfTrue, bytecode.OpJumpIfNone:
		return idx == 1
	case bytecode.OpJumpIfNe:
		return idx == 1 || idx == 2
	case bytecode.OpJumpIfNeConst:
		return idx == 1
	case bytecode.OpJumpIfEq:
		return idx == 1 || idx == 2
	case bytecode.OpJumpIfEqConst:
		return idx == 1
	case bytecode.OpJumpIfMissingProperty:
		return idx == 1 || idx == 2
	case bytecode.OpJumpIfMissingPropertyConst:
		return idx == 1
	case bytecode.OpIterNext:
		return idx == 1
	case bytecode.OpIterNextTimeout:
		return idx == 1 || idx == 2
	case bytecode.OpIterSkip, bytecode.OpIterLimit:
		return idx == 1 || idx == 2
	case bytecode.OpConcat:
		return idx == 0 || idx == 1
	case bytecode.OpAddConst:
		return idx == 0 || idx == 1
	case bytecode.OpLoadNone, bytecode.OpLoadZero, bytecode.OpLoadBool, bytecode.OpLoadConst, bytecode.OpLoadParam, bytecode.OpRand:
		return idx == 0
	case bytecode.OpLoadArray, bytecode.OpLoadObject:
		return idx == 0
	case bytecode.OpLoadAggregateKey:
		return idx == 0 || idx == 1
	case bytecode.OpMatchLoadPropertyConst:
		return idx == 0 || idx == 1
	case bytecode.OpAggregateUpdate:
		return idx == 0 || idx == 1
	case bytecode.OpAggregateGroupUpdate:
		return idx == 0 || idx == 1 || idx == 2
	case bytecode.OpFlatten:
		return idx == 0 || idx == 1
	case bytecode.OpDataSet, bytecode.OpDataSetCollector, bytecode.OpDataSetSorter, bytecode.OpDataSetMultiSorter:
		return idx == 0
	case bytecode.OpCounterInc:
		return idx == 0
	case bytecode.OpFail:
		return false
	case bytecode.OpFailTimeout:
		return false
	case bytecode.OpIncr, bytecode.OpDecr, bytecode.OpClose, bytecode.OpSleep, bytecode.OpReturn:
		return idx == 0
	default:
		return true
	}
}
