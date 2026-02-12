package vm

import "github.com/MontFerret/ferret/pkg/vm/internal/data"

type Instruction struct {
	Opcode   Opcode
	Operands [3]Operand

	inlineShapeID uint64
	inlineSlot    int

	inlineSetShape     *data.Shape
	inlineSetNextShape *data.Shape
}

func NewInstruction(opcode Opcode, operands ...Operand) Instruction {
	var ops [3]Operand

	switch len(operands) {
	case 3:
		ops = [3]Operand{operands[0], operands[1], operands[2]}
	case 2:
		ops = [3]Operand{operands[0], operands[1], 0}
	case 1:
		ops = [3]Operand{operands[0], 0, 0}
	default:
		ops = [3]Operand{0, 0, 0}
	}

	return Instruction{
		Opcode:   opcode,
		Operands: ops,
	}
}
