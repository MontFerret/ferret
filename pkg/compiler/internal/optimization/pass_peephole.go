package optimization

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

const PeepholePassName = "peephole"

// PeepholePass performs simple local rewrites to remove redundant instructions.
type PeepholePass struct{}

// NewPeepholePass creates a new peephole pass.
func NewPeepholePass() Pass {
	return &PeepholePass{}
}

// Name returns the pass name.
func (p *PeepholePass) Name() string {
	return PeepholePassName
}

// Requires returns dependencies for the pass.
func (p *PeepholePass) Requires() []string {
	return []string{}
}

// Run executes the peephole optimizations on the program.
func (p *PeepholePass) Run(ctx *PassContext) (*PassResult, error) {
	if ctx == nil || ctx.Program == nil || len(ctx.Program.Bytecode) == 0 {
		return &PassResult{Modified: false}, nil
	}

	prog := ctx.Program
	bytecodeLen := len(prog.Bytecode)
	targets := collectJumpTargets(prog.Bytecode, prog.CatchTable)
	var liveness map[int]*LivenessInfo
	var blockByInstruction []*BasicBlock

	if ctx.CFG != nil {
		liveness = computeLiveness(ctx.CFG)
		blockByInstruction = buildInstructionBlockMap(ctx.CFG, bytecodeLen)
	}

	keep := make([]bool, bytecodeLen)

	for i := range keep {
		keep[i] = true
	}

	modified := false

	for i := 0; i < bytecodeLen; i++ {
		if !keep[i] {
			continue
		}

		inst := prog.Bytecode[i]

		if i+1 < bytecodeLen && inst.Opcode == bytecode.OpEq && !targets[i] {
			next := prog.Bytecode[i+1]

			if next.Opcode == bytecode.OpJumpIfFalse && inst.Operands[0].IsRegister() && next.Operands[1].IsRegister() && next.Operands[1].Register() == inst.Operands[0].Register() {
				cmpReg := inst.Operands[0].Register()

				if !regLiveAfterInstruction(prog.Bytecode, i+1, cmpReg, blockByInstruction, liveness) {
					if i > 0 && inst.Operands[2].IsRegister() {
						prev := prog.Bytecode[i-1]
						if prev.Opcode == bytecode.OpLoadConst && prev.Operands[0].IsRegister() && inst.Operands[2].Register() == prev.Operands[0].Register() && targets[i-1] {
							continue
						}
					}

					if i > 0 && keep[i-1] && inst.Operands[2].IsRegister() {
						prev := prog.Bytecode[i-1]
						if prev.Opcode == bytecode.OpLoadConst && prev.Operands[0].IsRegister() && prev.Operands[1].IsConstant() && !targets[i-1] {
							constReg := prev.Operands[0].Register()

							if inst.Operands[2].Register() == constReg && !regLiveAfterInstruction(prog.Bytecode, i+1, constReg, blockByInstruction, liveness) {
								next.Opcode = bytecode.OpJumpIfNeConst
								next.Operands[1] = inst.Operands[1]
								next.Operands[2] = prev.Operands[1]
								prog.Bytecode[i+1] = next
								keep[i] = false
								keep[i-1] = false
								modified = true
								continue
							}
						}
					}

					if inst.Operands[2].IsConstant() {
						next.Opcode = bytecode.OpJumpIfNeConst
					} else {
						next.Opcode = bytecode.OpJumpIfNe
					}
					next.Operands[1] = inst.Operands[1]
					next.Operands[2] = inst.Operands[2]
					prog.Bytecode[i+1] = next
					keep[i] = false
					modified = true
					continue
				}
			}
		}

		if i+1 < bytecodeLen && inst.Opcode == bytecode.OpLoadConst && inst.Operands[0].IsRegister() && inst.Operands[1].IsConstant() {
			if !targets[i] {
				next := prog.Bytecode[i+1]

				if next.Opcode == bytecode.OpAdd && next.Operands[0].IsRegister() && next.Operands[1].IsRegister() && next.Operands[2].IsRegister() {
					tmpReg := inst.Operands[0].Register()
					dstReg := next.Operands[0].Register()

					if next.Operands[2].Register() == tmpReg {
						if dstReg == tmpReg || !regLiveAfterInstruction(prog.Bytecode, i+1, tmpReg, blockByInstruction, liveness) {
							next.Opcode = bytecode.OpAddConst
							next.Operands[2] = inst.Operands[1]
							prog.Bytecode[i+1] = next
							keep[i] = false
							modified = true

							continue
						}
					}
				}
			}
		}

		if isSelfMove(inst) && !targets[i] {
			keep[i] = false
			modified = true
			continue
		}

		if i+1 >= bytecodeLen || targets[i] {
			continue
		}

		if !isPureDef(inst.Opcode) {
			continue
		}

		defReg := defRegister(inst)
		if defReg == 0 {
			continue
		}

		next := prog.Bytecode[i+1]
		uses, defs := instructionUseDef(next)

		if regIn(defs, defReg) && !regIn(uses, defReg) {
			keep[i] = false
			modified = true
		}
	}

	if !modified {
		return &PassResult{Modified: false}, nil
	}

	indexMap := make([]int, bytecodeLen)
	newCode := make([]bytecode.Instruction, 0, bytecodeLen)

	for i, inst := range prog.Bytecode {
		if !keep[i] {
			indexMap[i] = -1
			continue
		}

		indexMap[i] = len(newCode)
		newCode = append(newCode, inst)
	}

	for i := range newCode {
		inst := newCode[i]
		if !isJumpOpcode(inst.Opcode) {
			continue
		}

		oldTarget := int(inst.Operands[0])
		if oldTarget < 0 || oldTarget >= bytecodeLen {
			continue
		}

		newTarget := remapIndexForward(indexMap, keep, oldTarget)
		if newTarget >= 0 && newTarget != oldTarget {
			inst.Operands[0] = bytecode.Operand(newTarget)
			newCode[i] = inst
		}
	}

	prog.Bytecode = newCode
	remapDebugSpans(prog, keep)
	remapLabels(prog, indexMap)
	remapCatchTable(prog, indexMap, keep)

	return &PassResult{
		Modified: true,
		Metadata: map[string]any{},
	}, nil
}

func collectJumpTargets(code []bytecode.Instruction, catches []bytecode.Catch) map[int]bool {
	targets := make(map[int]bool)
	for _, inst := range code {
		if !isJumpOpcode(inst.Opcode) {
			continue
		}

		target := int(inst.Operands[0])
		if target >= 0 {
			targets[target] = true
		}
	}

	for _, entry := range catches {
		if entry[2] > 0 {
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
		bytecode.OpIterNext,
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
	if inst.Opcode != bytecode.OpMove {
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

func remapIndexForward(indexMap []int, keep []bool, old int) int {
	if old < 0 || old >= len(indexMap) {
		return -1
	}
	if indexMap[old] >= 0 {
		return indexMap[old]
	}
	for i := old + 1; i < len(indexMap); i++ {
		if keep[i] {
			return indexMap[i]
		}
	}
	return -1
}

func remapIndexBackward(indexMap []int, keep []bool, old int) int {
	if old < 0 || old >= len(indexMap) {
		return -1
	}
	if indexMap[old] >= 0 {
		return indexMap[old]
	}
	for i := old - 1; i >= 0; i-- {
		if keep[i] {
			return indexMap[i]
		}
	}
	return -1
}

func remapCatchTable(prog *bytecode.Program, indexMap []int, keep []bool) {
	if prog == nil || len(prog.CatchTable) == 0 {
		return
	}

	updated := make([]bytecode.Catch, 0, len(prog.CatchTable))
	for _, entry := range prog.CatchTable {
		start := remapIndexForward(indexMap, keep, entry[0])
		end := remapIndexBackward(indexMap, keep, entry[1])
		if start < 0 || end < 0 || start > end {
			continue
		}
		jump := entry[2]
		if jump > 0 {
			jump = remapIndexForward(indexMap, keep, jump)
		}
		updated = append(updated, bytecode.Catch{start, end, jump})
	}

	prog.CatchTable = updated
}

func remapDebugSpans(prog *bytecode.Program, keep []bool) {
	if prog == nil || len(prog.Metadata.DebugSpans) == 0 {
		return
	}
	if len(prog.Metadata.DebugSpans) != len(keep) {
		return
	}
	updated := make([]file.Span, 0, len(prog.Metadata.DebugSpans))
	for i, span := range prog.Metadata.DebugSpans {
		if keep[i] {
			updated = append(updated, span)
		}
	}
	prog.Metadata.DebugSpans = updated
}

func remapLabels(prog *bytecode.Program, indexMap []int) {
	if prog == nil || len(prog.Metadata.Labels) == 0 {
		return
	}
	updated := make(map[int]string, len(prog.Metadata.Labels))
	for old, name := range prog.Metadata.Labels {
		if old < 0 || old >= len(indexMap) {
			continue
		}
		if idx := indexMap[old]; idx >= 0 {
			updated[idx] = name
		}
	}
	prog.Metadata.Labels = updated
}
