package runtime

import "github.com/MontFerret/ferret/pkg/runtime/core"

type Frame struct {
	registers []core.Value
	parent    *Frame
	pc        int
}

func newFrame(size, pc int, parent *Frame) *Frame {
	return &Frame{
		registers: make([]core.Value, size),
		parent:    parent,
		pc:        pc,
	}
}
