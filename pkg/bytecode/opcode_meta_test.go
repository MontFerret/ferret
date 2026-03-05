package bytecode

import "testing"

func TestOpcodeInfoCompleteness(t *testing.T) {
	for op := Opcode(0); op <= OpFail; op++ {
		info := OpcodeInfoOf(op)

		if info.Class == OpcodeClassUnknown {
			t.Fatalf("opcode %d (%s) has unknown class", op, op)
		}
	}
}

func TestOpcodeInfoCallMetadata(t *testing.T) {
	type tc struct {
		op        Opcode
		kind      CallKind
		protected bool
	}

	cases := []tc{
		{op: OpHCall, kind: CallKindHost, protected: false},
		{op: OpProtectedHCall, kind: CallKindHost, protected: true},
		{op: OpCall, kind: CallKindUser, protected: false},
		{op: OpProtectedCall, kind: CallKindUser, protected: true},
		{op: OpTailCall, kind: CallKindTail, protected: false},
	}

	for _, c := range cases {
		info := OpcodeInfoOf(c.op)
		if info.CallKind != c.kind {
			t.Fatalf("opcode %s: expected call kind %d, got %d", c.op, c.kind, info.CallKind)
		}

		if info.CallArgEncoding != CallArgEncodingRegisterRange {
			t.Fatalf("opcode %s: expected range call arg encoding", c.op)
		}

		if info.ProtectedCall != c.protected {
			t.Fatalf("opcode %s: expected protected=%v, got %v", c.op, c.protected, info.ProtectedCall)
		}
	}
}

func TestVisitCallArgumentRegisters(t *testing.T) {
	regs := make([]int, 0)

	VisitCallArgumentRegisters(OpCall, NewRegister(2), NewRegister(4), func(reg int) {
		regs = append(regs, reg)
	})

	expected := []int{2, 3, 4}
	if len(regs) != len(expected) {
		t.Fatalf("expected %d regs, got %d", len(expected), len(regs))
	}

	for i := range expected {
		if regs[i] != expected[i] {
			t.Fatalf("expected reg %d at pos %d, got %d", expected[i], i, regs[i])
		}
	}

	regs = regs[:0]
	VisitCallArgumentRegisters(OpCall, 0, 0, func(reg int) {
		regs = append(regs, reg)
	})

	if len(regs) != 0 {
		t.Fatalf("expected no call args for empty range, got %v", regs)
	}
}

func TestJumpTargetOperandIndex(t *testing.T) {
	if JumpTargetOperandIndex(OpJump) != 0 {
		t.Fatalf("expected OpJump target operand index 0")
	}

	if JumpTargetOperandIndex(OpJumpIfTrue) != 0 {
		t.Fatalf("expected OpJumpIfTrue target operand index 0")
	}

	if JumpTargetOperandIndex(OpIterNext) != 0 {
		t.Fatalf("expected OpIterNext target operand index 0")
	}

	if JumpTargetOperandIndex(OpAdd) != -1 {
		t.Fatalf("expected OpAdd to have no jump target")
	}
}
