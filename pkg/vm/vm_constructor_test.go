package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestNewWithOptions_InitializesFieldsFromProgramAndConfig(t *testing.T) {
	program := &bytecode.Program{
		Registers: 6,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(1)),
		},
		CatchTable: []bytecode.Catch{
			{0, 1, 42},
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Registers: 2},
				{Registers: 5},
			},
		},
	}

	instance := NewWithOptions(
		program,
		WithShapeCacheLimit(17),
		WithFastObjectDictThreshold(23),
		WithRunSafetyMode(RunSafetyFast),
	)

	if instance.program != program {
		t.Fatal("expected VM to keep source program reference")
	}

	if instance.registers == nil {
		t.Fatal("expected register file to be initialized")
	}

	if got, want := instance.registers.Size(), program.Registers; got != want {
		t.Fatalf("unexpected register file size: got %d, want %d", got, want)
	}

	if got, want := instance.runSafetyMode, RunSafetyFast; got != want {
		t.Fatalf("unexpected run safety mode: got %d, want %d", got, want)
	}

	if got, want := instance.fastObjectDictThreshold, 23; got != want {
		t.Fatalf("unexpected fast object dict threshold: got %d, want %d", got, want)
	}

	if instance.cache == nil {
		t.Fatal("expected cache to be initialized")
	}

	bytecodeLen := len(program.Bytecode)
	if got := len(instance.cache.HostFunctions); got != bytecodeLen {
		t.Fatalf("unexpected host function cache size: got %d, want %d", got, bytecodeLen)
	}

	if got := len(instance.cache.Regexps); got != bytecodeLen {
		t.Fatalf("unexpected regexp cache size: got %d, want %d", got, bytecodeLen)
	}

	if got := len(instance.cache.LoadKeyICs); got != bytecodeLen {
		t.Fatalf("unexpected load key IC cache size: got %d, want %d", got, bytecodeLen)
	}

	if got := len(instance.cache.LoadKeyConstICs); got != bytecodeLen {
		t.Fatalf("unexpected load key const IC cache size: got %d, want %d", got, bytecodeLen)
	}

	if got, want := len(instance.instructions), bytecodeLen; got != want {
		t.Fatalf("unexpected instruction wrapper size: got %d, want %d", got, want)
	}

	for i := range program.Bytecode {
		if got, want := instance.instructions[i].Instruction, program.Bytecode[i]; got != want {
			t.Fatalf("unexpected wrapped instruction at %d: got %+v, want %+v", i, got, want)
		}
	}

	wantCatchByPC := []int{0, 0}
	for i := range wantCatchByPC {
		if got := instance.catchByPC[i]; got != wantCatchByPC[i] {
			t.Fatalf("unexpected catch mapping at pc %d: got %d, want %d", i, got, wantCatchByPC[i])
		}
	}

	if got, want := len(instance.regPool.buckets), 6; got != want {
		t.Fatalf("unexpected reg pool bucket count: got %d, want %d", got, want)
	}
}

func TestBuildCatchByPC_EmptyBytecodeReturnsNil(t *testing.T) {
	got := buildCatchByPC(0, []bytecode.Catch{
		{0, 0, 1},
	})

	if got != nil {
		t.Fatalf("expected nil catch index for empty bytecode, got %#v", got)
	}
}

func TestBuildCatchByPC_ClampsAndKeepsFirstMatch(t *testing.T) {
	got := buildCatchByPC(5, []bytecode.Catch{
		{-3, 0, 10},
		{2, 100, 20},
		{3, 3, 30},
	})

	want := []int{0, -1, 1, 1, 1}
	if len(got) != len(want) {
		t.Fatalf("unexpected catch index size: got %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected catch index at pc %d: got %d, want %d", i, got[i], want[i])
		}
	}
}

func TestMaxUDFRegisters_ReturnsMaxOrZero(t *testing.T) {
	if got := maxUDFRegisters(nil); got != 0 {
		t.Fatalf("expected zero for nil list, got %d", got)
	}

	udfs := []bytecode.UDF{
		{Registers: 3},
		{Registers: 9},
		{Registers: 1},
	}

	if got := maxUDFRegisters(udfs); got != 9 {
		t.Fatalf("unexpected max UDF register count: got %d, want %d", got, 9)
	}
}
