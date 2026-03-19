package optimization

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

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
	return []string{LivenessAnalysisPassName}
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
		state.liveness = retrieveLiveness(ctx)
		state.blockByInstruction = buildInstructionBlockMap(ctx.CFG, bytecodeLen)
	}

	return state
}

// retrieveLiveness attempts to reuse liveness data from a prior LivenessAnalysisPass
// stored in the pipeline metadata. Falls back to computing liveness directly if
// metadata is not available.
func retrieveLiveness(ctx *PassContext) map[int]*LivenessInfo {
	if meta, ok := ctx.Metadata[LivenessAnalysisPassName].(map[string]any); ok {
		if raw, ok := meta[LivenessAnalysis]; ok {
			if liveness, ok := raw.(map[int]*LivenessInfo); ok {
				return liveness
			}
		}
	}

	return computeLiveness(ctx.CFG)
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
