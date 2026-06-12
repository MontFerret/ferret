package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestDebugTrapOpcodeIsNotKnownBytecode(t *testing.T) {
	if !isKnownOpcode(bytecode.OpReturn) {
		t.Fatal("expected return opcode to be known")
	}

	if isKnownOpcode(debugTrapOpcode) {
		t.Fatalf("debug trap opcode %d collides with a bytecode opcode", debugTrapOpcode)
	}
}
