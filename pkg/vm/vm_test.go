package vm

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	rtdiagnostics "github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

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

	state := &instance.state
	if state == nil {
		t.Fatal("expected run state")
	}

	return state
}

func mustNewEnvironment(t *testing.T, opts ...EnvironmentOption) *Environment {
	t.Helper()

	env, err := NewEnvironment(opts)
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	return env
}

func assertUnresolvedFunctionError(t *testing.T, err error) *RuntimeError {
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

	return rtErr
}

func assertIntArrayValue(t *testing.T, got runtime.Value, want ...runtime.Int) {
	t.Helper()

	arr, ok := got.(*runtime.Array)
	if !ok {
		t.Fatalf("expected runtime.Array, got %T", got)
	}

	expected := make([]runtime.Value, len(want))
	for i, value := range want {
		expected[i] = value
	}

	if arr.Compare(runtime.NewArrayWith(expected...)) != 0 {
		t.Fatalf("unexpected array value: got %v, want %v", got, runtime.NewArrayWith(expected...))
	}
}

type trackingCloser struct {
	name   string
	closed int
}

func newTrackingCloser(name string) *trackingCloser {
	return &trackingCloser{name: name}
}

func (c *trackingCloser) Close() error {
	c.closed++
	return nil
}

func (c *trackingCloser) String() string {
	return c.name
}

func (c *trackingCloser) Hash() uint64 {
	return 0
}

func (c *trackingCloser) Copy() runtime.Value {
	return c
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
	defer state.endRun()

	if got, want := len(state.registers), program.Registers; got != want {
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
	if got, want := len(instance.cache.HostFunctions), len(instance.plan.hostCallDescriptors); got != want {
		t.Fatalf("unexpected host function cache size: got %d, want %d", got, want)
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

	if got, want := len(instance.plan.instructions), bytecodeLen; got != want {
		t.Fatalf("unexpected instruction wrapper size: got %d, want %d", got, want)
	}

	for i := range program.Bytecode {
		if got, want := instance.plan.instructions[i].Instruction, program.Bytecode[i]; got != want {
			t.Fatalf("unexpected wrapped instruction at %d: got %+v, want %+v", i, got, want)
		}
	}

	wantCatchByPC := []int{0, 0}
	for i := range wantCatchByPC {
		if got := state.catchByPC[i]; got != wantCatchByPC[i] {
			t.Fatalf("unexpected catch mapping at pc %d: got %d, want %d", i, got, wantCatchByPC[i])
		}
	}

	reg := state.windows.Acquire(5)
	if got, want := len(reg), 5; got != want {
		t.Fatalf("unexpected pooled register size: got %d, want %d", got, want)
	}

	state.windows.Release(reg)
	reused := state.windows.Acquire(5)
	if got, want := len(reused), 5; got != want {
		t.Fatalf("unexpected reused register size: got %d, want %d", got, want)
	}

	if &reg[0] != &reused[0] {
		t.Fatal("expected register pool to reuse buffers")
	}
}

func TestVMCloseIsIdempotent(t *testing.T) {
	instance := mustNewVM(t, newTestProgram(
		1,
		nil,
		bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	))

	if err := instance.Close(); err != nil {
		t.Fatalf("expected first close to succeed, got: %v", err)
	}

	if err := instance.Close(); err != nil {
		t.Fatalf("expected repeated close to succeed, got: %v", err)
	}

	if !instance.closed {
		t.Fatal("expected VM to be marked closed")
	}

	if instance.cache != nil {
		t.Fatal("expected VM close to clear cache reference")
	}

	if instance.program != nil {
		t.Fatal("expected VM close to clear program reference")
	}

	if len(instance.plan.instructions) != 0 {
		t.Fatal("expected VM close to clear exec plan")
	}
}

func TestVMRunAfterCloseReturnsInvalidOperation(t *testing.T) {
	instance := mustNewVM(t, newTestProgram(
		1,
		nil,
		bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	))

	if err := instance.Close(); err != nil {
		t.Fatalf("expected close to succeed, got: %v", err)
	}

	_, err := instance.Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected run on closed VM to fail")
	}

	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected invalid operation for run after close, got: %v", err)
	}

	if !strings.Contains(err.Error(), "vm is closed") {
		t.Fatalf("expected closed VM message, got: %v", err)
	}
}

func TestVMCloseClosesTrackedOwnedResources(t *testing.T) {
	instance := mustNewVM(t, newTestProgram(
		1,
		nil,
		bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	))

	closer := newTrackingCloser("owned")
	instance.state.owned.Track(closer)

	if err := instance.Close(); err != nil {
		t.Fatalf("expected close to succeed, got: %v", err)
	}

	if got, want := closer.closed, 1; got != want {
		t.Fatalf("expected close to release owned resources once, got %d", got)
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
	defer state.endRun()

	state.registers = activeRegs
	state.frames.Push(frame.CallFrame{
		ReturnPC:         10,
		ReturnDest:       bytecode.NewRegister(0),
		CallerRegisters:  lowerRegs,
		RecoveryBoundary: false,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:         20,
		ReturnDest:       bytecode.NewRegister(1),
		CallerRegisters:  protectedRegs,
		RecoveryBoundary: true,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:         30,
		ReturnDest:       bytecode.NewRegister(0),
		CallerRegisters:  aboveRegs1,
		RecoveryBoundary: false,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:         40,
		ReturnDest:       bytecode.NewRegister(0),
		CallerRegisters:  aboveRegs2,
		RecoveryBoundary: false,
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

	if &state.registers[0] != &protectedRegs[0] {
		t.Fatal("expected unwind to restore the protected caller register window")
	}

	remaining := state.frames.Top()
	if remaining == nil {
		t.Fatal("expected remaining frame after unwind")
	}

	if got, want := remaining.ReturnPC, 10; got != want {
		t.Fatalf("unexpected surviving frame returnPC: got %d, want %d", got, want)
	}

	if got, want := state.registers[1], runtime.None; got != want {
		t.Fatalf("expected protected return destination to be reset, got %v", got)
	}

	reused4 := state.windows.Acquire(4)
	if len(reused4) != 4 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused4), 4)
	}
	if &reused4[0] != &aboveRegs1[0] {
		t.Fatal("expected frame registers of size 4 to be reclaimed")
	}

	reused5 := state.windows.Acquire(5)
	if len(reused5) != 5 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused5), 5)
	}
	if &reused5[0] != &aboveRegs2[0] {
		t.Fatal("expected frame registers of size 5 to be reclaimed")
	}

	reused6 := state.windows.Acquire(6)
	if len(reused6) != 6 {
		t.Fatalf("unexpected pooled registers length: got %d, want %d", len(reused6), 6)
	}
	if &reused6[0] != &activeRegs[0] {
		t.Fatal("expected active registers of size 6 to be reclaimed")
	}
}

func TestUnwindToProtected_ClosesOwnedResourcesOfUnwoundFramesOnly(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Registers: 6},
			},
		},
	})

	protectedOwned := newTrackingCloser("protected")
	aboveOwned1 := newTrackingCloser("above-1")
	aboveOwned2 := newTrackingCloser("above-2")
	activeOwned := newTrackingCloser("active")

	lowerRegs := mem.NewRegisterFile(2)
	protectedRegs := mem.NewRegisterFile(3)
	aboveRegs1 := mem.NewRegisterFile(4)
	aboveRegs2 := mem.NewRegisterFile(5)
	activeRegs := mem.NewRegisterFile(6)
	protectedRegs[0] = protectedOwned
	aboveRegs1[0] = aboveOwned1
	aboveRegs2[0] = aboveOwned2
	activeRegs[0] = activeOwned

	protectedResources := mem.OwnedResources{}
	protectedResources.Track(protectedOwned)
	aboveResources1 := mem.OwnedResources{}
	aboveResources1.Track(aboveOwned1)
	aboveResources2 := mem.OwnedResources{}
	aboveResources2.Track(aboveOwned2)

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	state.registers = activeRegs
	state.owned.Track(activeOwned)
	state.frames.Push(frame.CallFrame{
		ReturnPC:         10,
		ReturnDest:       bytecode.NewRegister(0),
		CallerRegisters:  lowerRegs,
		RecoveryBoundary: false,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:         20,
		ReturnDest:       bytecode.NewRegister(1),
		CallerRegisters:  protectedRegs,
		OwnedResources:   protectedResources,
		RecoveryBoundary: true,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:        30,
		ReturnDest:      bytecode.NewRegister(0),
		CallerRegisters: aboveRegs1,
		OwnedResources:  aboveResources1,
	})
	state.frames.Push(frame.CallFrame{
		ReturnPC:        40,
		ReturnDest:      bytecode.NewRegister(0),
		CallerRegisters: aboveRegs2,
		OwnedResources:  aboveResources2,
	})

	if ok := state.unwindToProtected(); !ok {
		t.Fatal("expected protected unwind to succeed")
	}

	if got := activeOwned.closed; got != 1 {
		t.Fatalf("expected active frame resource to close once, got %d", got)
	}

	if got := aboveOwned2.closed; got != 1 {
		t.Fatalf("expected top unwound frame resource to close once, got %d", got)
	}

	if got := aboveOwned1.closed; got != 1 {
		t.Fatalf("expected second unwound frame resource to close once, got %d", got)
	}

	if got := protectedOwned.closed; got != 0 {
		t.Fatalf("expected protected caller resource to remain open, got %d closes", got)
	}

	if !state.owned.Owns(protectedOwned) {
		t.Fatal("expected protected caller ownership to remain active after unwind")
	}
}

func TestSetCallResult_RaisesPendingFailureBeforeResolution(t *testing.T) {
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
	defer state.endRun()

	state.pc = 1
	state.registers[1] = runtime.True

	state.setCallResult(
		state.pc,
		bytecode.OpHCall,
		bytecode.NewRegister(1),
		runtime.True,
		errors.New("boom"),
	)

	if !state.hasFailure() {
		t.Fatal("expected pending failure to be recorded")
	}

	if got := state.registers[1]; got != runtime.True {
		t.Fatalf("expected destination to remain unchanged before resolution, got %v", got)
	}

	if got, want := state.pc, 1; got != want {
		t.Fatalf("expected pc to remain unchanged before resolution: got %d, want %d", got, want)
	}

	if action := state.resolveFailure(); action != errContinue {
		t.Fatalf("expected caught error to continue, got %v", action)
	}

	if got := state.registers[1]; got != runtime.None {
		t.Fatalf("expected destination to be reset to none, got %v", got)
	}

	if got, want := state.pc, 0; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}

	if state.hasFailure() {
		t.Fatal("expected pending failure to be cleared after resolution")
	}
}

func TestSetCallResult_ProtectedFailureBypassesCatchJump(t *testing.T) {
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
	defer state.endRun()

	state.pc = 1
	state.registers[1] = runtime.True

	state.setCallResult(
		state.pc,
		bytecode.OpProtectedHCall,
		bytecode.NewRegister(1),
		runtime.True,
		errors.New("boom"),
	)

	if action := state.resolveFailure(); action != errContinue {
		t.Fatalf("expected protected failure to continue, got %v", action)
	}

	if got := state.registers[1]; got != runtime.None {
		t.Fatalf("expected protected fallback to set runtime.None, got %v", got)
	}

	if got, want := state.pc, 1; got != want {
		t.Fatalf("expected protected recovery to bypass catch jump and keep pc %d, got %d", want, got)
	}
}

func TestResolveFailure_InvariantReturnsWithoutPanic(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
	}, WithPanicPolicy(PanicPropagate))

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	invariantErr := rtdiagnostics.NewInvariantError("boom", errors.New("cause"))
	state.raiseInvariant(invariantErr)

	if action := state.resolveFailure(); action != errReturn {
		t.Fatalf("expected invariant resolution to return, got %v", action)
	}

	var gotInvariant *rtdiagnostics.InvariantError
	if got := state.failureError(); !errors.As(got, &gotInvariant) {
		t.Fatalf("expected failure error to remain invariant, got %T", got)
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
	defer state.endRun()

	state.pc = 1
	state.raiseRuntime(errors.New("boom"), recoverDefault, bytecode.Operand(1), runtime.True, true)

	if got := state.registers[1]; got == runtime.True {
		t.Fatalf("expected fallback to be deferred until resolution, got %v", got)
	}

	action := state.resolveFailure()
	if action == errReturn {
		t.Fatalf("expected caught error to be swallowed, got %v", action)
	}

	val := state.registers[1]
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
	defer state.endRun()

	state.pc = 1
	state.raiseRuntime(errors.New("boom"), recoverDefault, bytecode.NoopOperand, nil, false)

	action := state.resolveFailure()
	if action == errReturn {
		t.Fatalf("expected caught error to be swallowed, got %v", action)
	}

	if got, want := state.pc, 2; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}

func TestHandleErrorWithCatch_UsesFailureOriginPC(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 0},
		},
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	state.pc = 2
	state.raiseRuntimeAt(1, errors.New("boom"), recoverDefault, bytecode.NoopOperand, nil, false)

	action := state.resolveFailure()
	if action != errContinue {
		t.Fatalf("expected caught error to continue, got %v", action)
	}

	if got, want := state.pc, 0; got != want {
		t.Fatalf("expected catch jump target %d from failure origin pc, got %d", want, got)
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
	defer state.endRun()

	state.pc = 1
	wantErr := errors.New("boom")

	state.raiseRuntime(wantErr, recoverDefault, bytecode.Operand(1), runtime.True, true)
	action := state.resolveFailure()
	if action != errReturn {
		t.Fatalf("expected original error to be returned, got %v", action)
	}

	if got := state.failureError(); !errors.Is(got, wantErr) {
		t.Fatalf("expected pending failure error to be preserved, got %v", got)
	}

	val := state.registers[1]
	if val == runtime.True {
		t.Fatalf("expected fallback value to be ignored, got %v", val)
	}

	if got, want := state.pc, 1; got != want {
		t.Fatalf("expected pc to stay unchanged at %d, got %d", want, got)
	}
}

func TestSetOrOptional_GenericErrorUsesDefaultResolverInOptionalMode(t *testing.T) {
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
	defer state.endRun()

	state.pc = 1
	state.registers[1] = runtime.True

	state.setOrOptional(state.pc, bytecode.NewRegister(1), runtime.True, errors.New("boom"), true)
	if !state.hasFailure() {
		t.Fatal("expected member error to raise pending failure")
	}

	if got := state.registers[1]; got != runtime.True {
		t.Fatalf("expected destination to remain unchanged before resolution, got %v", got)
	}

	action := state.resolveFailure()
	if action != errContinue {
		t.Fatalf("expected catch path to continue, got %v", action)
	}

	if got := state.registers[1]; got != runtime.None {
		t.Fatalf("expected destination to be reset to none, got %v", got)
	}

	if got, want := state.pc, 0; got != want {
		t.Fatalf("expected generic error to use default resolver and jump to catch %d, got %d", want, got)
	}
}

func TestSetOrOptional_NotFoundContinuesWithNone(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	state.registers[1] = runtime.True

	state.setOrOptional(0, bytecode.NewRegister(1), runtime.True, runtime.ErrNotFound, false)
	if !state.hasFailure() {
		t.Fatal("expected not found path to raise pending failure")
	}

	action := state.resolveFailure()
	if action != errContinue {
		t.Fatalf("expected not found to continue, got %v", action)
	}

	if got := state.registers[1]; got != runtime.None {
		t.Fatalf("expected destination to be reset to none, got %v", got)
	}
}

func TestReturnToCaller_ClosesDiscardedOwnedValuesButPreservesBorrowedArgs(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	borrowed := newTrackingCloser("borrowed")
	produced := newTrackingCloser("produced")

	callerRegs := mem.NewRegisterFile(2)
	callerRegs[0] = borrowed
	callerOwned := mem.OwnedResources{}
	callerOwned.Track(borrowed)

	activeRegs := mem.NewRegisterFile(3)
	activeRegs[1] = borrowed
	activeRegs[2] = produced
	state.registers = activeRegs
	state.owned.Track(produced)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        7,
		ReturnDest:      bytecode.NewRegister(1),
		CallerRegisters: callerRegs,
		OwnedResources:  callerOwned,
	})

	if ok := state.returnToCaller(runtime.True); !ok {
		t.Fatal("expected return to caller to succeed")
	}

	if got := produced.closed; got != 1 {
		t.Fatalf("expected produced callee value to close once, got %d", got)
	}

	if got := borrowed.closed; got != 0 {
		t.Fatalf("expected borrowed argument to remain open, got %d closes", got)
	}

	if !state.owned.Owns(borrowed) {
		t.Fatal("expected caller to keep ownership of borrowed argument")
	}

	if got := state.registers[1]; got != runtime.True {
		t.Fatalf("expected return destination to receive caller result, got %v", got)
	}
}

func TestReturnToCaller_TransfersReturnedOwnedCloserToCaller(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
	})

	state := mustAcquireRunState(t, instance)

	returned := newTrackingCloser("returned")
	activeRegs := mem.NewRegisterFile(1)
	activeRegs[0] = returned
	state.registers = activeRegs
	state.owned.Track(returned)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        9,
		ReturnDest:      bytecode.NewRegister(1),
		CallerRegisters: mem.NewRegisterFile(2),
	})

	if ok := state.returnToCaller(returned); !ok {
		t.Fatal("expected return to caller to succeed")
	}

	if got := returned.closed; got != 0 {
		t.Fatalf("expected returned closer to remain open after transfer, got %d closes", got)
	}

	if !state.owned.Owns(returned) {
		t.Fatal("expected caller to own transferred return value")
	}

	if got := state.registers[1]; got != returned {
		t.Fatalf("expected caller return destination to receive returned value, got %v", got)
	}

	state.endRun()

	if got := returned.closed; got != 1 {
		t.Fatalf("expected root cleanup to close transferred return value once, got %d", got)
	}
}

func TestSetOrOptional_NullDereferenceContinuesWithNone(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	state.registers[1] = runtime.True

	err := rtdiagnostics.MemberAccessErrorOf(
		runtime.None,
		rtdiagnostics.MemberAccessProperty,
		runtime.NewString("foo"),
	)

	state.setOrOptional(0, bytecode.NewRegister(1), runtime.True, err, true)
	if state.hasFailure() {
		t.Fatal("expected null-dereference miss to short-circuit without pending failure")
	}

	if got := state.registers[1]; got != runtime.None {
		t.Fatalf("expected destination to be set to runtime.None for null-dereference miss, got %v", got)
	}
}

func TestLoadKeyCached_FastObjectMissingReturnsNoneWithoutError(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
	})

	ctx := context.Background()
	obj := data.NewFastObject(nil, 0)
	if err := obj.Set(ctx, runtime.NewString("present"), runtime.NewInt(1)); err != nil {
		t.Fatalf("setup fast object failed: %v", err)
	}

	out, err := instance.loadKeyCached(ctx, 0, obj, runtime.NewString("missing"))
	if err != nil {
		t.Fatalf("expected missing fast-object key to return None without error, got %v", err)
	}

	if out != runtime.None {
		t.Fatalf("expected runtime.None on missing fast-object key, got %v", out)
	}

	if instance.cache.LoadKeyICs[0] == nil {
		t.Fatal("expected load-key cache to record fast-object miss")
	}

	out, err = instance.loadKeyCached(ctx, 0, obj, runtime.NewString("missing"))
	if err != nil {
		t.Fatalf("expected cached missing fast-object key to return None without error, got %v", err)
	}

	if out != runtime.None {
		t.Fatalf("expected runtime.None from cached fast-object miss, got %v", out)
	}
}

func TestLoadKeyConstCached_FastObjectMissingReturnsNoneWithoutError(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
	})

	ctx := context.Background()
	obj := data.NewFastObject(nil, 0)
	if err := obj.Set(ctx, runtime.NewString("present"), runtime.NewInt(1)); err != nil {
		t.Fatalf("setup fast object failed: %v", err)
	}

	inst := &instance.plan.instructions[0]
	out, err := instance.loadKeyConstCached(ctx, 0, inst, obj, runtime.NewString("missing"))
	if err != nil {
		t.Fatalf("expected missing const fast-object key to return None without error, got %v", err)
	}

	if out != runtime.None {
		t.Fatalf("expected runtime.None on missing const fast-object key, got %v", out)
	}

	if got, want := inst.InlineSlot, -1; got != want {
		t.Fatalf("expected inline slot %d for cached miss, got %d", want, got)
	}

	out, err = instance.loadKeyConstCached(ctx, 0, inst, obj, runtime.NewString("missing"))
	if err != nil {
		t.Fatalf("expected inline cached missing const key to return None without error, got %v", err)
	}

	if out != runtime.None {
		t.Fatalf("expected runtime.None on inline cached const miss, got %v", out)
	}
}

func TestSetOrOptional_GenericErrorReturnsWithoutCatch(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	wantErr := errors.New("boom")
	initial := runtime.NewInt(42)
	state.registers[1] = initial

	state.setOrOptional(0, bytecode.NewRegister(1), runtime.True, wantErr, true)
	if !state.hasFailure() {
		t.Fatal("expected generic member error to raise pending failure")
	}

	action := state.resolveFailure()
	if action != errReturn {
		t.Fatalf("expected unresolved generic error to return, got %v", action)
	}

	if got := state.failureError(); !errors.Is(got, wantErr) {
		t.Fatalf("expected failure error to be preserved, got %v", got)
	}

	if got := state.registers[1]; got != initial {
		t.Fatalf("expected destination to remain unchanged when error returns, got %v", got)
	}
}

func TestRaiseRuntime_FirstFailureWinsUntilResolved(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	errA := errors.New("A")
	errB := errors.New("B")

	state.raiseRuntime(errA, recoverDefault, bytecode.NoopOperand, nil, false)
	state.raiseRuntime(errB, recoverDefault, bytecode.NoopOperand, nil, false)

	if got := state.failureError(); !errors.Is(got, errA) {
		t.Fatalf("expected first raised error to remain pending, got %v", got)
	}

	if action := state.resolveFailure(); action != errReturn {
		t.Fatalf("expected unresolved failure to return, got %v", action)
	}

	if got := state.failureError(); !errors.Is(got, errA) {
		t.Fatalf("expected returned failure to remain original error, got %v", got)
	}
}

func TestRunStateLifecycleClearsPendingFailure(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		},
	})

	state := mustAcquireRunState(t, instance)
	state.raiseRuntime(errors.New("boom"), recoverDefault, bytecode.NoopOperand, nil, false)
	if !state.hasFailure() {
		t.Fatal("expected pending failure before lifecycle reset")
	}

	state.startRun(NewDefaultEnvironment())
	if state.hasFailure() {
		t.Fatal("expected prepareRun to clear pending failure")
	}

	state.raiseRuntime(errors.New("boom"), recoverDefault, bytecode.NoopOperand, nil, false)
	state.endRun()

	reused := mustAcquireRunState(t, instance)
	defer state.endRun()
	if reused != state {
		t.Fatal("expected state reuse from pool")
	}

	if reused.hasFailure() {
		t.Fatal("expected reset to clear pending failure")
	}
}

func TestTailCallUdf_ReusedWindowResetsNonArgSlotsToNone(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  8,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{
					Name:        "F",
					DisplayName: "f",
					Entry:       0,
					Registers:   6,
					Params:      2,
				},
			},
		},
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	reg := state.registers
	for i := range reg {
		reg[i] = runtime.NewInt(100 + i)
	}

	reg[1] = runtime.NewInt(0)  // UDF id
	reg[3] = runtime.NewInt(10) // arg1
	reg[4] = runtime.NewInt(20) // arg2

	state.frames.Push(frame.CallFrame{
		ReturnPC:        99,
		ReturnDest:      bytecode.NewRegister(0),
		CallerRegisters: make([]runtime.Value, 2),
	})
	state.pc = 5

	oldPtr := &reg[0]
	desc := &callDescriptor{
		DisplayName: "f",
		Dst:         bytecode.NewRegister(1),
		ID:          0,
		ArgStart:    3,
		ArgCount:    2,
		CallSitePC:  4,
	}
	if err := tailCallUdf(state, desc, &instance.program.Functions.UserDefined[0]); err != nil {
		t.Fatalf("unexpected tail call error: %v", err)
	}

	newRegs := state.registers
	if got, want := len(newRegs), 6; got != want {
		t.Fatalf("unexpected tail-call register window size: got %d, want %d", got, want)
	}

	if &newRegs[0] != oldPtr {
		t.Fatal("expected tail call to reuse active register window")
	}

	if got, want := newRegs[1], runtime.NewInt(10); got != want {
		t.Fatalf("unexpected first tail-call argument: got %v, want %v", got, want)
	}

	if got, want := newRegs[2], runtime.NewInt(20); got != want {
		t.Fatalf("unexpected second tail-call argument: got %v, want %v", got, want)
	}

	for _, idx := range []int{0, 3, 4, 5} {
		if got := newRegs[idx]; got != runtime.None {
			t.Fatalf("expected non-argument slot %d to be runtime.None, got %v", idx, got)
		}
	}
}

func TestTailCallUdf_ReusedWindowTransfersOwnedArgsAndClosesDiscardedValues(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  8,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{
					Name:        "F",
					DisplayName: "f",
					Entry:       0,
					Registers:   6,
					Params:      2,
				},
			},
		},
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	discarded := newTrackingCloser("discarded")
	ownedArg := newTrackingCloser("owned-arg")
	borrowedArg := newTrackingCloser("borrowed-arg")

	reg := state.registers
	reg[0] = discarded
	reg[3] = ownedArg
	reg[4] = borrowedArg
	state.owned.Track(discarded)
	state.owned.Track(ownedArg)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        99,
		ReturnDest:      bytecode.NewRegister(0),
		CallerRegisters: mem.NewRegisterFile(2),
	})

	desc := &callDescriptor{
		DisplayName: "f",
		Dst:         bytecode.NewRegister(1),
		ID:          0,
		ArgStart:    3,
		ArgCount:    2,
		CallSitePC:  4,
	}
	if err := tailCallUdf(state, desc, &instance.program.Functions.UserDefined[0]); err != nil {
		t.Fatalf("unexpected tail call error: %v", err)
	}

	if got := discarded.closed; got != 1 {
		t.Fatalf("expected discarded value to close once, got %d", got)
	}

	if got := ownedArg.closed; got != 0 {
		t.Fatalf("expected owned argument to remain open after transfer, got %d closes", got)
	}

	if got := borrowedArg.closed; got != 0 {
		t.Fatalf("expected borrowed argument to remain open, got %d closes", got)
	}

	if !state.owned.Owns(ownedArg) {
		t.Fatal("expected tail call to transfer ownership of direct argument")
	}

	if state.owned.Owns(borrowedArg) {
		t.Fatal("did not expect borrowed argument to become owned")
	}

	if got := state.registers[1]; got != ownedArg {
		t.Fatalf("unexpected transferred argument in first slot: got %v", got)
	}

	if got := state.registers[2]; got != borrowedArg {
		t.Fatalf("unexpected borrowed argument in second slot: got %v", got)
	}
}

func TestTailCallUdf_FreshWindowTransfersOwnedArgsAndClosesDiscardedValues(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  4,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{
					Name:        "F",
					DisplayName: "f",
					Entry:       0,
					Registers:   6,
					Params:      2,
				},
			},
		},
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	discarded := newTrackingCloser("discarded")
	ownedArg := newTrackingCloser("owned-arg")
	borrowedArg := newTrackingCloser("borrowed-arg")

	reg := mem.NewRegisterFile(4)
	reg[0] = discarded
	reg[2] = ownedArg
	reg[3] = borrowedArg
	state.registers = reg
	state.owned.Track(discarded)
	state.owned.Track(ownedArg)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        99,
		ReturnDest:      bytecode.NewRegister(0),
		CallerRegisters: mem.NewRegisterFile(2),
	})

	oldPtr := &reg[0]
	desc := &callDescriptor{
		DisplayName: "f",
		Dst:         bytecode.NewRegister(1),
		ID:          0,
		ArgStart:    2,
		ArgCount:    2,
		CallSitePC:  4,
	}
	if err := tailCallUdf(state, desc, &instance.program.Functions.UserDefined[0]); err != nil {
		t.Fatalf("unexpected tail call error: %v", err)
	}

	if got := discarded.closed; got != 1 {
		t.Fatalf("expected discarded value to close once, got %d", got)
	}

	if got := ownedArg.closed; got != 0 {
		t.Fatalf("expected owned argument to remain open after transfer, got %d closes", got)
	}

	if got := borrowedArg.closed; got != 0 {
		t.Fatalf("expected borrowed argument to remain open, got %d closes", got)
	}

	if !state.owned.Owns(ownedArg) {
		t.Fatal("expected tail call to transfer ownership of direct argument")
	}

	if state.owned.Owns(borrowedArg) {
		t.Fatal("did not expect borrowed argument to become owned")
	}

	if &state.registers[0] == oldPtr {
		t.Fatal("expected fresh-window tail call to install a new register window")
	}

	if got := state.registers[1]; got != ownedArg {
		t.Fatalf("unexpected transferred argument in first slot: got %v", got)
	}

	if got := state.registers[2]; got != borrowedArg {
		t.Fatalf("unexpected borrowed argument in second slot: got %v", got)
	}
}

func TestOpClose_DoesNotDoubleCloseTrackedValueAtRunEnd(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpClose, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{
			runtime.NewString("MAKE"),
		},
	}

	instance := mustNewVM(t, program)
	value := newTrackingCloser("close-me")
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("MAKE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return value, nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	result, err := instance.Run(context.Background(), env)
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}

	if got := value.closed; got != 1 {
		t.Fatalf("expected explicit close to run exactly once, got %d closes", got)
	}

	if got := result; got != runtime.ZeroInt {
		t.Fatalf("unexpected return value: got %v, want %v", got, runtime.ZeroInt)
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
			{1, 1, 3},
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

	_, err := NewWith(program)
	if err == nil {
		t.Fatal("expected initialization error")
	}

	var initErrs *rtdiagnostics.InitializationErrorSet
	if !errors.As(err, &initErrs) {
		t.Fatalf("expected initialization error set, got %T", err)
	}

	if got, want := initErrs.Size(), 1; got != want {
		t.Fatalf("unexpected initialization error count: got %d, want %d", got, want)
	}

	if got := initErrs.First().Cause; !errors.Is(got, ErrInvalidFunctionName) {
		t.Fatalf("expected invalid function name cause, got %v", got)
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

	_, err := NewWith(program)
	if err == nil {
		t.Fatal("expected initialization error")
	}

	var initErrs *rtdiagnostics.InitializationErrorSet
	if !errors.As(err, &initErrs) {
		t.Fatalf("expected initialization error set, got %T", err)
	}

	if got, want := initErrs.Size(), 1; got != want {
		t.Fatalf("unexpected initialization error count: got %d, want %d", got, want)
	}

	if got := initErrs.First().Cause; !errors.Is(got, ErrInvalidFunctionName) {
		t.Fatalf("expected invalid function name cause, got %v", got)
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

	warmup := &rtdiagnostics.WarmupErrorSet{}
	warmup.Add(ErrInvalidFunctionName, 0, bytecode.NewRegister(0))

	err := rtdiagnostics.WrapRuntimeError(program, 1, nil, warmup)
	if err == nil {
		t.Fatal("expected wrapped error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	var rtErrSet *rtdiagnostics.RuntimeErrorSet
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
	defer state.endRun()

	state.frames.Push(frame.CallFrame{
		ReturnPC:         9,
		ReturnDest:       bytecode.NewRegister(1),
		CallerRegisters:  protectedRegs,
		RecoveryBoundary: true,
	})
	state.pc = 1

	state.raiseRuntime(errors.New("boom"), recoverDefault, bytecode.NoopOperand, nil, false)
	action := state.resolveFailure()
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
	defer state.endRun()

	state.registers = activeRegs
	state.frames.Push(frame.CallFrame{
		ReturnPC:         7,
		ReturnDest:       bytecode.NewRegister(1),
		CallerRegisters:  lowerRegs,
		RecoveryBoundary: true,
	})
	state.pc = 1

	state.raiseRuntime(errors.New("boom"), recoverDefault, bytecode.NoopOperand, nil, false)
	action := state.resolveFailure()
	if action != errContinue {
		t.Fatalf("expected continue, got %v", action)
	}

	if got, want := state.pc, 7; got != want {
		t.Fatalf("unexpected unwind target: got %d, want %d", got, want)
	}

	if &state.registers[0] != &lowerRegs[0] {
		t.Fatal("expected protected unwind to resume in the caller register window")
	}

	if got, want := state.frames.Len(), 0; got != want {
		t.Fatalf("expected stack to be unwound, got len %d", got)
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

	base := rtdiagnostics.NewRuntimeError(
		program,
		1,
		UncaughtError,
		"Runtime error",
		"",
		"",
		"existing note",
	)

	stack := []frame.TraceEntry{
		{CallSitePC: 0, FnID: 2, FnName: "inner"},
		{CallSitePC: 0, FnID: 1, FnName: "outer"},
	}

	wrapped := rtdiagnostics.WrapRuntimeError(program, 1, stack, base)

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

	wrappedTwice := rtdiagnostics.WrapRuntimeError(program, 1, stack, wrapped)
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

	base := &rtdiagnostics.RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    UncaughtError,
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

	wrapped := rtdiagnostics.WrapRuntimeError(program, 1, stack, base)

	var rtErr *RuntimeError
	if !errors.As(wrapped, &rtErr) {
		t.Fatalf("expected runtime error, got %T", wrapped)
	}

	formatted := rtErr.Format()
	if strings.Count(formatted, "called from (#1)") != 1 {
		t.Fatalf("expected legacy stack label to be preserved without duplication, got:\n%s", formatted)
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
