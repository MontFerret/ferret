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
