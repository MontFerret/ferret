package runtime

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Frame struct {
	registers []core.Value
	parent    *Frame
	pc        int
}

func newFrame(size, pc int, parent *Frame) *Frame {
	registers := make([]core.Value, size)
	registers[ResultOperand] = values.None

	return &Frame{
		registers: registers,
		parent:    parent,
		pc:        pc,
	}
}
