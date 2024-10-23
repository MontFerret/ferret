package runtime

type Instruction struct {
	Opcode   Opcode
	Operands [3]int
}
