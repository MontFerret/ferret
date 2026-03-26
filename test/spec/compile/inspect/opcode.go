package inspect

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

func FindFirstOpcodeIndex(code []bytecode.Instruction, op bytecode.Opcode) (int, bool) {
	for i, inst := range code {
		if inst.Opcode == op {
			return i, true
		}
	}

	return -1, false
}

func CountOpcode(prog *bytecode.Program, op bytecode.Opcode) int {
	if prog == nil {
		return 0
	}

	count := 0
	for _, inst := range prog.Bytecode {
		if inst.Opcode == op {
			count++
		}
	}

	return count
}

func HasOpcode(program *bytecode.Program, op bytecode.Opcode) bool {
	_, ok := FindFirstOpcodeIndex(program.Bytecode, op)

	return ok
}

func LastRegisterDefOpcodeBefore(code []bytecode.Instruction, before int, reg int) (bytecode.Opcode, bool) {
	for i := before - 1; i >= 0; i-- {
		inst := code[i]

		if !inst.Operands[0].IsRegister() {
			continue
		}

		if inst.Operands[0].Register() == reg {
			return inst.Opcode, true
		}
	}

	return bytecode.OpMove, false
}

func HasFunctionCallOpcode(program *bytecode.Program) bool {
	for _, instruction := range program.Bytecode {
		switch instruction.Opcode {
		case bytecode.OpHCall,
			bytecode.OpProtectedHCall,
			bytecode.OpCall,
			bytecode.OpProtectedCall,
			bytecode.OpTailCall:
			return true
		}
	}

	return false
}
