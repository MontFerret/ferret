package bytecode

import "fmt"

type Instruction struct {
	Opcode   Opcode     `json:"opcode"`
	Operands [3]Operand `json:"operands"`
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

func (i Instruction) String() string {
	return fmt.Sprintf("%s %d %d %d", i.Opcode, i.Operands[0], i.Operands[1], i.Operands[2])
}
