package inspect

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func FindApplyQueryDescriptorSize(code []bytecode.Instruction, applyIdx int) (int, bool) {
	if applyIdx < 0 || applyIdx >= len(code) {
		return 0, false
	}

	queryReg := code[applyIdx].Operands[2]
	if !queryReg.IsRegister() {
		return 0, false
	}

	for i := applyIdx - 1; i >= 0; i-- {
		inst := code[i]
		if inst.Opcode != bytecode.OpLoadArray {
			continue
		}

		if !inst.Operands[0].IsRegister() || inst.Operands[0].Register() != queryReg.Register() {
			continue
		}

		return int(inst.Operands[1]), true
	}

	return 0, false
}
