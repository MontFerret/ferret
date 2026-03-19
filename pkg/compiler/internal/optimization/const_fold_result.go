package optimization

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type constFoldResult struct {
	newConsts map[int]runtime.Value
	modified  bool
}

func newConstFoldResult() constFoldResult {
	return constFoldResult{
		newConsts: make(map[int]runtime.Value),
	}
}

func (r *constFoldResult) setConst(reg int, val runtime.Value) {
	r.newConsts[reg] = val
}

func (r *constFoldResult) rewriteWithConst(inst *bytecode.Instruction, dst int, val runtime.Value, env constFoldEnv) {
	r.setConst(dst, val)
	r.modified = replaceWithConstLoad(inst, dst, val, env.program, env.constIndex) || r.modified
}
