package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type execInstruction struct {
	InlineSetShape     *data.Shape
	InlineSetNextShape *data.Shape
	bytecode.Instruction
	InlineShapeID uint64
	InlineSlot    int
}

// isBackedge resolves jump metadata only for a branch that was taken, keeping
// the instruction representation compact for native straight-line execution.
func (i *execInstruction) isBackedge(pc int) bool {
	if i.Opcode == bytecode.OpMatchLoadPropertyConst {
		return i.InlineSlot <= pc
	}

	targetIndex := bytecode.JumpTargetOperandIndex(i.Opcode)
	return targetIndex >= 0 && int(i.Operands[targetIndex]) <= pc
}
