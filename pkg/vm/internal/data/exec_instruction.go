package data

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

type ExecInstruction struct {
	bytecode.Instruction

	InlineShapeID      uint64
	InlineSlot         int
	InlineSetShape     *Shape
	InlineSetNextShape *Shape
}
