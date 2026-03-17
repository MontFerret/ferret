package vm

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Alias Behavior Tests
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func TestLifecycle_SameCloserInMultipleRegisters(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	shared := newTrackingCloser("shared")

	// Store same closer in multiple registers
	state.writeProducedRegister(bytecode.NewRegister(0), shared)
	state.copyRegister(bytecode.NewRegister(1), bytecode.NewRegister(0))
	state.copyRegister(bytecode.NewRegister(2), bytecode.NewRegister(0))

	if !state.owned.Owns(shared) {
		t.Fatal("expected VM to own shared closer")
	}

	// Clear one alias - closer should remain alive
	state.clearRegister(bytecode.NewRegister(0))

	if !state.owned.Owns(shared) {
		t.Fatal("expected VM to still own closer after clearing one of three aliases")
	}

	if got := countDeferredClosers(&state.deferred); got != 0 {
		t.Fatalf("expected no deferred cleanup with live aliases, got %d", got)
	}

	if got := shared.closed; got != 0 {
		t.Fatalf("expected closer to remain open with live aliases, got %d closes", got)
	}

	// Clear second alias
	state.clearRegister(bytecode.NewRegister(1))

	if !state.owned.Owns(shared) {
		t.Fatal("expected VM to still own closer after clearing two of three aliases")
	}

	// Clear final alias - should trigger deferred cleanup
	state.clearRegister(bytecode.NewRegister(2))

	if state.owned.Owns(shared) {
		t.Fatal("expected ownership to end after clearing final alias")
	}

	if got := countDeferredClosers(&state.deferred); got != 1 {
		t.Fatalf("expected deferred cleanup after clearing final alias, got %d", got)
	}

	if got := shared.closed; got != 0 {
		t.Fatalf("expected closer to remain deferred until explicit cleanup, got %d closes", got)
	}
}

func TestLifecycle_OverwriteOneAliasKeepsOthersAlive(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	shared := newTrackingCloser("shared")
	replacement := newTrackingCloser("replacement")

	// Create multiple aliases
	state.writeProducedRegister(bytecode.NewRegister(0), shared)
	state.copyRegister(bytecode.NewRegister(1), bytecode.NewRegister(0))
	state.copyRegister(bytecode.NewRegister(2), bytecode.NewRegister(0))

	// Overwrite one alias with different value
	state.writeProducedRegister(bytecode.NewRegister(1), replacement)

	if !state.owned.Owns(shared) {
		t.Fatal("expected shared closer to remain owned after overwriting one alias")
	}

	if !state.owned.Owns(replacement) {
		t.Fatal("expected replacement closer to be owned")
	}

	if got := shared.closed; got != 0 {
		t.Fatalf("expected shared closer to remain open, got %d closes", got)
	}

	if got := countDeferredClosers(&state.deferred); got != 0 {
		t.Fatalf("expected no deferred cleanup with live aliases, got %d", got)
	}
}

func TestLifecycle_OverwriteLastAliasTriggersDeferred(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	original := newTrackingCloser("original")
	replacement := newTrackingCloser("replacement")

	// Single alias
	state.writeProducedRegister(bytecode.NewRegister(0), original)

	// Overwrite it
	state.writeProducedRegister(bytecode.NewRegister(1), replacement)
	state.writeBorrowedRegister(bytecode.NewRegister(0), runtime.None)

	if state.owned.Owns(original) {
		t.Fatal("expected original closer ownership to end after overwriting final alias")
	}

	if got := countDeferredClosers(&state.deferred); got != 1 {
		t.Fatalf("expected deferred cleanup for original, got %d", got)
	}

	if got := original.closed; got != 0 {
		t.Fatalf("expected original to be deferred not closed, got %d closes", got)
	}
}

func TestLifecycle_AliasesAcrossFrameBoundaries(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
	})

	state := mustAcquireRunState(t, instance)

	callerResource := newTrackingCloser("caller-owned")
	calleeResource := newTrackingCloser("callee-owned")

	// Caller frame owns one resource
	callerRegs := mem.NewRegisterFile(2)
	callerRegs[0] = callerResource
	callerOwned := mem.OwnedResources{}
	callerOwned.Track(callerResource)

	// Callee frame has passed-in arg (not owned) and its own resource
	state.registers[0] = callerResource // borrowed arg
	state.registers[1] = calleeResource
	state.registers[2] = calleeResource
	state.owned.Track(calleeResource)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        10,
		ReturnDest:      bytecode.NewRegister(1),
		CallerRegisters: callerRegs,
		OwnedResources:  callerOwned,
	})

	// Clear one of callee's resource aliases
	state.clearRegister(bytecode.NewRegister(1))

	// Callee still owns its resource (has another alias)
	if !state.owned.Owns(calleeResource) {
		t.Fatal("expected callee frame to still own its resource (has another alias)")
	}

	// Clear second alias
	state.clearRegister(bytecode.NewRegister(2))

	// Callee loses ownership of its resource
	if state.owned.Owns(calleeResource) {
		t.Fatal("expected callee frame ownership to end after clearing final alias")
	}

	if got := countDeferredClosers(&state.deferred); got != 1 {
		t.Fatalf("expected deferred cleanup for callee resource, got %d", got)
	}

	// Return to caller
	if ok := state.returnToCaller(runtime.True); !ok {
		t.Fatal("expected return to caller to succeed")
	}

	// Caller still owns its resource
	if !state.owned.Owns(callerResource) {
		t.Fatal("expected caller to still own its resource after callee return")
	}

	// Callee's resource should be deferred but not yet closed
	if got := calleeResource.closed; got != 0 {
		t.Fatalf("expected callee resource to be deferred not closed, got %d closes", got)
	}

	if got := callerResource.closed; got != 0 {
		t.Fatalf("expected caller resource to remain open, got %d closes", got)
	}

	// Deferred cleanup happens at endRun
	state.endRun()

	if got := calleeResource.closed; got != 1 {
		t.Fatalf("expected callee resource to close at run end, got %d closes", got)
	}

	if got := callerResource.closed; got != 1 {
		t.Fatalf("expected caller resource to close at run end, got %d closes", got)
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Explicit Close Tests
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func TestLifecycle_ExplicitCloseOneAliasLeavesOtherStale(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(1), bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(2), bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpClose, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{
			runtime.NewString("MAKE"),
		},
	}

	instance := mustNewVM(t, program)
	value := newTrackingCloser("explicit-close")
	env := mustNewEnvironment(t, WithFunction("MAKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return value, nil
	}))

	result := mustRunResult(t, instance, env)

	// Explicit close happened during run
	if got := value.closed; got != 1 {
		t.Fatalf("expected explicit close to happen once, got %d closes", got)
	}

	// Stale aliases should be ignored
	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	// No second close
	if got := value.closed; got != 1 {
		t.Fatalf("expected stale alias not to trigger second close, got %d closes", got)
	}
}

func TestLifecycle_ExplicitCloseRemovesOwnership(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
	})

	state := mustAcquireRunState(t, instance)
	defer state.endRun()

	closer := newTrackingCloser("explicit")

	state.writeProducedRegister(bytecode.NewRegister(0), closer)
	state.copyRegister(bytecode.NewRegister(1), bytecode.NewRegister(0))

	if !state.owned.Owns(closer) {
		t.Fatal("expected VM to own closer before explicit close")
	}

	// Simulate OpClose
	released, ok := state.owned.Release(closer)
	if !ok {
		t.Fatal("expected release to succeed")
	}

	if err := released.Close(); err != nil {
		t.Fatalf("expected explicit close to succeed, got %v", err)
	}

	if got := closer.closed; got != 1 {
		t.Fatalf("expected closer to close once, got %d", got)
	}

	// Ownership is removed
	if state.owned.Owns(closer) {
		t.Fatal("expected ownership to be removed after explicit close")
	}

	// Clearing stale aliases should not schedule deferred cleanup
	state.clearRegister(bytecode.NewRegister(0))
	state.clearRegister(bytecode.NewRegister(1))

	if got := countDeferredClosers(&state.deferred); got != 0 {
		t.Fatalf("expected no deferred cleanup for already-closed resource, got %d", got)
	}

	// No second close
	if got := closer.closed; got != 1 {
		t.Fatalf("expected closer to remain closed exactly once, got %d", got)
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Transfer Behavior Tests
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func TestLifecycle_ReturnedOwnedResourceNotClosedByFrameEnd(t *testing.T) {
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
	resource := newTrackingCloser("returned")
	env := mustNewEnvironment(t, WithFunction("MAKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return resource, nil
	}))

	result := mustRunResult(t, instance, env)

	// Frame ended, but resource is not closed (transferred to Result)
	if got := resource.closed; got != 0 {
		t.Fatalf("expected returned resource to remain open after frame end, got %d closes", got)
	}

	if got := result.Root(); got != resource {
		t.Fatalf("expected result root to be returned resource, got %v", got)
	}
}

func TestLifecycle_ReturnedOwnedResourceClosedByResultClose(t *testing.T) {
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
	resource := newTrackingCloser("returned")
	env := mustNewEnvironment(t, WithFunction("MAKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return resource, nil
	}))

	result := mustRunResult(t, instance, env)

	if got := resource.closed; got != 0 {
		t.Fatalf("expected returned resource to remain open before result close, got %d closes", got)
	}

	// Result.Close() should close the transferred resource
	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	if got := resource.closed; got != 1 {
		t.Fatalf("expected Result.Close() to close returned resource once, got %d closes", got)
	}
}

func TestLifecycle_TransferWithMultipleAliases(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  3,
	})

	state := mustAcquireRunState(t, instance)

	shared := newTrackingCloser("shared")

	// Create multiple aliases in callee
	activeRegs := mem.NewRegisterFile(3)
	activeRegs[0] = shared
	activeRegs[1] = shared
	activeRegs[2] = shared
	state.registers = activeRegs
	state.owned.Track(shared)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        5,
		ReturnDest:      bytecode.NewRegister(0),
		CallerRegisters: mem.NewRegisterFile(1),
	})

	// Return one alias
	if ok := state.returnToCaller(shared); !ok {
		t.Fatal("expected return to caller to succeed")
	}

	// Transferred to caller, should not be closed during cleanup
	if got := shared.closed; got != 0 {
		t.Fatalf("expected transferred resource to remain open, got %d closes", got)
	}

	if !state.owned.Owns(shared) {
		t.Fatal("expected caller to own transferred resource")
	}

	// End run should close it
	state.endRun()

	if got := shared.closed; got != 1 {
		t.Fatalf("expected end run cleanup to close transferred resource once, got %d closes", got)
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Failure Path Tests
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func TestLifecycle_RuntimeErrorAfterOwnedResourceCreation(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewString("MAKE"),
			runtime.NewString("FAIL"),
		},
	}

	instance := mustNewVM(t, program)
	resource := newTrackingCloser("leaked")
	testErr := errors.New("intentional failure")

	env := mustNewEnvironment(t,
		WithFunction("MAKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return resource, nil
		}),
		WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return nil, testErr
		}),
	)

	_, err := instance.Run(context.Background(), env)
	if !errors.Is(err, testErr) {
		t.Fatalf("expected intentional failure, got %v", err)
	}

	// Resource should be cleaned up despite error
	if got := resource.closed; got != 1 {
		t.Fatalf("expected run-end cleanup to close resource once, got %d closes", got)
	}
}

func TestLifecycle_NoLeakedOwnershipAfterError(t *testing.T) {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewString("FAIL"),
		},
	}

	instance := mustNewVM(t, program)
	leaked := newTrackingCloser("leaked")
	testErr := errors.New("fail immediately")

	callCount := 0
	env := mustNewEnvironment(t, WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		callCount++
		if callCount == 1 {
			return nil, testErr
		}
		return leaked, nil
	}))

	// First run fails
	_, err := instance.Run(context.Background(), env)
	if !errors.Is(err, testErr) {
		t.Fatalf("expected test error, got %v", err)
	}

	// Second run succeeds - verify no leaked state
	result, err := instance.Run(context.Background(), env)
	if err != nil {
		t.Fatalf("expected second run to succeed, got %v", err)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	if got := leaked.closed; got != 1 {
		t.Fatalf("expected clean second run to close resource once, got %d closes", got)
	}
}

func TestLifecycle_ErrorDuringFrameReturnCleansUpCallee(t *testing.T) {
	instance := mustNewVM(t, &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  2,
	})

	state := mustAcquireRunState(t, instance)

	calleeOwned := newTrackingCloser("callee-owned")

	// Callee has owned resource
	activeRegs := mem.NewRegisterFile(1)
	activeRegs[0] = calleeOwned
	state.registers = activeRegs
	state.owned.Track(calleeOwned)

	state.frames.Push(frame.CallFrame{
		ReturnPC:        10,
		ReturnDest:      bytecode.NewRegister(0),
		CallerRegisters: mem.NewRegisterFile(2),
	})

	// Simulate error before return
	state.raiseRuntime(errors.New("error before return"), recoverDefault, bytecode.NoopOperand, nil, false)

	// End run should clean up callee resources
	state.endRun()

	if got := calleeOwned.closed; got != 1 {
		t.Fatalf("expected error path cleanup to close callee resource once, got %d closes", got)
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Materialize Interaction Tests
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func TestLifecycle_MaterializerAdoptsCloserViaResult(t *testing.T) {
	container := newTrackingCloser("container")
	nested := newTrackingCloser("nested")
	result := newResult(container)
	result.AdoptValue(container)

	_, err := Materialize[string](result, func(root runtime.Value) (Materialized[string], error) {
		// Materializer discovers nested resource and returns it for adoption
		return Materialized[string]{
			Value:   "materialized",
			Closers: []io.Closer{nested},
		}, nil
	})
	if err != nil {
		t.Fatalf("expected materialization to succeed, got %v", err)
	}

	if got := container.closed; got != 0 {
		t.Fatalf("expected container to remain open before result close, got %d closes", got)
	}

	if got := nested.closed; got != 0 {
		t.Fatalf("expected adopted closer to remain open before result close, got %d closes", got)
	}

	// Result.Close() should close both
	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	if got := container.closed; got != 1 {
		t.Fatalf("expected container to close once, got %d closes", got)
	}

	if got := nested.closed; got != 1 {
		t.Fatalf("expected adopted closer to close once, got %d closes", got)
	}
}

func TestLifecycle_SecondMaterializeFailsPredictably(t *testing.T) {
	result := newResult(runtime.NewInt(42))

	// First materialization succeeds
	_, err := Materialize[int](result, func(val runtime.Value) (Materialized[int], error) {
		return Materialized[int]{Value: 42}, nil
	})
	if err != nil {
		t.Fatalf("expected first materialization to succeed, got %v", err)
	}

	// Second materialization fails
	_, err = Materialize[int](result, func(val runtime.Value) (Materialized[int], error) {
		t.Fatal("materializer should not be called for second materialization")
		return Materialized[int]{}, nil
	})
	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected second materialization to fail with ErrInvalidOperation, got %v", err)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}
}

func TestLifecycle_MaterializeErrorDoesNotLeakClosers(t *testing.T) {
	adopted := newTrackingCloser("adopted-on-error")
	result := newResult(runtime.True)

	materializeErr := errors.New("materialize failed")
	_, err := Materialize[bool](result, func(root runtime.Value) (Materialized[bool], error) {
		// Return error - closers in return value are NOT adopted when error is returned
		return Materialized[bool]{
			Closers: []io.Closer{adopted},
		}, materializeErr
	})
	if !errors.Is(err, materializeErr) {
		t.Fatalf("expected materialize error to propagate, got %v", err)
	}

	// When materializer fails, closers are NOT adopted (error path returns early)
	// So Result.Close() will NOT close them
	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	// Closer was never adopted due to error, so it's not closed by Result
	if got := adopted.closed; got != 0 {
		t.Fatalf("expected non-adopted closer to remain open (materializer failed before adoption), got %d closes", got)
	}
}

func TestLifecycle_MaterializeAdoptsMultipleClosers(t *testing.T) {
	closer1 := newTrackingCloser("closer1")
	closer2 := newTrackingCloser("closer2")
	closer3 := newTrackingCloser("closer3")

	result := newResult(runtime.True)

	_, err := Materialize[string](result, func(root runtime.Value) (Materialized[string], error) {
		return Materialized[string]{
			Value:   "multi",
			Closers: []io.Closer{closer1, closer2, closer3},
		}, nil
	})
	if err != nil {
		t.Fatalf("expected materialization to succeed, got %v", err)
	}

	if err := result.Close(); err != nil {
		t.Fatalf("expected result close to succeed, got %v", err)
	}

	// All should be closed
	if got := closer1.closed; got != 1 {
		t.Fatalf("expected closer1 to close once, got %d", got)
	}

	if got := closer2.closed; got != 1 {
		t.Fatalf("expected closer2 to close once, got %d", got)
	}

	if got := closer3.closed; got != 1 {
		t.Fatalf("expected closer3 to close once, got %d", got)
	}
}
