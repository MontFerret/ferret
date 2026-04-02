package optimization

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

func collectJumpTargets(prog *bytecode.Program) map[int]bool {
	targets := make(map[int]bool)
	if prog == nil {
		return targets
	}

	for _, inst := range prog.Bytecode {
		if !isJumpOpcode(inst.Opcode) {
			continue
		}

		target := int(inst.Operands[0])
		if target >= 0 {
			targets[target] = true
		}
	}

	for _, target := range prog.Metadata.MatchFailTargets {
		if target >= 0 {
			targets[target] = true
		}
	}

	for _, entry := range prog.CatchTable {
		if entry[2] >= 0 {
			targets[entry[2]] = true
		}
	}

	return targets
}

func isJumpOpcode(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpJump,
		bytecode.OpJumpIfFalse,
		bytecode.OpJumpIfTrue,
		bytecode.OpJumpIfNone,
		bytecode.OpJumpIfNe,
		bytecode.OpJumpIfNeConst,
		bytecode.OpJumpIfEq,
		bytecode.OpJumpIfEqConst,
		bytecode.OpJumpIfMissingProperty,
		bytecode.OpJumpIfMissingPropertyConst,
		bytecode.OpMatchLoadPropertyConst,
		bytecode.OpIterNext,
		bytecode.OpIterNextTimeout,
		bytecode.OpIterSkip,
		bytecode.OpIterLimit:
		return true
	default:
		return false
	}
}

func isPureDef(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpLoadConst,
		bytecode.OpLoadParam,
		bytecode.OpLoadNone,
		bytecode.OpLoadZero,
		bytecode.OpLoadBool,
		bytecode.OpMove:
		return true
	default:
		return false
	}
}

func isSelfMove(inst bytecode.Instruction) bool {
	if inst.Opcode != bytecode.OpMove && inst.Opcode != bytecode.OpMoveTracked {
		return false
	}

	dst, src := inst.Operands[0], inst.Operands[1]

	return dst.IsRegister() && src.IsRegister() && dst.Register() == src.Register()
}

func defRegister(inst bytecode.Instruction) int {
	_, defs := instructionUseDef(inst)

	if len(defs) == 1 {
		return defs[0]
	}

	return 0
}

func regIn(list []int, reg int) bool {
	for _, r := range list {
		if r == reg {
			return true
		}
	}

	return false
}

// regUsedAfter is the fallback path used when CFG/liveness information is not
// available. It performs a linear scan from start to determine whether reg is
// used before being redefined. The primary entry point is regLiveAfterInstruction,
// which prefers block-level liveness when a CFG is present.
func regUsedAfter(code []bytecode.Instruction, start int, reg int) bool {
	for i := start + 1; i < len(code); i++ {
		uses, defs := instructionUseDef(code[i])

		if regIn(uses, reg) {
			return true
		}

		if regIn(defs, reg) {
			return false
		}
	}

	return false
}

func samePureDef(a, b bytecode.Instruction) bool {
	if a.Opcode != b.Opcode {
		return false
	}

	switch a.Opcode {
	case bytecode.OpLoadConst, bytecode.OpLoadParam, bytecode.OpLoadBool, bytecode.OpMove:
		return a.Operands[0] == b.Operands[0] && a.Operands[1] == b.Operands[1]
	case bytecode.OpLoadNone, bytecode.OpLoadZero:
		return a.Operands[0] == b.Operands[0]
	default:
		return false
	}
}

func buildInstructionBlockMap(cfg *ControlFlowGraph, codeLen int) []*BasicBlock {
	if cfg == nil || codeLen <= 0 {
		return nil
	}

	byInstruction := make([]*BasicBlock, codeLen)

	for _, block := range cfg.Blocks {
		if block == nil || block == cfg.Exit || len(block.Instructions) == 0 {
			continue
		}

		start := block.Start
		end := block.End

		if start < 0 {
			start = 0
		}

		if end >= codeLen {
			end = codeLen - 1
		}

		for idx := start; idx <= end; idx++ {
			byInstruction[idx] = block
		}
	}

	return byInstruction
}

func regLiveAfterInstruction(code []bytecode.Instruction, start int, reg int, blockByInstruction []*BasicBlock, liveness map[int]*LivenessInfo) bool {
	if reg <= 0 || start < 0 || start >= len(code) {
		return false
	}

	if len(blockByInstruction) == len(code) {
		block := blockByInstruction[start]

		if block != nil {
			startInBlock := start - block.Start

			for i := startInBlock + 1; i < len(block.Instructions); i++ {
				uses, defs := instructionUseDef(block.Instructions[i])

				if regIn(uses, reg) {
					return true
				}

				if regIn(defs, reg) {
					return false
				}
			}

			if info := liveness[block.ID]; info != nil {
				return info.LiveOut[reg]
			}

			return false
		}
	}

	return regUsedAfter(code, start, reg)
}
