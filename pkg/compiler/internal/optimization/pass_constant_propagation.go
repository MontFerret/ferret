package optimization

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const ConstantPropagationPassName = "constant-propagation"

type (
	// ConstantPropagationPass performs a simple constant propagation and folding.
	// It is conservative across control-flow merges: a register is considered constant
	// only if all predecessors agree on the same constant value.
	ConstantPropagationPass struct{}

	constState map[int]runtime.Value

	constFoldEnv struct {
		state      constState
		program    *bytecode.Program
		constIndex map[string]int
		bg         context.Context
	}
)

// NewConstantPropagationPass creates a new constant propagation pass.
func NewConstantPropagationPass() Pass {
	return &ConstantPropagationPass{}
}

func (p *ConstantPropagationPass) Name() string {
	return ConstantPropagationPassName
}

func (p *ConstantPropagationPass) Requires() []string {
	return []string{}
}

func (p *ConstantPropagationPass) Run(ctx *PassContext) (*PassResult, error) {
	if ctx == nil || ctx.Program == nil || ctx.CFG == nil {
		return &PassResult{Modified: false}, nil
	}

	if len(ctx.Program.Bytecode) == 0 {
		return &PassResult{Modified: false}, nil
	}

	constIndex := buildConstIndex(ctx.Program.Constants)
	in := make(map[int]constState)
	out := make(map[int]constState)

	modified := false
	changed := true
	bg := context.Background()

	for changed {
		changed = false

		for _, block := range ctx.CFG.Blocks {
			if block == ctx.CFG.Exit {
				continue
			}

			inState := meetPreds(block, out)
			prevIn := in[block.ID]

			if !statesEqual(prevIn, inState) {
				in[block.ID] = inState
				changed = true
			}

			state := copyState(inState)

			for i := 0; i < len(block.Instructions); i++ {
				inst := block.Instructions[i]

				if applyConstFolding(&inst, state, ctx.Program, constIndex, bg) {
					block.Instructions[i] = inst
					ctx.Program.Bytecode[block.Start+i] = inst
					modified = true
					changed = true
				}
			}

			if !statesEqual(out[block.ID], state) {
				out[block.ID] = state
				changed = true
			}
		}
	}

	return &PassResult{
		Modified: modified,
		Metadata: map[string]any{},
	}, nil
}

func copyState(in constState) constState {
	if len(in) == 0 {
		return constState{}
	}

	out := make(constState, len(in))

	for k, v := range in {
		out[k] = v
	}

	return out
}

func statesEqual(a, b constState) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bv, ok := b[k]; !ok || !constEqual(v, bv) {
			return false
		}
	}

	return true
}

func meetPreds(block *BasicBlock, out map[int]constState) constState {
	if len(block.Predecessors) == 0 {
		return constState{}
	}

	var base constState

	for _, pred := range block.Predecessors {
		if pred == nil {
			continue
		}

		base = out[pred.ID]

		break
	}

	if len(base) == 0 {
		return constState{}
	}

	result := make(constState, len(base))
	for reg, val := range base {
		keep := true

		for _, pred := range block.Predecessors {
			ps := out[pred.ID]
			pv, ok := ps[reg]

			if !ok || !constEqual(val, pv) {
				keep = false
				break
			}
		}

		if keep {
			result[reg] = val
		}
	}

	return result
}

func applyConstFolding(inst *bytecode.Instruction, state constState, program *bytecode.Program, constIndex map[string]int, bg context.Context) bool {
	defs := defsOnly(*inst)
	env := constFoldEnv{
		state:      state,
		program:    program,
		constIndex: constIndex,
		bg:         bg,
	}
	result := foldInstructionByOpcode(inst, env)
	applyConstUpdates(state, defs, result.newConsts)

	return result.modified
}

func applyConstUpdates(state constState, defs []int, newConsts map[int]runtime.Value) {
	for _, reg := range defs {
		if val, ok := newConsts[reg]; ok {
			state[reg] = val
		} else {
			delete(state, reg)
		}
	}
}

func defsOnly(inst bytecode.Instruction) []int {
	_, defs := instructionUseDef(inst)

	return defs
}
