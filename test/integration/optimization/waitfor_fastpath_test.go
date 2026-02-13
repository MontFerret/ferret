package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestWaitforFastPath_TrueSkipsSleep(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR TRUE TIMEOUT 1s`)
	assertNoOpcode(t, prog, vm.OpSleep)
	out := execOptimized(t, prog)
	if out != true {
		t.Fatalf("expected true, got %v", out)
	}
}

func TestWaitforFastPath_FalseTimeoutIsSingleSleep(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR FALSE TIMEOUT 10ms`)
	assertHasOpcode(t, prog, vm.OpSleep)
	assertNoOpcode(t, prog, vm.OpJump)
	assertNoOpcode(t, prog, vm.OpJumpIfTrue)
	assertNoOpcode(t, prog, vm.OpJumpIfFalse)
	out := execOptimized(t, prog)
	if out != false {
		t.Fatalf("expected false, got %v", out)
	}
}

func TestWaitforFastPath_ValueNoneTimeout(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR VALUE NONE TIMEOUT 10ms`)
	assertHasOpcode(t, prog, vm.OpSleep)
	assertNoOpcode(t, prog, vm.OpJump)
	assertNoOpcode(t, prog, vm.OpJumpIfTrue)
	assertNoOpcode(t, prog, vm.OpJumpIfFalse)
	out := execOptimized(t, prog)
	if out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}

func TestWaitforFastPath_ExistsEmptyArrayTimeout(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR EXISTS [] TIMEOUT 10ms`)
	assertHasOpcode(t, prog, vm.OpSleep)
	assertNoOpcode(t, prog, vm.OpJump)
	assertNoOpcode(t, prog, vm.OpJumpIfTrue)
	assertNoOpcode(t, prog, vm.OpJumpIfFalse)
	out := execOptimized(t, prog)
	if out != false {
		t.Fatalf("expected false, got %v", out)
	}
}

func TestWaitforFastPath_ExistsNonEmptyArrayImmediate(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR EXISTS [1] TIMEOUT 10ms`)
	assertNoOpcode(t, prog, vm.OpSleep)
	out := execOptimized(t, prog)
	if out != true {
		t.Fatalf("expected true, got %v", out)
	}
}

func TestWaitforFastPath_ExistsEmptyObjectTimeout(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR EXISTS {} TIMEOUT 10ms`)
	assertHasOpcode(t, prog, vm.OpSleep)
	assertNoOpcode(t, prog, vm.OpJump)
	assertNoOpcode(t, prog, vm.OpJumpIfTrue)
	assertNoOpcode(t, prog, vm.OpJumpIfFalse)
	out := execOptimized(t, prog)
	if out != false {
		t.Fatalf("expected false, got %v", out)
	}
}

func TestWaitforFastPath_ExistsNonEmptyObjectImmediate(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR EXISTS { foo: 1 } TIMEOUT 10ms`)
	assertNoOpcode(t, prog, vm.OpSleep)
	out := execOptimized(t, prog)
	if out != true {
		t.Fatalf("expected true, got %v", out)
	}
}

func TestWaitforFastPath_ExistsEmptyStringTimeout(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR EXISTS "" TIMEOUT 10ms`)
	assertHasOpcode(t, prog, vm.OpSleep)
	assertNoOpcode(t, prog, vm.OpJump)
	assertNoOpcode(t, prog, vm.OpJumpIfTrue)
	assertNoOpcode(t, prog, vm.OpJumpIfFalse)
	out := execOptimized(t, prog)
	if out != false {
		t.Fatalf("expected false, got %v", out)
	}
}

func TestWaitforFastPath_ExistsNonEmptyStringImmediate(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR EXISTS "ok" TIMEOUT 10ms`)
	assertNoOpcode(t, prog, vm.OpSleep)
	out := execOptimized(t, prog)
	if out != true {
		t.Fatalf("expected true, got %v", out)
	}
}

func TestWaitforFastPath_ValueArrayImmediate(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR VALUE [1] TIMEOUT 10ms`)
	assertNoOpcode(t, prog, vm.OpSleep)
	out := execOptimized(t, prog)
	arr, ok := out.([]any)
	if !ok || len(arr) != 1 || arr[0] != float64(1) {
		t.Fatalf("expected [1], got %v", out)
	}
}

func TestWaitforFastPath_ValueObjectImmediate(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR VALUE { foo: 1 } TIMEOUT 10ms`)
	assertNoOpcode(t, prog, vm.OpSleep)
	out := execOptimized(t, prog)
	obj, ok := out.(map[string]any)
	if !ok || obj["foo"] != float64(1) {
		t.Fatalf("expected object with foo=1, got %v", out)
	}
}

func TestWaitforFastPath_ValueStringImmediate(t *testing.T) {
	prog := compileOptimized(t, `RETURN WAITFOR VALUE "ok" TIMEOUT 10ms`)
	assertNoOpcode(t, prog, vm.OpSleep)
	out := execOptimized(t, prog)
	if out != "ok" {
		t.Fatalf("expected ok, got %v", out)
	}
}
