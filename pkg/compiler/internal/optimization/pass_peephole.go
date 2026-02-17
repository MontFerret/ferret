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
