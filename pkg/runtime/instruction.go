package runtime

import "fmt"

type Instruction struct {
	Opcode   Opcode
	Operands [3]Operand
}

func (i Instruction) String() string {
	return fmt.Sprintf("%d %s %s %s", i.Opcode, i.Operands[0], i.Operands[1], i.Operands[2])
}
