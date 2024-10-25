package runtime

type Operand int

const ResultOperand = Operand(0)

func NewConstantOperand(idx int) Operand {
	return Operand(-idx - 1)
}

func NewRegisterOperand(idx int) Operand {
	return Operand(idx)
}

func (op Operand) IsRegister() bool {
	return op >= 0
}

func (op Operand) IsConstant() bool {
	return op < 0
}

func (op Operand) Register() int {
	return int(op)
}

func (op Operand) Constant() int {
	idx := -(op + 1)

	return int(idx)
}
