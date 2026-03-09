package vm

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
)

func compileProgram(t *testing.T, source string) *bytecode.Program {
	t.Helper()

	c := compiler.New()
	program, err := c.Compile(file.NewSource("test", source))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	return program
}

func TestPanicPolicyRecoversPanics(t *testing.T) {
	program := compileProgram(t, "RETURN PANIC_FN()")

	instance := New(program)
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("PANIC_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			panic("panic in host function")
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	_, err = instance.Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected runtime error in strict mode")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected wrapped runtime error, got %T", err)
	}
}

func TestPanicPolicyPropagatesPanics(t *testing.T) {
	program := compileProgram(t, "RETURN PANIC_FN()")

	instance := NewWith(program, WithPanicPolicy(PanicPropagate))
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("PANIC_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			panic("panic in host function")
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected panic to propagate in fast mode")
		}
	}()

	_, _ = instance.Run(context.Background(), env)
}

func TestPanicPolicyPropagateStillWrapsReturnedErrors(t *testing.T) {
	program := compileProgram(t, "RETURN FAIL_FN()")

	instance := NewWith(program, WithPanicPolicy(PanicPropagate))
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("FAIL_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.None, errors.New("boom")
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	_, err = instance.Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected wrapped runtime error, got %T", err)
	}
}

func TestNewWith_InitializesFieldsFromProgramAndConfig(t *testing.T) {
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

	instance := NewWith(
		program,
		WithShapeCacheLimit(17),
		WithFastObjectDictThreshold(23),
		WithPanicPolicy(PanicPropagate),
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

	if got, want := instance.options.panicPolicy, PanicPropagate; got != want {
		t.Fatalf("unexpected panic policy: got %d, want %d", got, want)
	}

	if got, want := instance.options.fastObjectDictThreshold, 23; got != want {
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

	reg := instance.frames.AcquireRegisters(5)
	if got, want := len(reg), 5; got != want {
		t.Fatalf("unexpected pooled register size: got %d, want %d", got, want)
	}

	instance.frames.ReleaseRegisters(reg)
	reused := instance.frames.AcquireRegisters(5)
	if got, want := len(reused), 5; got != want {
		t.Fatalf("unexpected reused register size: got %d, want %d", got, want)
	}

	if &reg[0] != &reused[0] {
		t.Fatal("expected register pool to reuse buffers")
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

func TestOpLoadParam_UsesBoundSlots(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  4,
		Params:     []string{"foo", "bar"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(2), bytecode.Operand(2)),
			bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
		},
	}

	env := NewDefaultEnvironment()
	env.Params["foo"] = runtime.NewInt(1)
	env.Params["bar"] = runtime.NewInt(2)

	out, err := New(program).Run(context.Background(), env)
	if err != nil {
		t.Fatalf("expected successful run, got %v", err)
	}

	if out != runtime.NewInt(3) {
		t.Fatalf("unexpected result: got %v, want %v", out, runtime.NewInt(3))
	}
}

func TestOpLoadParam_MissingParamsPreserveRuntimeError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Params:     []string{"foo", "bar"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	env := NewDefaultEnvironment()
	env.Params["foo"] = runtime.NewInt(1)

	_, err := New(program).Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected missing parameter error")
	}

	if !strings.Contains(err.Error(), "Missing parameter") {
		t.Fatalf("unexpected error: %v", err)
	}

	cause := errors.Unwrap(err)
	if cause == nil {
		cause = err
	}

	if !strings.Contains(cause.Error(), "@bar") {
		t.Fatalf("expected missing parameter name in error, got %v (cause: %v)", err, cause)
	}
}

func TestUnwindToProtected_ReclaimsDiscardedFrameRegisters(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Registers: 6},
			},
		},
	})

	lowerRegs := make([]runtime.Value, 2)
	protectedRegs := make([]runtime.Value, 3)
	aboveRegs1 := make([]runtime.Value, 4)
	aboveRegs2 := make([]runtime.Value, 5)
	activeRegs := make([]runtime.Value, 6)

	protectedRegs[1] = runtime.True

	instance.registers.Values = activeRegs
	instance.frames.Push(frame.CallFrame{
		ReturnPC:   10,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  lowerRegs,
		Protected:  false,
	})
	instance.frames.Push(frame.CallFrame{
		ReturnPC:   20,
		ReturnDest: bytecode.NewRegister(1),
		Registers:  protectedRegs,
		Protected:  true,
	})
	instance.frames.Push(frame.CallFrame{
		ReturnPC:   30,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  aboveRegs1,
		Protected:  false,
	})
	instance.frames.Push(frame.CallFrame{
		ReturnPC:   40,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  aboveRegs2,
		Protected:  false,
	})

	if ok := instance.unwindToProtected(); !ok {
		t.Fatal("expected protected unwind to succeed")
	}

	if got, want := instance.pc, 20; got != want {
		t.Fatalf("unexpected pc after unwind: got %d, want %d", got, want)
	}

	if got, want := instance.frames.Len(), 1; got != want {
		t.Fatalf("unexpected frame depth after unwind: got %d, want %d", got, want)
	}

	remaining := instance.frames.Top()
	if remaining == nil {
		t.Fatal("expected remaining frame after unwind")
	}

	if got, want := remaining.ReturnPC, 10; got != want {
		t.Fatalf("unexpected surviving frame returnPC: got %d, want %d", got, want)
	}

	if got, want := instance.registers.Values[1], runtime.None; got != want {
		t.Fatalf("expected protected return destination to be reset, got %v", got)
	}

	reused4 := instance.frames.AcquireRegisters(4)
	if len(reused4) != 4 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused4), 4)
	}
	if &reused4[0] != &aboveRegs1[0] {
		t.Fatal("expected frame registers of size 4 to be reclaimed")
	}

	reused5 := instance.frames.AcquireRegisters(5)
	if len(reused5) != 5 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused5), 5)
	}
	if &reused5[0] != &aboveRegs2[0] {
		t.Fatal("expected frame registers of size 5 to be reclaimed")
	}

	reused6 := instance.frames.AcquireRegisters(6)
	if len(reused6) != 6 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused6), 6)
	}
	if &reused6[0] != &activeRegs[0] {
		t.Fatal("expected active registers of size 6 to be reclaimed")
	}
}

func TestRunReturnsUnresolvedFunctionWhenHostCacheEntryIsMissing(t *testing.T) {
	c := compiler.New()
	program, err := c.Compile(file.NewSource("test", "RETURN TEST(1)"))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	instance := New(program)
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.True, nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	if _, err := instance.Run(context.Background(), env); err != nil {
		t.Fatalf("first run failed: %v", err)
	}

	hostPC := -1
	for i, inst := range program.Bytecode {
		if inst.Opcode == bytecode.OpHCall || inst.Opcode == bytecode.OpProtectedHCall {
			hostPC = i
			break
		}
	}

	if hostPC < 0 {
		t.Fatal("host call opcode not found")
	}

	instance.cache.HostFunctions[hostPC] = nil

	_, err = instance.Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected unresolved function error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if rtErr.Message != "Unresolved function" {
		t.Fatalf("expected unresolved function message, got %q", rtErr.Message)
	}
}

func TestSetCallResult_AppliesCatchJumpZeroAndFallbackValue(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 0},
		},
	})

	instance.pc = 1
	instance.registers.Values[1] = runtime.True

	err := instance.setCallResult(
		bytecode.OpHCall,
		bytecode.NewRegister(1),
		runtime.True,
		errors.New("boom"),
	)

	if err != nil {
		t.Fatalf("expected caught error to be swallowed, got %v", err)
	}

	if got := instance.registers.Values[1]; got != runtime.None {
		t.Fatalf("expected destination to be reset to none, got %v", got)
	}

	if got, want := instance.pc, 0; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_AppliesJumpTargetZero(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 0},
		},
	})

	instance.pc = 1
	called := false

	err := instance.handleErrorWithCatch(errors.New("boom"), func() {
		called = true
	})

	if err != nil {
		t.Fatalf("expected caught error to be swallowed, got %v", err)
	}

	if !called {
		t.Fatal("expected onCatch callback to be called")
	}

	if got, want := instance.pc, 0; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_AppliesPositiveJumpTarget(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 2},
		},
	})

	instance.pc = 1

	err := instance.handleErrorWithCatch(errors.New("boom"), nil)
	if err != nil {
		t.Fatalf("expected caught error to be swallowed, got %v", err)
	}

	if got, want := instance.pc, 2; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_ReturnsErrorOutsideCatchRegion(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{0, 0, 1},
		},
	})

	instance.pc = 1
	called := false
	wantErr := errors.New("boom")

	err := instance.handleErrorWithCatch(wantErr, func() {
		called = true
	})

	if err != wantErr {
		t.Fatalf("expected original error to be returned, got %v", err)
	}

	if called {
		t.Fatal("expected onCatch callback not to be called")
	}

	if got, want := instance.pc, 1; got != want {
		t.Fatalf("expected pc to stay unchanged at %d, got %d", want, got)
	}
}

func TestOpFail_UncaughtReturnsRuntimeError(t *testing.T) {
	instance := New(&bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpFail, bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewString("boom"),
		},
	})

	_, err := instance.Run(context.Background(), nil)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if !strings.Contains(rtErr.Format(), "boom") {
		t.Fatalf("expected runtime error to include fail message, got:\n%s", rtErr.Format())
	}
}

func TestOpFail_CaughtUsesCatchJumpTarget(t *testing.T) {
	instance := New(&bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpFail, bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewString("boom"),
			runtime.NewInt(7),
		},
		CatchTable: []bytecode.Catch{
			{2, 2, 3},
		},
	})

	result, err := instance.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected fail to be caught, got %v", err)
	}

	if result != runtime.NewInt(7) {
		t.Fatalf("expected catch jump target to continue execution, got %v", result)
	}
}

func TestOpFail_InvalidMessageTypeReturnsTypeError(t *testing.T) {
	instance := New(&bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpFail, bytecode.NewConstant(0)),
		},
		Constants: []runtime.Value{
			runtime.NewInt(1),
		},
	})

	_, err := instance.Run(context.Background(), nil)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if !strings.Contains(strings.ToLower(rtErr.Format()), "invalid type") {
		t.Fatalf("expected invalid type error, got:\n%s", rtErr.Format())
	}
}
