package data

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

type ExecInstruction struct {
	InlineSetShape     *Shape
	InlineSetNextShape *Shape
	bytecode.Instruction
	InlineShapeID uint64
	InlineSlot    int
}
