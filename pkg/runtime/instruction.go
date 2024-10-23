package runtime

type Instruction struct {
	Opcode   OpCode
	Operands [3]int
}
