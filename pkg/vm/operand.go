package vm

import "fmt"

type Operand int

// NoopOperand is a reserved operand for no operation and final results.
const NoopOperand = Operand(0)

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

func (op Operand) String() string {
	if op.IsRegister() {
		return fmt.Sprintf("R%d", op.Register())
	}

	return fmt.Sprintf("C%d", op.Constant())
}
