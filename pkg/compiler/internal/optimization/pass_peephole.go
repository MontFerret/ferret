package optimization

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

const PeepholePassName = "peephole"

type (
	// PeepholePass performs simple local rewrites to remove redundant instructions.
	PeepholePass struct{}

	peepholeRunState struct {
		prog               *bytecode.Program
		targets            map[int]bool
		liveness           map[int]*LivenessInfo
		blockByInstruction []*BasicBlock
		keep               []bool
		bytecodeLen        int
	}
)

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

	state := newPeepholeRunState(ctx)
	modified := runPeepholeScanAndRewrite(state)

	if !modified {
		return &PassResult{Modified: false}, nil
	}

	applyPeepholeCompactionAndRemap(state)

	return &PassResult{
		Modified: true,
		Metadata: map[string]any{},
	}, nil
}

func newPeepholeRunState(ctx *PassContext) *peepholeRunState {
	prog := ctx.Program
	bytecodeLen := len(prog.Bytecode)
	state := &peepholeRunState{
		prog:        prog,
		bytecodeLen: bytecodeLen,
		targets:     collectJumpTargets(prog),
		keep:        make([]bool, bytecodeLen),
	}

	for i := range state.keep {
		state.keep[i] = true
	}

	if ctx.CFG != nil {
		state.liveness = computeLiveness(ctx.CFG)
		state.blockByInstruction = buildInstructionBlockMap(ctx.CFG, bytecodeLen)
	}

	return state
}

func runPeepholeScanAndRewrite(state *peepholeRunState) bool {
	modified := false

	for i := 0; i < state.bytecodeLen; i++ {
		if !state.keep[i] {
			continue
		}

		if tryRewriteComparisonJump(state, i) {
			modified = true
			continue
		}

		if tryRewriteAddConst(state, i) {
			modified = true
			continue
		}

		if tryRemoveSelfMove(state, i) {
			modified = true
			continue
		}

		if tryRemoveDeadOrRedundantPureDefs(state, i) {
			modified = true
		}
	}

	return modified
}

func tryRewriteComparisonJump(state *peepholeRunState, index int) bool {
	if index+1 >= state.bytecodeLen || state.targets[index] {
		return false
	}

	inst := state.prog.Bytecode[index]
	next := state.prog.Bytecode[index+1]

	if !isComparisonJumpPair(inst, next) {
		return false
	}

	cmpReg := inst.Operands[0].Register()
	if regLiveAfterInstruction(state.prog.Bytecode, index+1, cmpReg, state.blockByInstruction, state.liveness) {
		return false
	}

	if hasTargetedPreviousConstLoad(state, index, inst) {
		return false
	}

	newOp, newOpConst, ok := resolveComparisonJumpOpcode(inst.Opcode, next.Opcode)
	if !ok {
		return false
	}

	if rewritten, ok := rewriteJumpUsingPreviousConstLoad(state, index, inst, next, newOpConst); ok {
		state.prog.Bytecode[index+1] = rewritten
		state.keep[index] = false
		state.keep[index-1] = false
		return true
	}

	if inst.Operands[2].IsConstant() {
		next.Opcode = newOpConst
	} else {
		next.Opcode = newOp
	}
	next.Operands[1] = inst.Operands[1]
	next.Operands[2] = inst.Operands[2]
	state.prog.Bytecode[index+1] = next
	state.keep[index] = false

	return true
}

func isComparisonJumpPair(compareInst, jumpInst bytecode.Instruction) bool {
	if compareInst.Opcode != bytecode.OpEq && compareInst.Opcode != bytecode.OpNe {
		return false
	}

	if jumpInst.Opcode != bytecode.OpJumpIfFalse && jumpInst.Opcode != bytecode.OpJumpIfTrue {
		return false
	}

	if !compareInst.Operands[0].IsRegister() || !jumpInst.Operands[1].IsRegister() {
		return false
	}

	return jumpInst.Operands[1].Register() == compareInst.Operands[0].Register()
}

func resolveComparisonJumpOpcode(compareOp, jumpOp bytecode.Opcode) (bytecode.Opcode, bytecode.Opcode, bool) {
	if (compareOp != bytecode.OpEq && compareOp != bytecode.OpNe) || (jumpOp != bytecode.OpJumpIfFalse && jumpOp != bytecode.OpJumpIfTrue) {
		return 0, 0, false
	}

	if (compareOp == bytecode.OpEq && jumpOp == bytecode.OpJumpIfTrue) || (compareOp == bytecode.OpNe && jumpOp == bytecode.OpJumpIfFalse) {
		return bytecode.OpJumpIfEq, bytecode.OpJumpIfEqConst, true
	}

	return bytecode.OpJumpIfNe, bytecode.OpJumpIfNeConst, true
}

func hasTargetedPreviousConstLoad(state *peepholeRunState, index int, inst bytecode.Instruction) bool {
	if index == 0 || !inst.Operands[2].IsRegister() {
		return false
	}

	prev := state.prog.Bytecode[index-1]
	if prev.Opcode != bytecode.OpLoadConst || !prev.Operands[0].IsRegister() {
		return false
	}

	return inst.Operands[2].Register() == prev.Operands[0].Register() && state.targets[index-1]
}

func rewriteJumpUsingPreviousConstLoad(state *peepholeRunState, index int, compareInst, jumpInst bytecode.Instruction, jumpConstOp bytecode.Opcode) (bytecode.Instruction, bool) {
	if index == 0 || !state.keep[index-1] || !compareInst.Operands[2].IsRegister() {
		return bytecode.Instruction{}, false
	}

	prev := state.prog.Bytecode[index-1]
	if prev.Opcode != bytecode.OpLoadConst || !prev.Operands[0].IsRegister() || !prev.Operands[1].IsConstant() || state.targets[index-1] {
		return bytecode.Instruction{}, false
	}

	constReg := prev.Operands[0].Register()
	if compareInst.Operands[2].Register() != constReg {
		return bytecode.Instruction{}, false
	}

	if regLiveAfterInstruction(state.prog.Bytecode, index+1, constReg, state.blockByInstruction, state.liveness) {
		return bytecode.Instruction{}, false
	}

	jumpInst.Opcode = jumpConstOp
	jumpInst.Operands[1] = compareInst.Operands[1]
	jumpInst.Operands[2] = prev.Operands[1]

	return jumpInst, true
}

func tryRewriteAddConst(state *peepholeRunState, index int) bool {
	if index+1 >= state.bytecodeLen || state.targets[index] {
		return false
	}

	loadInst := state.prog.Bytecode[index]
	if loadInst.Opcode != bytecode.OpLoadConst || !loadInst.Operands[0].IsRegister() || !loadInst.Operands[1].IsConstant() {
		return false
	}

	next := state.prog.Bytecode[index+1]
	if next.Opcode != bytecode.OpAdd || !next.Operands[0].IsRegister() || !next.Operands[1].IsRegister() || !next.Operands[2].IsRegister() {
		return false
	}

	tempReg := loadInst.Operands[0].Register()
	if next.Operands[2].Register() != tempReg {
		return false
	}

	dstReg := next.Operands[0].Register()
	if dstReg != tempReg && regLiveAfterInstruction(state.prog.Bytecode, index+1, tempReg, state.blockByInstruction, state.liveness) {
		return false
	}

	next.Opcode = bytecode.OpAddConst
	next.Operands[2] = loadInst.Operands[1]
	state.prog.Bytecode[index+1] = next
	state.keep[index] = false

	return true
}

func tryRemoveSelfMove(state *peepholeRunState, index int) bool {
	if state.targets[index] || !isSelfMove(state.prog.Bytecode[index]) {
		return false
	}

	state.keep[index] = false

	return true
}

func tryRemoveDeadOrRedundantPureDefs(state *peepholeRunState, index int) bool {
	if index+1 >= state.bytecodeLen || state.targets[index] {
		return false
	}

	inst := state.prog.Bytecode[index]
	if !isPureDef(inst.Opcode) {
		return false
	}

	defReg := defRegister(inst)
	if defReg == 0 {
		return false
	}

	next := state.prog.Bytecode[index+1]
	uses, defs := instructionUseDef(next)

	if isPureDef(next.Opcode) && defRegister(next) == defReg && samePureDef(inst, next) && !state.targets[index+1] {
		state.keep[index+1] = false
		return true
	}

	if regIn(defs, defReg) && !regIn(uses, defReg) {
		state.keep[index] = false
		return true
	}

	return false
}

func applyPeepholeCompactionAndRemap(state *peepholeRunState) {
	newCode, indexMap := compactPeepholeInstructions(state.prog.Bytecode, state.keep)
	remapPeepholeJumps(newCode, indexMap, state.keep, state.bytecodeLen)

	state.prog.Bytecode = newCode
	remapAggregateSelectorSlots(state.prog, state.keep)
	remapMatchFailTargets(state.prog, indexMap, state.keep)
	remapDebugSpans(state.prog, state.keep)
	remapLabels(state.prog, indexMap)
	remapUdfEntries(state.prog, indexMap, state.keep)
	remapCatchTable(state.prog, indexMap, state.keep)
}

func compactPeepholeInstructions(code []bytecode.Instruction, keep []bool) ([]bytecode.Instruction, []int) {
	indexMap := make([]int, len(code))
	newCode := make([]bytecode.Instruction, 0, len(code))

	for i, inst := range code {
		if !keep[i] {
			indexMap[i] = -1
			continue
		}

		indexMap[i] = len(newCode)
		newCode = append(newCode, inst)
	}

	return newCode, indexMap
}

func remapPeepholeJumps(code []bytecode.Instruction, indexMap []int, keep []bool, oldCodeLen int) {
	for i := range code {
		inst := code[i]
		if !isJumpOpcode(inst.Opcode) {
			continue
		}
		if inst.Opcode == bytecode.OpMatchLoadPropertyConst {
			continue
		}

		oldTarget := int(inst.Operands[0])
		if oldTarget < 0 || oldTarget >= oldCodeLen {
			continue
		}

		newTarget := remapIndexForward(indexMap, keep, oldTarget)
		if newTarget >= 0 && newTarget != oldTarget {
			inst.Operands[0] = bytecode.Operand(newTarget)
			code[i] = inst
		}
	}
}

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
		if jump >= 0 {
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

func remapAggregateSelectorSlots(prog *bytecode.Program, keep []bool) {
	if prog == nil || len(prog.Metadata.AggregateSelectorSlots) == 0 {
		return
	}

	if len(prog.Metadata.AggregateSelectorSlots) != len(keep) {
		return
	}

	updated := make([]int, 0, len(prog.Metadata.AggregateSelectorSlots))

	for i, slot := range prog.Metadata.AggregateSelectorSlots {
		if keep[i] {
			updated = append(updated, slot)
		}
	}

	prog.Metadata.AggregateSelectorSlots = updated
}

func remapMatchFailTargets(prog *bytecode.Program, indexMap []int, keep []bool) {
	if prog == nil || len(prog.Metadata.MatchFailTargets) == 0 {
		return
	}

	if len(prog.Metadata.MatchFailTargets) != len(keep) {
		return
	}

	updated := make([]int, 0, len(prog.Metadata.MatchFailTargets))

	for i, target := range prog.Metadata.MatchFailTargets {
		if !keep[i] {
			continue
		}

		if target >= 0 {
			target = remapIndexForward(indexMap, keep, target)
		}

		updated = append(updated, target)
	}

	prog.Metadata.MatchFailTargets = updated
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

func remapUdfEntries(prog *bytecode.Program, indexMap []int, keep []bool) {
	if prog == nil || len(prog.Functions.UserDefined) == 0 {
		return
	}

	for i := range prog.Functions.UserDefined {
		entry := prog.Functions.UserDefined[i].Entry
		if entry < 0 || entry >= len(indexMap) {
			continue
		}

		newEntry := indexMap[entry]
		if newEntry < 0 {
			newEntry = remapIndexForward(indexMap, keep, entry)
		}

		if newEntry >= 0 {
			prog.Functions.UserDefined[i].Entry = newEntry
		}
	}
}
