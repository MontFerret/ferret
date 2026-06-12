package vm

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

const debugTrapOpcode = bytecode.Opcode(255)

func init() {
	if isKnownOpcode(debugTrapOpcode) {
		panic("vm debug trap opcode collides with bytecode opcode")
	}
}

func isKnownOpcode(op bytecode.Opcode) bool {
	return bytecode.OpcodeInfoOf(op).Class != bytecode.OpcodeClassUnknown
}
