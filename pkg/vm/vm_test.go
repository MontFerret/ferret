package vm

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
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

func mustNewVM(t *testing.T, program *bytecode.Program, opts ...Option) *VM {
	t.Helper()

	instance, err := NewWith(program, opts...)
	if err != nil {
		t.Fatalf("vm init failed: %v", err)
	}

	return instance
}

func mustAcquireRunState(t *testing.T, instance *VM) *execState {
	t.Helper()

	state := instance.acquireRunState()
	if state == nil {
		t.Fatal("expected run state")
	}

	return state
}

func TestPanicPolicyRecoversPanics(t *testing.T) {
	program := compileProgram(t, "RETURN PANIC_FN()")

	instance := mustNewVM(t, program)
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

	instance := mustNewVM(t, program, WithPanicPolicy(PanicPropagate))
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

	instance := mustNewVM(t, program, WithPanicPolicy(PanicPropagate))
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
		ISAVersion: bytecode.Version,
		Registers:  6,
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

	instance := mustNewVM(t,
		program,
		WithShapeCacheLimit(17),
		WithFastObjectDictThreshold(23),
		WithPanicPolicy(PanicPropagate),
	)

	if instance.program != program {
		t.Fatal("expected VM to keep source program reference")
	}

	state := mustAcquireRunState(t, instance)
	defer instance.releaseRunState(state)

	if state.registers == nil {
		t.Fatal("expected register file to be initialized")
	}

	if got, want := state.registers.Size(), program.Registers; got != want {
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

	reg := state.frames.AcquireRegisters(5)
	if got, want := len(reg), 5; got != want {
		t.Fatalf("unexpected pooled register size: got %d, want %d", got, want)
	}

	state.frames.ReleaseRegisters(reg)
	reused := state.frames.AcquireRegisters(5)
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

	out, err := mustNewVM(t, program).Run(context.Background(), env)
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

	_, err := mustNewVM(t, program).Run(context.Background(), env)
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

func TestRun_MissingParamPrecedesWarmupHostResolution(t *testing.T) {
	program := compileProgram(t, "RETURN MISSING_FN(@foo)")

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected missing parameter error")
	}

	if !strings.Contains(err.Error(), "Missing parameter") {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(err.Error(), "Unresolved function") {
		t.Fatalf("expected missing parameter to be reported before host warmup failure, got %v", err)
	}

	cause := errors.Unwrap(err)
	if cause == nil {
		cause = err
	}

	if !strings.Contains(cause.Error(), "@foo") {
		t.Fatalf("expected missing parameter name in error, got %v (cause: %v)", err, cause)
	}
}

func TestUnwindToProtected_ReclaimsDiscardedFrameRegisters(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
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

	state := mustAcquireRunState(t, instance)
	defer instance.releaseRunState(state)

	state.registers.Values = activeRegs
	state.frames.Push(frame.CallFrame{
		ReturnPC:   10,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  lowerRegs,
		Protected:  false,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:   20,
		ReturnDest: bytecode.NewRegister(1),
		Registers:  protectedRegs,
		Protected:  true,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:   30,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  aboveRegs1,
		Protected:  false,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:   40,
		ReturnDest: bytecode.NewRegister(0),
		Registers:  aboveRegs2,
		Protected:  false,
	})

	if ok := state.unwindToProtected(); !ok {
		t.Fatal("expected protected unwind to succeed")
	}

	if got, want := state.pc, 20; got != want {
		t.Fatalf("unexpected pc after unwind: got %d, want %d", got, want)
	}

	if got, want := state.frames.Len(), 1; got != want {
		t.Fatalf("unexpected frame depth after unwind: got %d, want %d", got, want)
	}

	remaining := state.frames.Top()
	if remaining == nil {
		t.Fatal("expected remaining frame after unwind")
	}

	if got, want := remaining.ReturnPC, 10; got != want {
		t.Fatalf("unexpected surviving frame returnPC: got %d, want %d", got, want)
	}

	if got, want := state.registers.Values[1], runtime.None; got != want {
		t.Fatalf("expected protected return destination to be reset, got %v", got)
	}

	reused4 := state.frames.AcquireRegisters(4)
	if len(reused4) != 4 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused4), 4)
	}
	if &reused4[0] != &aboveRegs1[0] {
		t.Fatal("expected frame registers of size 4 to be reclaimed")
	}

	reused5 := state.frames.AcquireRegisters(5)
	if len(reused5) != 5 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused5), 5)
	}
	if &reused5[0] != &aboveRegs2[0] {
		t.Fatal("expected frame registers of size 5 to be reclaimed")
	}

	reused6 := state.frames.AcquireRegisters(6)
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

	instance := mustNewVM(t, program)
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
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 0},
		},
	})

	state := mustAcquireRunState(t, instance)
	defer instance.releaseRunState(state)

	state.pc = 1
	state.registers.Values[1] = runtime.True

	action := state.setCallResult(
		bytecode.OpHCall,
		bytecode.NewRegister(1),
		runtime.True,
		errors.New("boom"),
	)

	if action == errReturn {
		t.Fatalf("expected caught error to be swallowed, got %v", action)
	}

	if got := state.registers.Values[1]; got != runtime.None {
		t.Fatalf("expected destination to be reset to none, got %v", got)
	}

	if got, want := state.pc, 0; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_AppliesJumpTargetZero(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 0},
		},
	})

	state := mustAcquireRunState(t, instance)
	defer instance.releaseRunState(state)

	state.pc = 1
	action := state.applyCatch(bytecode.Operand(1), runtime.True, errors.New("boom"))
	if action == errReturn {
		t.Fatalf("expected caught error to be swallowed, got %v", action)
	}

	val := state.registers.Values[1]
	if val != runtime.True {
		t.Fatalf("expected fallback value to be set in destination register, got %v", val)
	}

	if got, want := state.pc, 0; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_AppliesPositiveJumpTarget(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 2},
		},
	})

	state := mustAcquireRunState(t, instance)
	defer instance.releaseRunState(state)

	state.pc = 1
	action := state.applyCatch(bytecode.NoopOperand, nil, errors.New("boom"))
	if action == errReturn {
		t.Fatalf("expected caught error to be swallowed, got %v", action)
	}

	if got, want := state.pc, 2; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_ReturnsErrorOutsideCatchRegion(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{0, 0, 1},
		},
	})

	state := mustAcquireRunState(t, instance)
	defer instance.releaseRunState(state)

	state.pc = 1
	wantErr := errors.New("boom")

	action := state.applyCatch(bytecode.Operand(1), runtime.True, wantErr)
	if action != errReturn {
		t.Fatalf("expected original error to be returned, got %v", action)
	}

	val := state.registers.Values[1]
	if val == runtime.True {
		t.Fatalf("expected fallback value to be ignored, got %v", val)
	}

	if got, want := state.pc, 1; got != want {
		t.Fatalf("expected pc to stay unchanged at %d, got %d", want, got)
	}
}

func TestOpFail_UncaughtReturnsRuntimeError(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
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
	instance := mustNewVM(t, &bytecode.Program{
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
	instance := mustNewVM(t, &bytecode.Program{
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

func TestWarmupClearsStaleHostCacheAcrossEnvironments(t *testing.T) {
	program := compileProgram(t, "RETURN F()")
	instance := mustNewVM(t, program)

	envWithFn, err := NewEnvironment([]EnvironmentOption{
		WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(7), nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	out, err := instance.Run(context.Background(), envWithFn)
	if err != nil {
		t.Fatalf("first run failed: %v", err)
	}

	if got, want := out, runtime.NewInt(7); got != want {
		t.Fatalf("unexpected first run result: got %v, want %v", got, want)
	}

	_, err = instance.Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected unresolved function after env switch")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if got, want := rtErr.Message, "Unresolved function"; got != want {
		t.Fatalf("unexpected message: got %q, want %q", got, want)
	}
}

func TestWarmupRebindTouchesOnlyHostCallSlots(t *testing.T) {
	program := compileProgram(t, "RETURN F()")
	instance := mustNewVM(t, program)

	hostPC := -1
	nonHostPC := -1
	for i, inst := range program.Bytecode {
		if inst.Opcode == bytecode.OpHCall || inst.Opcode == bytecode.OpProtectedHCall {
			if hostPC < 0 {
				hostPC = i
			}
			continue
		}

		if nonHostPC < 0 {
			nonHostPC = i
		}
	}

	if hostPC < 0 || nonHostPC < 0 {
		t.Fatalf("expected host and non-host opcodes, got host=%d nonHost=%d", hostPC, nonHostPC)
	}

	sentinel := &mem.CachedHostFunction{
		FnV: func(_ context.Context, _ ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(77), nil
		},
	}
	instance.cache.HostFunctions[nonHostPC] = sentinel

	envA, err := NewEnvironment([]EnvironmentOption{
		WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(1), nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	envB, err := NewEnvironment([]EnvironmentOption{
		WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(2), nil
		}),
		WithFunction("G", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(3), nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	if _, err := instance.Run(context.Background(), envA); err != nil {
		t.Fatalf("first run failed: %v", err)
	}

	instance.cache.HostFunctions[nonHostPC] = sentinel
	if _, err := instance.Run(context.Background(), envB); err != nil {
		t.Fatalf("second run failed: %v", err)
	}

	if got := instance.cache.HostFunctions[nonHostPC]; got != sentinel {
		t.Fatalf("unexpected non-host cache mutation: got %p, want %p", got, sentinel)
	}

	if instance.cache.HostFunctions[hostPC] == nil {
		t.Fatal("expected host call slot to be rebound")
	}
}

func TestStrictWarmupFailsProtectedMissingHostCallForDefaultAndBuiltEnvironment(t *testing.T) {
	program := compileProgram(t, "RETURN MISSING_FN()?")

	builtEnv, err := NewEnvironment(nil)
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	cases := []struct {
		env  *Environment
		name string
	}{
		{name: "default", env: NewDefaultEnvironment()},
		{name: "built", env: builtEnv},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mustNewVM(t, program).Run(context.Background(), tc.env)
			if err == nil {
				t.Fatal("expected unresolved function error from strict warmup")
			}

			var rtErr *RuntimeError
			if !errors.As(err, &rtErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}

			if got, want := rtErr.Message, "Unresolved function"; got != want {
				t.Fatalf("unexpected message: got %q, want %q", got, want)
			}
		})
	}
}

func TestStrictWarmupFailsOnDeadCodeUnresolvedHostCall(t *testing.T) {
	program := compileProgram(t, "RETURN false ? MISSING_FN() : 1")

	envWithDummy, err := NewEnvironment([]EnvironmentOption{
		WithFunction("DUMMY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	cases := []struct {
		env  *Environment
		name string
	}{
		{name: "default", env: NewDefaultEnvironment()},
		{name: "dummy", env: envWithDummy},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mustNewVM(t, program).Run(context.Background(), tc.env)
			if err == nil {
				t.Fatal("expected unresolved function error from strict warmup")
			}

			var rtErr *RuntimeError
			if !errors.As(err, &rtErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}

			if got, want := rtErr.Message, "Unresolved function"; got != want {
				t.Fatalf("unexpected message: got %q, want %q", got, want)
			}
		})
	}
}

func TestStrictWarmupAggregatesMissingHostFunctions(t *testing.T) {
	program := compileProgram(t, `
LET a = MISSING_A()
LET b = MISSING_B()
RETURN a + b
`)

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected warmup error set")
	}

	var rtErrSet *RuntimeErrorSet
	if !errors.As(err, &rtErrSet) {
		t.Fatalf("expected runtime error set, got %T", err)
	}

	if got, want := rtErrSet.Size(), 2; got != want {
		t.Fatalf("unexpected error set size: got %d, want %d", got, want)
	}

	unresolved := 0
	for _, rtErr := range rtErrSet.Errors() {
		if rtErr.Message == "Unresolved function" {
			unresolved++
		}
	}

	if got, want := unresolved, 2; got != want {
		t.Fatalf("unexpected unresolved error count: got %d, want %d", got, want)
	}
}

func TestStrictWarmupInvalidTargetReturnsInvalidFunctionName(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
	}

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if got, want := rtErr.Message, "Invalid function name"; got != want {
		t.Fatalf("unexpected message: got %q, want %q", got, want)
	}
}

func TestStrictWarmupFailureIsRepeatableUntilEnvironmentFixed(t *testing.T) {
	program := compileProgram(t, "RETURN F()")
	instance := mustNewVM(t, program)

	assertUnresolved := func(err error) {
		t.Helper()

		if err == nil {
			t.Fatal("expected unresolved function error")
		}

		var rtErr *RuntimeError
		if !errors.As(err, &rtErr) {
			t.Fatalf("expected runtime error, got %T", err)
		}

		if got, want := rtErr.Message, "Unresolved function"; got != want {
			t.Fatalf("unexpected message: got %q, want %q", got, want)
		}
	}

	_, err := instance.Run(context.Background(), NewDefaultEnvironment())
	assertUnresolved(err)

	_, err = instance.Run(context.Background(), NewDefaultEnvironment())
	assertUnresolved(err)

	validEnv, err := NewEnvironment([]EnvironmentOption{
		WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(7), nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	out, err := instance.Run(context.Background(), validEnv)
	if err != nil {
		t.Fatalf("expected successful run after env fix, got %v", err)
	}

	if got, want := out, runtime.NewInt(7); got != want {
		t.Fatalf("unexpected output: got %v, want %v", got, want)
	}
}

func TestResetDrainsLeakedFramesBetweenFailedRuns(t *testing.T) {
	program := compileProgram(t, `
FUNC inner() (
	RETURN 1 / 0
)

FUNC outer() (
	RETURN inner()
)

RETURN outer()
`)

	instance := mustNewVM(t, program)

	runAndCheck := func(label string) int {
		t.Helper()

		_, err := instance.Run(context.Background(), NewDefaultEnvironment())
		if err == nil {
			t.Fatalf("%s: expected runtime error", label)
		}

		var rtErr *RuntimeError
		if !errors.As(err, &rtErr) {
			t.Fatalf("%s: expected runtime error, got %T", label, err)
		}

		if rtErr.Kind != DivideByZero {
			t.Fatalf("%s: unexpected error kind: got %s, want %s", label, rtErr.Kind, DivideByZero)
		}

		return strings.Count(rtErr.Format(), "called from")
	}

	stackDepthFirst := runAndCheck("first run")
	stackDepthSecond := runAndCheck("second run")
	if stackDepthSecond != stackDepthFirst {
		t.Fatalf("expected stable stack depth across repeated failed runs: first=%d second=%d", stackDepthFirst, stackDepthSecond)
	}
}

func TestRunReenterSameVMUsesIsolatedRunState(t *testing.T) {
	program := compileProgram(t, "RETURN REENTER()")
	instance := mustNewVM(t, program)

	var (
		env   *Environment
		depth int
		err   error
	)

	env, err = NewEnvironment([]EnvironmentOption{
		WithFunction("REENTER", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			if depth > 0 {
				return runtime.NewInt(42), nil
			}

			depth++
			defer func() {
				depth--
			}()

			return instance.Run(ctx, env)
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	out, err := instance.Run(context.Background(), env)
	if err != nil {
		t.Fatalf("outer run failed: %v", err)
	}

	if got, want := out, runtime.NewInt(42); got != want {
		t.Fatalf("unexpected first output: got %v, want %v", got, want)
	}

	out, err = instance.Run(context.Background(), env)
	if err != nil {
		t.Fatalf("second outer run failed: %v", err)
	}

	if got, want := out, runtime.NewInt(42); got != want {
		t.Fatalf("unexpected second output: got %v, want %v", got, want)
	}
}

func TestHostNilResultIsNormalizedToNone(t *testing.T) {
	program := compileProgram(t, "RETURN NIL_FN()")
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("NIL_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return nil, nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	out, err := mustNewVM(t, program).Run(context.Background(), env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != runtime.None {
		t.Fatalf("expected none, got %v", out)
	}
}

func TestModuloTypeErrorNotMisclassifiedAsModuloByZero(t *testing.T) {
	program := compileProgram(t, `RETURN 5 % "x"`)

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if rtErr.Kind != diagnostics.TypeError {
		t.Fatalf("unexpected kind: got %s, want %s", rtErr.Kind, diagnostics.TypeError)
	}

	formatted := strings.ToLower(rtErr.Format())
	if strings.Contains(formatted, "modulo by zero") {
		t.Fatalf("expected non-modulo classification, got:\n%s", rtErr.Format())
	}
}

func TestInvalidFunctionNameHasNonEmptyKind(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{runtime.NewInt(123)},
	}

	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("DUMMY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	_, err = mustNewVM(t, program).Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if rtErr.Kind == "" {
		t.Fatalf("expected non-empty diagnostic kind, got %q", rtErr.Kind)
	}

	if got, want := rtErr.Message, "Invalid function name"; got != want {
		t.Fatalf("unexpected message: got %q, want %q", got, want)
	}
}

func TestWrapRuntimeErrorSingleWarmupFailureReturnsRuntimeError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     file.NewSource("test", "RETURN 1"),
		Metadata: bytecode.Metadata{
			DebugSpans: []file.Span{{Start: 0, End: 6}},
		},
	}

	warmup := &diagnostic.WarmupErrorSet{}
	warmup.Add(ErrInvalidFunctionName, 0, bytecode.NewRegister(0))

	err := diagnostic.WrapRuntimeError(program, 1, nil, warmup)
	if err == nil {
		t.Fatal("expected wrapped error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	var rtErrSet *RuntimeErrorSet
	if errors.As(err, &rtErrSet) {
		t.Fatalf("expected single runtime error, got set")
	}
}

func TestNearestBoundaryPrefersCatchOverProtectedUnwind(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 3},
		},
	})

	protectedRegs := make([]runtime.Value, 2)
	state := mustAcquireRunState(t, instance)
	defer instance.releaseRunState(state)

	state.frames.Push(frame.CallFrame{
		ReturnPC:   9,
		ReturnDest: bytecode.NewRegister(1),
		Registers:  protectedRegs,
		Protected:  true,
	})
	state.pc = 1

	action := state.applyProtected(errors.New("boom"))
	if action != errContinue {
		t.Fatalf("expected continue, got %v", action)
	}

	if got, want := state.pc, 3; got != want {
		t.Fatalf("expected catch jump to win, got %d", got)
	}

	if got, want := state.frames.Len(), 1; got != want {
		t.Fatalf("expected protected frame to remain, got %d", got)
	}
}

func TestNearestBoundaryUsesProtectedUnwindWithoutCatch(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
	})

	lowerRegs := make([]runtime.Value, 2)
	activeRegs := make([]runtime.Value, 2)
	state := mustAcquireRunState(t, instance)
	defer instance.releaseRunState(state)

	state.registers.Values = activeRegs
	state.frames.Push(frame.CallFrame{
		ReturnPC:   7,
		ReturnDest: bytecode.NewRegister(1),
		Registers:  lowerRegs,
		Protected:  true,
	})
	state.pc = 1

	action := state.applyProtected(errors.New("boom"))
	if action != errContinue {
		t.Fatalf("expected continue, got %v", action)
	}

	if got, want := state.pc, 7; got != want {
		t.Fatalf("unexpected unwind target: got %d, want %d", got, want)
	}

	if got, want := state.frames.Len(), 0; got != want {
		t.Fatalf("expected stack to be unwound, got len %d", got)
	}
}

func TestRuntimeErrorIncludesUDFCallStackContext(t *testing.T) {
	program := compileProgram(t, `
FUNC inner() (
	RETURN @x.foo
)
FUNC middle() (
	LET value = inner()
	RETURN value
)
FUNC outer() (
	LET value = middle()
	RETURN value
)
RETURN outer()
`)

	env := NewDefaultEnvironment()
	env.Params["x"] = runtime.None

	_, err := mustNewVM(t, program).Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	formatted := rtErr.Format()
	if !strings.Contains(formatted, "called from inner (#1)") {
		t.Fatalf("expected VM call stack context, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "VM stack: outer -> middle -> inner") {
		t.Fatalf("expected additive VM stack note, got:\n%s", formatted)
	}
}

func TestWrapRuntimeErrorPreservesExistingNoteAndDoesNotDuplicateStackDetails(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     file.NewSource("test", "RETURN 1"),
		Metadata: bytecode.Metadata{
			DebugSpans: []file.Span{{Start: 0, End: 6}},
		},
	}

	base := diagnostic.NewRuntimeError(
		program,
		1,
		diagnostic.UncaughtError,
		"Runtime error",
		"",
		"",
		"existing note",
	)

	stack := []frame.TraceEntry{
		{CallSitePC: 0, FnID: 2, FnName: "inner"},
		{CallSitePC: 0, FnID: 1, FnName: "outer"},
	}

	wrapped := diagnostic.WrapRuntimeError(program, 1, stack, base)

	var rtErr *RuntimeError
	if !errors.As(wrapped, &rtErr) {
		t.Fatalf("expected runtime error, got %T", wrapped)
	}

	if !strings.Contains(rtErr.Note, "existing note") {
		t.Fatalf("expected existing note to be preserved, got %q", rtErr.Note)
	}

	if !strings.Contains(rtErr.Note, "VM stack: outer -> inner") {
		t.Fatalf("expected VM stack note, got %q", rtErr.Note)
	}

	wrappedTwice := diagnostic.WrapRuntimeError(program, 1, stack, wrapped)
	if !errors.As(wrappedTwice, &rtErr) {
		t.Fatalf("expected runtime error, got %T", wrappedTwice)
	}

	if strings.Count(rtErr.Note, "VM stack:") != 1 {
		t.Fatalf("expected idempotent stack note append, got %q", rtErr.Note)
	}

	if strings.Count(rtErr.Format(), "called from inner (#1)") != 1 {
		t.Fatalf("expected idempotent stack spans, got:\n%s", rtErr.Format())
	}
}

func TestWrapRuntimeErrorRecognizesLegacyCallStackLabelWithoutDuplicatingSpans(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     file.NewSource("test", "RETURN 1"),
		Metadata: bytecode.Metadata{
			DebugSpans: []file.Span{{Start: 0, End: 6}},
		},
	}

	base := &diagnostic.RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    diagnostic.UncaughtError,
			Message: "Runtime error",
			Source:  program.Source,
			Spans: []diagnostics.ErrorSpan{
				diagnostics.NewSecondaryErrorSpan(file.Span{Start: 0, End: 6}, "called from (#1)"),
				diagnostics.NewMainErrorSpan(file.Span{Start: 0, End: 6}, ""),
			},
		},
	}

	stack := []frame.TraceEntry{
		{CallSitePC: 0, FnID: 1, FnName: "boo"},
	}

	wrapped := diagnostic.WrapRuntimeError(program, 1, stack, base)

	var rtErr *RuntimeError
	if !errors.As(wrapped, &rtErr) {
		t.Fatalf("expected runtime error, got %T", wrapped)
	}

	formatted := rtErr.Format()
	if strings.Count(formatted, "called from (#1)") != 1 {
		t.Fatalf("expected legacy stack label to be preserved without duplication, got:\n%s", formatted)
	}
}

func TestRuntimeErrorSingleUdfStackFormattingUsesSourceSpelling(t *testing.T) {
	program := compileProgram(t, `
FUNC boo() (
	LET a = 1
	LET b = 0
	RETURN a / b
)
RETURN boo()
`)

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	formatted := rtErr.Format()
	if !strings.Contains(formatted, "called from boo (#1)") {
		t.Fatalf("expected call-site label with source-spelling udf name, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "VM stack: boo") {
		t.Fatalf("expected source-spelling VM stack note, got:\n%s", formatted)
	}
}

func TestUdfRuntimeMessageUsesSourceSpellingName(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{
			runtime.NewInt(0),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{
					Name:        "BOO",
					DisplayName: "boo",
					Entry:       3,
					Registers:   2,
					Params:      1,
				},
			},
		},
	}

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	formatted := rtErr.Format()
	if !strings.Contains(formatted, "'boo'") {
		t.Fatalf("expected source-spelling name in runtime diagnostic, got:\n%s", formatted)
	}

	if strings.Contains(formatted, "'BOO'") {
		t.Fatalf("unexpected normalized name in runtime diagnostic:\n%s", formatted)
	}
}

func TestRecoveredPanicRuntimeErrorDoesNotLeakGoStackTrace(t *testing.T) {
	program := compileProgram(t, "RETURN PANIC_FN()")
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("PANIC_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			panic("panic in host function")
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	_, err = mustNewVM(t, program).Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	formatted := rtErr.Format()
	if strings.Contains(formatted, "goroutine ") || strings.Contains(formatted, "runtime/debug.Stack") {
		t.Fatalf("runtime error format leaked Go stack trace:\n%s", formatted)
	}
}

func TestInvariantInRecoverModeBecomesUnexpectedRuntimeError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(
				bytecode.OpDataSetCollector,
				bytecode.NewRegister(1),
				bytecode.Operand(bytecode.CollectorTypeAggregate),
				bytecode.Operand(99),
			),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if rtErr.Kind != diagnostics.UnexpectedError {
		t.Fatalf("unexpected kind: got %s, want %s", rtErr.Kind, diagnostics.UnexpectedError)
	}

	if !strings.Contains(strings.ToLower(rtErr.Message), "invalid aggregate plan index") {
		t.Fatalf("unexpected message: %q", rtErr.Message)
	}
}

func TestInvalidCollectorTypeInvariantInRecoverModeBecomesUnexpectedRuntimeError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(
				bytecode.OpDataSetCollector,
				bytecode.NewRegister(1),
				bytecode.Operand(255),
				bytecode.Operand(0),
			),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if rtErr.Kind != diagnostics.UnexpectedError {
		t.Fatalf("unexpected kind: got %s, want %s", rtErr.Kind, diagnostics.UnexpectedError)
	}

	if !strings.Contains(strings.ToLower(rtErr.Message), "invalid collector configuration") {
		t.Fatalf("unexpected message: %q", rtErr.Message)
	}
}

func TestInvariantInPropagateModePanics(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(
				bytecode.OpDataSetCollector,
				bytecode.NewRegister(1),
				bytecode.Operand(bytecode.CollectorTypeAggregate),
				bytecode.Operand(99),
			),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	instance := mustNewVM(t, program, WithPanicPolicy(PanicPropagate))

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for invariant in propagate mode")
		}
	}()

	_, _ = instance.Run(context.Background(), NewDefaultEnvironment())
}

func TestInvalidCollectorTypeInvariantInPropagateModePanics(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(
				bytecode.OpDataSetCollector,
				bytecode.NewRegister(1),
				bytecode.Operand(255),
				bytecode.Operand(0),
			),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	instance := mustNewVM(t, program, WithPanicPolicy(PanicPropagate))

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for invariant in propagate mode")
		}
	}()

	_, _ = instance.Run(context.Background(), NewDefaultEnvironment())
}
