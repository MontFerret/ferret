package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func ReadOperandValue(reg []runtime.Value, constants []runtime.Value, operand bytecode.Operand) runtime.Value {
	if operand.IsConstant() {
		return constants[operand.Constant()]
	}

	return reg[operand]
}
