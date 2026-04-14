package vm

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	rtdiagnostics "github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
	"github.com/MontFerret/ferret/v2/pkg/vm/test"
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

func mustRunResult(t *testing.T, instance *VM, env *Environment) *Result {
	t.Helper()

	result, err := instance.Run(context.Background(), env)
	if err != nil {
		t.Fatalf("expected successful run, got %v", err)
	}

	return result
}

func mustResultRootAndClose(t *testing.T, result *Result) runtime.Value {
	t.Helper()

	root := result.Root()
	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	return root
}

func countDeferredClosers(deferred *mem.DeferredClosers) int {
	if deferred == nil {
		return 0
	}

	count := 0
	deferred.ForEach(func(io.Closer) {
		count++
	})

	return count
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

	if got, want := rtErr.Message, "unresolved function"; got != want {
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

func TestResultMaterializeIsOneShot(t *testing.T) {
	result := newResult(runtime.NewInt(7))

	got, err := Materialize[int](result, func(value runtime.Value) (Materialized[int], error) {
		val, ok := value.(runtime.Int)
		if !ok {
			t.Fatalf("expected runtime.Int root, got %T", value)
		}

		return Materialized[int]{Value: int(val)}, nil
	})
	if err != nil {
		t.Fatalf("expected first materialization to succeed, got %v", err)
	}

	if got != 7 {
		t.Fatalf("unexpected materialized value: got %d, want 7", got)
	}

	_, err = Materialize[int](result, func(runtime.Value) (Materialized[int], error) {
		return Materialized[int]{Value: 0}, nil
	})
	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected second materialization to fail with invalid operation, got %v", err)
	}

	if !strings.Contains(err.Error(), "result is already materialized") {
		t.Fatalf("expected already materialized error, got %v", err)
	}
}

func TestResultMaterializeNilResultReturnsNilError(t *testing.T) {
	var result *Result
	called := false

	_, err := Materialize[int](result, func(runtime.Value) (Materialized[int], error) {
		called = true
		return Materialized[int]{Value: 0}, nil
	})
	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected nil result materialization to fail with invalid operation, got %v", err)
	}

	if !strings.Contains(err.Error(), "result is nil") {
		t.Fatalf("expected nil result error, got %v", err)
	}

	if called {
		t.Fatal("expected nil result materialization to fail before calling materializer")
	}
}

func TestResultMaterializeClosedResultReturnsClosedError(t *testing.T) {
	result := newResult(runtime.NewInt(7))
	if err := result.Close(); err != nil {
		t.Fatalf("expected close to succeed, got %v", err)
	}

	called := false

	_, err := Materialize[int](result, func(runtime.Value) (Materialized[int], error) {
		called = true
		return Materialized[int]{Value: 0}, nil
	})
	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected closed result materialization to fail with invalid operation, got %v", err)
	}

	if !strings.Contains(err.Error(), "result is closed") {
		t.Fatalf("expected closed result error, got %v", err)
	}

	if called {
		t.Fatal("expected closed result materialization to fail before calling materializer")
	}
}

func TestResultCloseIsIdempotent(t *testing.T) {
	value := newTrackingCloser("result-close")
	result := newResult(value)
	result.AdoptValue(value)

	if err := result.Close(); err != nil {
		t.Fatalf("expected first close to succeed, got %v", err)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected repeated close to succeed, got %v", err)
	}

	if got := value.closed; got != 1 {
		t.Fatalf("expected tracked closer to close once, got %d closes", got)
	}
}

func TestRunResultKeepsReturnedDirectCloserAliveUntilClose(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewString("MAKE"),
		},
	}

	instance := mustNewVM(t, program)
	value := newTrackingCloser("result-root")
	env := mustNewEnvironment(t, WithFunction("MAKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return value, nil
	}))

	result := mustRunResult(t, instance, env)
	if got := value.closed; got != 0 {
		t.Fatalf("expected returned direct closer to stay live until result close, got %d closes", got)
	}

	if got := result.Root(); got != value {
		t.Fatalf("expected result root to expose the live closer, got %v", got)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	if got := value.closed; got != 1 {
		t.Fatalf("expected result close to release returned direct closer once, got %d closes", got)
	}
}

func TestRun_BenchmarkResultModeReusesHandleAfterClose(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewString("MAKE"),
		},
	}

	instance := mustNewVM(t, program, WithTesting(test.WithBenchmarkMode()))
	closers := []*trackingCloser{
		newTrackingCloser("first"),
		newTrackingCloser("second"),
	}
	runCount := 0
	env := mustNewEnvironment(t, WithFunction("MAKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		value := closers[runCount]
		runCount++
		return value, nil
	}))

	result := mustRunResult(t, instance, env)
	if got := result.Root(); got != closers[0] {
		t.Fatalf("expected first run to expose first closer, got %v", got)
	}

	materialized, err := Materialize[string](result, func(value runtime.Value) (Materialized[string], error) {
		closer, ok := value.(*trackingCloser)
		if !ok {
			t.Fatalf("expected tracking closer root, got %T", value)
		}

		return Materialized[string]{Value: closer.name}, nil
	})
	if err != nil {
		t.Fatalf("expected benchmark-mode materialization to succeed, got %v", err)
	}

	if got, want := materialized, "first"; got != want {
		t.Fatalf("unexpected materialized value: got %q, want %q", got, want)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected first benchmark result close to succeed, got %v", err)
	}

	if got, want := closers[0].closed, 1; got != want {
		t.Fatalf("expected first closer to close once, got %d", got)
	}

	reused := mustRunResult(t, instance, env)
	if reused != result {
		t.Fatal("expected benchmark mode to reuse the same result handle")
	}

	if got := reused.Root(); got != closers[1] {
		t.Fatalf("expected second run to expose second closer, got %v", got)
	}

	if err := reused.Close(); err != nil {
		t.Fatalf("expected second benchmark result close to succeed, got %v", err)
	}

	if got, want := closers[1].closed, 1; got != want {
		t.Fatalf("expected second closer to close once, got %d", got)
	}
}

func TestRun_BenchmarkResultModeRequiresCloseBeforeNextRun(t *testing.T) {
	instance := mustNewVM(t, newTestProgram(
		1,
		nil,
		bytecode.NewInstruction(bytecode.OpLoadBool, bytecode.NewRegister(0), bytecode.Operand(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	), WithTesting(test.WithBenchmarkMode()))

	result := mustRunResult(t, instance, nil)

	_, err := instance.Run(context.Background(), nil)
	if err == nil {
		t.Fatal("expected second run without closing benchmark result to fail")
	}

	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected invalid operation, got %v", err)
	}

	if !strings.Contains(err.Error(), "benchmark result must be closed") {
		t.Fatalf("expected benchmark close requirement in error, got %v", err)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected benchmark result close to succeed, got %v", err)
	}

	reused, err := instance.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected run after closing benchmark result to succeed, got %v", err)
	}

	if err := reused.Close(); err != nil {
		t.Fatalf("expected reused benchmark result close to succeed, got %v", err)
	}
}

func TestRun_DefaultModeAllowsMultipleOpenResults(t *testing.T) {
	instance := mustNewVM(t, newTestProgram(
		1,
		nil,
		bytecode.NewInstruction(bytecode.OpLoadBool, bytecode.NewRegister(0), bytecode.Operand(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	))

	first := mustRunResult(t, instance, nil)
	second := mustRunResult(t, instance, nil)

	if first == second {
		t.Fatal("expected default run mode to allocate independent result handles")
	}

	if got := first.Root(); got != runtime.True {
		t.Fatalf("unexpected first result root: got %v", got)
	}

	if got := second.Root(); got != runtime.True {
		t.Fatalf("unexpected second result root: got %v", got)
	}

	if err := first.Close(); err != nil {
		t.Fatalf("expected first result close to succeed, got %v", err)
	}

	if err := second.Close(); err != nil {
		t.Fatalf("expected second result close to succeed, got %v", err)
	}
}

func TestFailedRunClosesDiscardedDeferredClosers(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpFail, bytecode.NewConstant(1)),
		},
		Constants: []runtime.Value{
			runtime.NewString("MAKE"),
			runtime.NewString("boom"),
		},
	}

	instance := mustNewVM(t, program)
	value := newTrackingCloser("discarded-before-failure")
	env := mustNewEnvironment(t, WithFunction("MAKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return value, nil
	}))

	_, err := instance.Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected run to fail")
	}

	if got := value.closed; got != 1 {
		t.Fatalf("expected failed run teardown to close discarded closer once, got %d closes", got)
	}
}

func TestResultCloseReleasesClosersAdoptedBeforeMaterializationFailure(t *testing.T) {
	value := newTrackingCloser("adopted-on-failure")
	result := newResult(runtime.NewArrayWith(value))

	materializeErr := errors.New("materialize failed")
	_, err := Materialize[string](result, func(root runtime.Value) (Materialized[string], error) {
		arr, ok := root.(*runtime.Array)
		if !ok {
			t.Fatalf("expected runtime.Array root, got %T", root)
		}

		item, itemErr := arr.At(context.Background(), runtime.ZeroInt)
		if itemErr != nil {
			t.Fatalf("failed to read array item: %v", itemErr)
		}

		result.AdoptValue(item)

		return Materialized[string]{}, materializeErr
	})
	if !errors.Is(err, materializeErr) {
		t.Fatalf("expected materialization error to be preserved, got %v", err)
	}

	if got := value.closed; got != 0 {
		t.Fatalf("expected adopted closer to remain open until result close, got %d closes", got)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	if got := value.closed; got != 1 {
		t.Fatalf("expected result close to release adopted closer once, got %d closes", got)
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

	out := mustResultRootAndClose(t, mustRunResult(t, mustNewVM(t, program), env))

	if out != runtime.NewInt(3) {
		t.Fatalf("unexpected result: got %v, want %v", out, runtime.NewInt(3))
	}
}

func TestOpLoadParam_MissingParamsPreserveRuntimeError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
		Params:     []string{"foo", "bar"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	env := NewDefaultEnvironment()
	env.Params["foo"] = runtime.NewInt(1)

	_, err := mustNewVM(t, program).Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected missing parameter error")
	}

	if !strings.Contains(err.Error(), "missing parameter") {
		t.Fatalf("unexpected error: %v", err)
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if !strings.Contains(rtErr.Note, "@bar") {
		t.Fatalf("expected missing parameter note to mention @bar, got %q", rtErr.Note)
	}
}

func TestWarmupMissingParamsSingleSiteReturnsRuntimeError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("single_missing_param.fql", "RETURN @foo"),
		Registers:  2,
		Params:     []string{"foo"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{
				{Start: 7, End: 11},
				{Start: 0, End: 11},
			},
		},
	}

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected missing parameter error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	var rtErrSet *rtdiagnostics.RuntimeErrorSet
	if errors.As(err, &rtErrSet) {
		t.Fatalf("expected single runtime error, got set")
	}

	if got, want := rtErr.Message, "missing parameter"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
	}

	if !strings.Contains(rtErr.Note, "@foo") {
		t.Fatalf("expected missing parameter note to mention @foo, got %q", rtErr.Note)
	}
}

func TestWarmupMissingParamsAggregateDifferentSlotsByCallsite(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  4,
		Params:     []string{"foo", "bar"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(2), bytecode.Operand(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected aggregated missing parameter error")
	}

	var rtErrSet *rtdiagnostics.RuntimeErrorSet
	if !errors.As(err, &rtErrSet) {
		t.Fatalf("expected runtime error set, got %T", err)
	}

	if got, want := rtErrSet.Size(), 2; got != want {
		t.Fatalf("unexpected runtime error set size: got %d, want %d", got, want)
	}

	first := rtErrSet.First()
	last := rtErrSet.Last()
	if first == nil || last == nil {
		t.Fatal("expected aggregated runtime errors")
	}

	if !strings.Contains(first.Note, "@foo") {
		t.Fatalf("expected first missing param note to mention @foo, got %q", first.Note)
	}

	if !strings.Contains(last.Note, "@bar") {
		t.Fatalf("expected second missing param note to mention @bar, got %q", last.Note)
	}
}

func TestWarmupMissingParamsAggregateRepeatedCallsites(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  4,
		Params:     []string{"foo"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(2), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected aggregated missing parameter error")
	}

	var rtErrSet *rtdiagnostics.RuntimeErrorSet
	if !errors.As(err, &rtErrSet) {
		t.Fatalf("expected runtime error set, got %T", err)
	}

	if got, want := rtErrSet.Size(), 2; got != want {
		t.Fatalf("unexpected runtime error set size: got %d, want %d", got, want)
	}

	for i, rtErr := range [](*RuntimeError){rtErrSet.First(), rtErrSet.Last()} {
		if rtErr == nil {
			t.Fatalf("expected runtime error at index %d", i)
		}

		if got, want := rtErr.Message, "missing parameter"; got != want {
			t.Fatalf("unexpected runtime error message at index %d: got %q, want %q", i, got, want)
		}

		if !strings.Contains(rtErr.Note, "@foo") {
			t.Fatalf("expected missing parameter note at index %d to mention @foo, got %q", i, rtErr.Note)
		}
	}
}

func TestWarmupHostResolutionPrecedesMissingParamAggregation(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("missing_host_precedence.fql", "RETURN MISSING_FN(@foo)"),
		Registers:  2,
		Params:     []string{"foo"},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{runtime.NewString("MISSING_FN")},
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{
				{Start: 7, End: 17},
				{Start: 7, End: 17},
				{Start: 18, End: 22},
				{Start: 0, End: 22},
			},
		},
	}

	_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
	if err == nil {
		t.Fatal("expected warmup error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if got, want := rtErr.Message, "unresolved function"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
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

	if got := activeOwned.closed; got != 0 {
		t.Fatalf("expected active frame resource to remain deferred until run end, got %d closes", got)
	}

	if got := aboveOwned2.closed; got != 0 {
		t.Fatalf("expected top unwound frame resource to remain deferred until run end, got %d closes", got)
	}

	if got := aboveOwned1.closed; got != 0 {
		t.Fatalf("expected second unwound frame resource to remain deferred until run end, got %d closes", got)
	}

	if got := protectedOwned.closed; got != 0 {
		t.Fatalf("expected protected caller resource to remain open, got %d closes", got)
	}

	if !state.owned.Owns(protectedOwned) {
		t.Fatal("expected protected caller ownership to remain active after unwind")
	}

	if got, want := countDeferredClosers(&state.deferred), 3; got != want {
		t.Fatalf("expected deferred queue size %d after unwind, got %d", want, got)
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

	if got := produced.closed; got != 0 {
		t.Fatalf("expected produced callee value to remain deferred until run end, got %d closes", got)
	}

	if got := borrowed.closed; got != 0 {
		t.Fatalf("expected borrowed argument to remain open, got %d closes", got)
	}

	if !state.owned.Owns(borrowed) {
		t.Fatal("expected caller to keep ownership of borrowed argument")
	}

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred queue size %d after return, got %d", want, got)
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

func TestReturnToCaller_KeepsOwnershipWhileCallerAliasRemains(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	shared := newTrackingCloser("shared")

	callerRegs := mem.NewRegisterFile(2)
	callerRegs[0] = shared
	callerOwned := mem.OwnedResources{}
	callerOwned.Track(shared)
	callerAliases := mem.AliasTracker{}
	if key, _, ok := mem.ResourceKeyOf(shared); ok {
		callerAliases.Inc(key)
	}

	activeRegs := mem.NewRegisterFile(1)
	activeRegs[0] = shared
	state.registers = activeRegs
	state.owned.Track(shared)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        11,
		ReturnDest:      bytecode.NewRegister(1),
		CallerRegisters: callerRegs,
		OwnedResources:  callerOwned,
		Aliases:         callerAliases,
	})

	if ok := state.returnToCaller(shared); !ok {
		t.Fatal("expected return to caller to succeed")
	}

	state.writeBorrowedRegister(bytecode.NewRegister(1), runtime.None)

	if !state.owned.Owns(shared) {
		t.Fatal("expected first caller alias to remain owned after clearing return slot")
	}

	if got := countDeferredClosers(&state.deferred); got != 0 {
		t.Fatalf("expected no deferred closers after clearing one of two aliases, got %d", got)
	}

	state.writeBorrowedRegister(bytecode.NewRegister(0), runtime.None)

	if state.owned.Owns(shared) {
		t.Fatal("expected ownership to end after clearing the final caller alias")
	}

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred queue size %d after clearing final alias, got %d", want, got)
	}

	if got := shared.closed; got != 0 {
		t.Fatalf("expected closer to remain deferred until cleanup, got %d closes", got)
	}
}

func TestReturnToCaller_MatchingDestinationDoesNotNeedExtraTracking(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	shared := newTrackingCloser("shared")

	callerRegs := mem.NewRegisterFile(2)
	callerRegs[1] = shared
	callerOwned := mem.OwnedResources{}
	callerOwned.Track(shared)

	activeRegs := mem.NewRegisterFile(1)
	activeRegs[0] = shared
	state.registers = activeRegs
	state.owned.Track(shared)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        13,
		ReturnDest:      bytecode.NewRegister(1),
		CallerRegisters: callerRegs,
		OwnedResources:  callerOwned,
	})

	if ok := state.returnToCaller(shared); !ok {
		t.Fatal("expected return to caller to succeed")
	}

	state.writeBorrowedRegister(bytecode.NewRegister(1), runtime.None)

	if state.owned.Owns(shared) {
		t.Fatal("expected matching return destination not to require extra ownership tracking")
	}

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred queue size %d after clearing the only caller alias, got %d", want, got)
	}

	if got := shared.closed; got != 0 {
		t.Fatalf("expected closer to remain deferred until cleanup, got %d closes", got)
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

func TestMatchLoadPropertyConst_FastObject(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
		Params:     []string{"obj"},
		Constants:  []runtime.Value{runtime.NewString("a")},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(1), bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpMatchLoadPropertyConst, bytecode.NewRegister(2), bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Metadata: bytecode.Metadata{
			MatchFailTargets: []int{-1, 3, -1, -1, -1},
		},
	}

	instance := mustNewVM(t, program)
	env := mustNewEnvironment(t)

	obj := data.NewFastObject(nil, 0)
	if err := obj.Set(context.Background(), runtime.NewString("a"), runtime.NewInt(7)); err != nil {
		t.Fatalf("setup fast object failed: %v", err)
	}
	env.Params["obj"] = obj

	root := mustResultRootAndClose(t, mustRunResult(t, instance, env))
	if got, want := root, runtime.NewInt(7); got != want {
		t.Fatalf("unexpected match result: got %v, want %v", got, want)
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

	if got := discarded.closed; got != 0 {
		t.Fatalf("expected discarded value to remain deferred until run end, got %d closes", got)
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

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred queue size %d after reused-window tail call, got %d", want, got)
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

	if got := discarded.closed; got != 0 {
		t.Fatalf("expected discarded value to remain deferred until run end, got %d closes", got)
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

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred queue size %d after fresh-window tail call, got %d", want, got)
	}
}

func TestTailCallUdf_DuplicateOwnedArgsStayLiveUntilLastAliasClears(t *testing.T) {
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

	shared := newTrackingCloser("shared")

	reg := state.registers
	reg[3] = shared
	reg[4] = shared
	state.owned.Track(shared)

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

	state.writeBorrowedRegister(bytecode.NewRegister(1), runtime.None)

	if !state.owned.Owns(shared) {
		t.Fatal("expected closer to remain owned after clearing one transferred alias")
	}

	if got := countDeferredClosers(&state.deferred); got != 0 {
		t.Fatalf("expected no deferred closers after clearing one of two aliases, got %d", got)
	}

	state.writeBorrowedRegister(bytecode.NewRegister(2), runtime.None)

	if state.owned.Owns(shared) {
		t.Fatal("expected ownership to end after clearing the final transferred alias")
	}

	if got, want := countDeferredClosers(&state.deferred), 1; got != want {
		t.Fatalf("expected deferred queue size %d after clearing final alias, got %d", want, got)
	}

	if got := shared.closed; got != 0 {
		t.Fatalf("expected closer to remain deferred until cleanup, got %d closes", got)
	}
}

func TestOpClose_DoesNotDoubleCloseTrackedValueAtRunEnd(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
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
	defer func() {
		_ = result.Close()
	}()

	if got := value.closed; got != 1 {
		t.Fatalf("expected explicit close to run exactly once, got %d closes", got)
	}

	if got := result.Root(); got != runtime.ZeroInt {
		t.Fatalf("unexpected return value: got %v, want %v", got, runtime.ZeroInt)
	}
}

func TestOpClose_DuplicateAliasClosesExactlyOnce(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(1), bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpClose, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
		},
		Constants: []runtime.Value{
			runtime.NewString("MAKE"),
		},
	}

	instance := mustNewVM(t, program)
	value := newTrackingCloser("close-me")
	env := mustNewEnvironment(t, WithFunction("MAKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return value, nil
	}))

	result := mustRunResult(t, instance, env)

	if got, want := value.closed, 1; got != want {
		t.Fatalf("expected explicit close to run exactly once with duplicate aliases, got %d closes", got)
	}

	if got := result.Root(); got != runtime.ZeroInt {
		t.Fatalf("unexpected return value: got %v, want %v", got, runtime.ZeroInt)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	if got, want := value.closed, 1; got != want {
		t.Fatalf("expected duplicate alias not to re-close at result cleanup, got %d closes", got)
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
	defer func() {
		_ = result.Close()
	}()

	if got := result.Root(); got != runtime.NewInt(7) {
		t.Fatalf("expected catch jump target to continue execution, got %v", got)
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

	if got, want := rtErr.Message, "invalid type"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
	}

	if got, want := rtErr.Note, "expected String, but got Int"; got != want {
		t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
	}

	if got, want := rtErr.Hint, "Convert the value to String before this operation"; got != want {
		t.Fatalf("unexpected runtime error hint: got %q, want %q", got, want)
	}

	if rtErr.Cause == nil || rtErr.Cause.Error() != runtime.ErrInvalidType.Error() {
		t.Fatalf("expected invalid type cause, got %v", rtErr.Cause)
	}
}

func TestToRuntimeError_InvalidArgumentTypeUsesArgumentSpanAndSeparatesNoteFromCause(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("test.fql", "F(1,true)"),
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(0)),
		},
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{
				{Start: 0, End: 9},
			},
			CallArgumentSpans: [][]source.Span{
				{
					{Start: 2, End: 3},
					{Start: 4, End: 8},
				},
			},
		},
	}

	rtErr := rtdiagnostics.ToRuntimeError(
		program,
		1,
		nil,
		runtime.ArgError(runtime.TypeErrorOf(runtime.True, runtime.TypeString), 1),
	)

	if rtErr == nil {
		t.Fatal("expected runtime error")
	}

	if got, want := rtErr.Kind, diagnostics.TypeError; got != want {
		t.Fatalf("unexpected runtime error kind: got %s, want %s", got, want)
	}

	if got, want := rtErr.Message, "invalid argument type"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
	}

	if got, want := rtErr.Note, "argument 2 expects String, but got Boolean"; got != want {
		t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
	}

	if got, want := rtErr.Hint, "Convert argument 2 to String before this call"; got != want {
		t.Fatalf("unexpected runtime error hint: got %q, want %q", got, want)
	}

	if rtErr.Cause == nil || rtErr.Cause.Error() != runtime.ErrInvalidType.Error() {
		t.Fatalf("expected invalid type cause, got %v", rtErr.Cause)
	}

	formatted := rtErr.Format()
	if !strings.Contains(formatted, "--> test.fql:1:5") {
		t.Fatalf("expected argument span highlight, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "argument 2 has incompatible type") {
		t.Fatalf("expected concrete argument label, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "Hint: Convert argument 2 to String before this call") {
		t.Fatalf("expected synthesized argument type hint in formatted output, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "Caused by: invalid type") {
		t.Fatalf("expected technical cause in formatted output, got:\n%s", formatted)
	}
}

func TestToRuntimeError_InvalidArgumentUsesDedicatedKind(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("test.fql", "F(-1)"),
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(0)),
		},
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{
				{Start: 0, End: 5},
			},
			CallArgumentSpans: [][]source.Span{
				{
					{Start: 2, End: 4},
				},
			},
		},
	}

	rtErr := rtdiagnostics.ToRuntimeError(
		program,
		1,
		nil,
		runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "must be positive"), 0),
	)

	if rtErr == nil {
		t.Fatal("expected runtime error")
	}

	if got, want := rtErr.Kind, InvalidArgument; got != want {
		t.Fatalf("unexpected runtime error kind: got %s, want %s", got, want)
	}

	if got, want := rtErr.Message, "invalid argument"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
	}

	if got, want := rtErr.Note, "argument 1 must be positive"; got != want {
		t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
	}

	if got, want := rtErr.Hint, "Pass argument 1 with a value that is positive"; got != want {
		t.Fatalf("unexpected runtime error hint: got %q, want %q", got, want)
	}

	if rtErr.Cause == nil || rtErr.Cause.Error() != runtime.ErrInvalidArgument.Error() {
		t.Fatalf("expected invalid argument cause, got %v", rtErr.Cause)
	}

	if formatted := rtErr.Format(); !strings.Contains(formatted, "argument 1 is invalid") {
		t.Fatalf("expected concrete invalid argument label, got:\n%s", formatted)
	}
}

func TestToRuntimeError_GenericInvalidArgumentWithoutStructuredDetailOmitsHint(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("test.fql", "F(1)"),
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(0)),
		},
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{
				{Start: 0, End: 4},
			},
		},
	}

	rtErr := rtdiagnostics.ToRuntimeError(
		program,
		1,
		nil,
		runtime.Error(runtime.ErrInvalidArgument, "constraint violated"),
	)

	if rtErr == nil {
		t.Fatal("expected runtime error")
	}

	if rtErr.Hint != "" {
		t.Fatalf("expected opaque invalid argument hint to be omitted, got %q", rtErr.Hint)
	}
}

func TestToRuntimeError_GenericArityErrorSynthesizesHint(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("test.fql", "F(1)"),
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(0)),
		},
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{
				{Start: 0, End: 4},
			},
		},
	}

	rtErr := rtdiagnostics.ToRuntimeError(program, 1, nil, runtime.ArityError(1, 2, 3))
	if rtErr == nil {
		t.Fatal("expected runtime error")
	}

	mainSpan := rtErr.Spans[0]
	if got, want := mainSpan.Label, "wrong number of arguments"; got != want {
		t.Fatalf("unexpected runtime error label: got %q, want %q", got, want)
	}

	if got, want := rtErr.Note, "expected number of arguments 2-3, but got 1"; got != want {
		t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
	}

	if got, want := rtErr.Hint, "Pass 2-3 arguments to this call"; got != want {
		t.Fatalf("unexpected runtime error hint: got %q, want %q", got, want)
	}
}

func TestToRuntimeError_MissingParameterPromotesDetailIntoNote(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("test.fql", "@foo"),
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadParam, bytecode.NewRegister(0), bytecode.Operand(1)),
		},
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{
				{Start: 0, End: 4},
			},
		},
	}

	rtErr := rtdiagnostics.ToRuntimeError(program, 1, nil, runtime.Error(ErrMissedParam, "@foo"))
	if rtErr == nil {
		t.Fatal("expected runtime error")
	}

	if got, want := rtErr.Kind, UnresolvedSymbol; got != want {
		t.Fatalf("unexpected runtime error kind: got %s, want %s", got, want)
	}

	if got, want := rtErr.Message, "missing parameter"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
	}

	if got, want := rtErr.Note, "this query requires parameter '@foo'"; got != want {
		t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
	}

	if got, want := rtErr.Hint, "Provide a value for @foo before executing this query"; got != want {
		t.Fatalf("unexpected runtime error hint: got %q, want %q", got, want)
	}

	if rtErr.Cause == nil || rtErr.Cause.Error() != ErrMissedParam.Error() {
		t.Fatalf("expected missing parameter cause, got %v", rtErr.Cause)
	}

	formatted := rtErr.Format()
	if !strings.Contains(formatted, "parameter '@foo' was not provided") {
		t.Fatalf("expected parameter-specific label, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, "Note: this query requires parameter '@foo'") {
		t.Fatalf("expected factual note in formatted output, got:\n%s", formatted)
	}
}

func TestRuntimeErrorFromPanicUsesStableMessageAndBugHint(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("test.fql", "RETURN 1"),
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{{Start: 0, End: 6}},
		},
	}

	err := rtdiagnostics.RuntimeErrorFromPanic(program, 1, nil, errors.New("boom"))

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if got, want := rtErr.Kind, diagnostics.UnexpectedError; got != want {
		t.Fatalf("unexpected runtime error kind: got %s, want %s", got, want)
	}

	if got, want := rtErr.Message, "unexpected runtime panic"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
	}

	if got, want := rtErr.Note, "panic value: boom"; got != want {
		t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
	}

	if !strings.Contains(rtErr.Hint, "internal VM bug") {
		t.Fatalf("expected bug-report hint, got %q", rtErr.Hint)
	}

	if rtErr.Cause == nil || rtErr.Cause.Error() != "boom" {
		t.Fatalf("expected panic cause to be preserved, got %v", rtErr.Cause)
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

func TestBuildExecPlanRejectsMissingAggregateSelectorSlot(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpAggregateUpdate, bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
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

	if got := initErrs.First().Cause; got == nil || !strings.Contains(got.Error(), "invalid aggregate selector slot") {
		t.Fatalf("expected invalid aggregate selector slot cause, got %v", got)
	}
}

func TestBuildExecPlanRejectsAggregateSelectorSlotLengthMismatch(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpAggregateUpdate, bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Metadata: bytecode.Metadata{
			AggregateSelectorSlots: []int{0},
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

	if got := initErrs.First().Cause; got == nil || !strings.Contains(got.Error(), "metadata length") {
		t.Fatalf("expected aggregate selector slot metadata length cause, got %v", got)
	}
}

func TestBuildExecPlanRejectsMissingMatchFailTarget(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpMatchLoadPropertyConst, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
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

	if got := initErrs.First().Cause; got == nil || !strings.Contains(got.Error(), "invalid match fail target") {
		t.Fatalf("expected invalid match fail target cause, got %v", got)
	}
}

func TestBuildExecPlanRejectsMatchFailTargetLengthMismatch(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpMatchLoadPropertyConst, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Metadata: bytecode.Metadata{
			MatchFailTargets: []int{1},
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

	if got := initErrs.First().Cause; got == nil || !strings.Contains(got.Error(), "match fail target metadata length") {
		t.Fatalf("expected match fail target metadata length cause, got %v", got)
	}
}

func TestWrapRuntimeErrorSingleWarmupFailureReturnsRuntimeError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Source:     source.New("test", "RETURN 1"),
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{{Start: 0, End: 6}},
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

	if got, want := rtErr.Kind, UnresolvedSymbol; got != want {
		t.Fatalf("unexpected runtime error kind: got %s, want %s", got, want)
	}

	if got, want := rtErr.Message, "invalid function name"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
	}

	if got, want := rtErr.Note, "host call target must resolve to a string function name"; got != want {
		t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
	}

	if rtErr.Cause == nil || rtErr.Cause.Error() != ErrInvalidFunctionName.Error() {
		t.Fatalf("expected invalid function name cause, got %v", rtErr.Cause)
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
		Source:     source.New("test", "RETURN 1"),
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{{Start: 0, End: 6}},
		},
	}

	base := rtdiagnostics.NewRuntimeError(
		program,
		1,
		UncaughtError,
		"runtime error",
		"",
		"",
		"existing note",
		nil,
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
		Source:     source.New("test", "RETURN 1"),
		Metadata: bytecode.Metadata{
			DebugSpans: []source.Span{{Start: 0, End: 6}},
		},
	}

	base := &rtdiagnostics.RuntimeError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    UncaughtError,
			Message: "Runtime error",
			Source:  program.Source,
			Spans: []diagnostics.ErrorSpan{
				diagnostics.NewSecondaryErrorSpan(source.Span{Start: 0, End: 6}, "called from (#1)"),
				diagnostics.NewMainErrorSpan(source.Span{Start: 0, End: 6}, ""),
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

	if got, want := rtErr.Kind, ArityError; got != want {
		t.Fatalf("unexpected runtime error kind: got %s, want %s", got, want)
	}

	if got, want := rtErr.Message, "invalid number of arguments"; got != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
	}

	if got, want := rtErr.Note, "UDF 'boo' expects 1 arguments, but got 0"; got != want {
		t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
	}

	formatted := rtErr.Format()
	if got, want := rtErr.Hint, "Pass 1 argument to boo"; got != want {
		t.Fatalf("unexpected runtime error hint: got %q, want %q", got, want)
	}

	if !strings.Contains(formatted, "Hint: Pass 1 argument to boo") {
		t.Fatalf("expected callable-aware arity hint, got:\n%s", formatted)
	}

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

	if got, want := rtErr.Message, "vm invariant violation"; got != want {
		t.Fatalf("unexpected message: %q", rtErr.Message)
	}

	if got, want := rtErr.Note, "invalid aggregate plan index"; got != want {
		t.Fatalf("unexpected note: %q", rtErr.Note)
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

	if got, want := rtErr.Message, "vm invariant violation"; got != want {
		t.Fatalf("unexpected message: %q", rtErr.Message)
	}

	if got, want := rtErr.Note, "invalid collector configuration"; got != want {
		t.Fatalf("unexpected note: %q", rtErr.Note)
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
