package vm

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func TestNewWith_InlinesHostCallIDs(t *testing.T) {
	program := newHostCallProgram(hostCallSpec{name: "F", args: []runtime.Value{runtime.NewInt(1)}})
	instance := mustNewVM(t, program)

	if len(instance.plan.hostCallDescriptors) == 0 {
		t.Fatal("expected host call warmup descriptors")
	}

	for i, binding := range instance.plan.hostCallDescriptors {
		if got, want := binding.ID, i; got != want {
			t.Fatalf("unexpected host binding id at index %d: got %d, want %d", i, got, want)
		}

		if binding.ID < 0 || binding.ID >= len(instance.cache.HostFunctions) {
			t.Fatalf("invalid binding id %d for pc %d", binding.ID, binding.PC)
		}

		if binding.PC < 0 || binding.PC >= len(instance.plan.instructions) {
			t.Fatalf("invalid callsite pc %d", binding.PC)
		}

		inst := instance.plan.instructions[binding.PC]
		if inst.Opcode != bytecode.OpHCall && inst.Opcode != bytecode.OpProtectedHCall {
			t.Fatalf("callsite pc %d does not point to host call opcode %d", binding.PC, inst.Opcode)
		}

		if got, want := inst.InlineSlot, binding.ID; got != want {
			t.Fatalf("unexpected inlined host id at pc %d: got %d, want %d", binding.PC, got, want)
		}
	}
}

func TestNewWith_HostCallIDsAreCompactAndOrdered(t *testing.T) {
	program := newHostCallProgram(
		hostCallSpec{name: "F", args: []runtime.Value{runtime.NewInt(1)}},
		hostCallSpec{name: "G", args: []runtime.Value{runtime.NewInt(2)}},
	)
	instance := mustNewVM(t, program)

	if got, want := len(instance.plan.hostCallDescriptors), 2; got != want {
		t.Fatalf("unexpected host binding count: got %d, want %d", got, want)
	}

	prevPC := -1
	for i, binding := range instance.plan.hostCallDescriptors {
		if got, want := binding.ID, i; got != want {
			t.Fatalf("unexpected host binding id at index %d: got %d, want %d", i, got, want)
		}
	}

	used := make([]bool, len(instance.plan.hostCallDescriptors))
	siteByPC := make(map[int]callDescriptor, len(instance.plan.hostCallDescriptors))
	for _, binding := range instance.plan.hostCallDescriptors {
		if binding.PC <= prevPC {
			t.Fatalf("host warmup pcs are not increasing: prev=%d, curr=%d", prevPC, binding.PC)
		}
		prevPC = binding.PC

		if binding.ID < 0 || binding.ID >= len(instance.plan.hostCallDescriptors) {
			t.Fatalf("invalid inlined host id at pc %d: %d", binding.PC, binding.ID)
		}

		used[binding.ID] = true
		siteByPC[binding.PC] = binding

		inst := instance.plan.instructions[binding.PC]
		if got, want := inst.InlineSlot, binding.ID; got != want {
			t.Fatalf("unexpected inlined host id at pc %d: got %d, want %d", binding.PC, got, want)
		}
	}

	for i, ok := range used {
		if !ok {
			t.Fatalf("host binding %d is never referenced by a callsite", i)
		}
	}

	hostCallsites := 0
	for pc, inst := range instance.plan.instructions {
		if inst.Opcode != bytecode.OpHCall && inst.Opcode != bytecode.OpProtectedHCall {
			continue
		}

		hostCallsites++

		if inst.InlineSlot < 0 || inst.InlineSlot >= len(instance.plan.hostCallDescriptors) {
			t.Fatalf("invalid inlined host id at pc %d: %d", pc, inst.InlineSlot)
		}

		site, ok := siteByPC[pc]
		if !ok {
			t.Fatalf("host callsite pc %d missing from warmup site metadata", pc)
		}

		if got, want := site.ID, inst.InlineSlot; got != want {
			t.Fatalf("host id mismatch at pc %d: got %d, want %d", pc, got, want)
		}
	}

	if got, want := hostCallsites, len(instance.plan.hostCallDescriptors); got != want {
		t.Fatalf("unexpected host callsite count: got %d, want %d", got, want)
	}
}

func TestNewWith_AssignsOneHostBindingPerCallsite(t *testing.T) {
	program := newHostCallProgram(
		hostCallSpec{name: "F", args: []runtime.Value{runtime.NewInt(1)}},
		hostCallSpec{name: "F", args: []runtime.Value{runtime.NewInt(2)}},
	)
	instance := mustNewVM(t, program)

	if got, want := len(instance.plan.hostCallDescriptors), 2; got != want {
		t.Fatalf("unexpected host binding count: got %d, want %d", got, want)
	}

	firstID := instance.plan.hostCallDescriptors[0].ID
	secondID := instance.plan.hostCallDescriptors[1].ID
	if firstID == secondID {
		t.Fatalf("expected distinct binding ids per callsite, got %d", firstID)
	}

	for i, binding := range instance.plan.hostCallDescriptors {
		if got, want := binding.ID, i; got != want {
			t.Fatalf("unexpected binding id for site %d: got %d, want %d", i, got, want)
		}

		inst := instance.plan.instructions[binding.PC]
		if got, want := inst.InlineSlot, binding.ID; got != want {
			t.Fatalf("unexpected inlined host id at pc %d: got %d, want %d", binding.PC, got, want)
		}
	}
}

func TestNewWith_AssignsDistinctBindingsAcrossBindClasses(t *testing.T) {
	program := newHostCallProgram(
		hostCallSpec{name: "F"},
		hostCallSpec{name: "F", args: []runtime.Value{runtime.NewInt(1)}},
	)
	instance := mustNewVM(t, program)

	if got, want := len(instance.plan.hostCallDescriptors), 2; got != want {
		t.Fatalf("unexpected host binding count: got %d, want %d", got, want)
	}

	firstID := instance.plan.hostCallDescriptors[0].ID
	secondID := instance.plan.hostCallDescriptors[1].ID
	if firstID == secondID {
		t.Fatalf("expected distinct binding ids for different bind classes, got %d", firstID)
	}
}

func TestRun_IgnoresInlineSlotForNonHostOpcodes(t *testing.T) {
	program := newHostCallProgram(hostCallSpec{name: "F", args: []runtime.Value{runtime.NewInt(1)}})
	instance := mustNewVM(t, program)

	nonHostPC := -1
	for pc, inst := range instance.plan.instructions {
		if inst.Opcode != bytecode.OpHCall && inst.Opcode != bytecode.OpProtectedHCall {
			nonHostPC = pc
			break
		}
	}

	if nonHostPC < 0 {
		t.Fatal("expected at least one non-host opcode")
	}

	instance.plan.instructions[nonHostPC].InlineSlot = 1 << 20

	env := mustNewEnvironment(t, WithFunction("F", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
		return args[0], nil
	}))

	out, err := instance.Run(context.Background(), env)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	assertRuntimeValueEquals(t, out, runtime.NewInt(1))
}

func TestRunReturnsUnresolvedFunctionWhenHostCacheEntryIsMissing(t *testing.T) {
	program := newHostCallProgram(hostCallSpec{name: "TEST", args: []runtime.Value{runtime.NewInt(1)}})
	instance := mustNewVM(t, program)

	env := mustNewEnvironment(t, WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.True, nil
	}))

	if _, err := instance.Run(context.Background(), env); err != nil {
		t.Fatalf("first run failed: %v", err)
	}

	hostPC := instance.plan.hostCallDescriptors[0].PC
	hostID := instance.plan.instructions[hostPC].InlineSlot
	if hostID < 0 || hostID >= len(instance.cache.HostFunctions) {
		t.Fatalf("invalid host id at pc %d: %d", hostPC, hostID)
	}

	instance.cache.HostFunctions[hostID] = mem.CachedHostFunction{}
	instance.cache.HostFunctions[hostID].Bound = false

	_, err := instance.Run(context.Background(), env)
	rtErr := assertUnresolvedFunctionError(t, err)
	if rtErr.Message != "Unresolved function" {
		t.Fatalf("expected unresolved function message, got %q", rtErr.Message)
	}
}

func TestWarmupRebindTouchesOnlyHostCallSlots(t *testing.T) {
	program := newHostCallProgram(hostCallSpec{name: "F"})
	instance := mustNewVM(t, program)

	if got, want := len(instance.plan.hostCallDescriptors), 1; got != want {
		t.Fatalf("unexpected host binding count: got %d, want %d", got, want)
	}

	hostID := instance.plan.hostCallDescriptors[0].ID
	if hostID < 0 || hostID >= len(instance.cache.HostFunctions) {
		t.Fatalf("invalid host id %d", hostID)
	}

	sentinel := mem.CachedHostFunction{
		FnV: func(_ context.Context, _ ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(77), nil
		},
	}
	instance.cache.HostFunctions[hostID] = sentinel
	instance.cache.HostFunctions[hostID].Bound = true

	envA := mustNewEnvironment(t, WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.NewInt(1), nil
	}))
	envB := mustNewEnvironment(t,
		WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(2), nil
		}),
		WithFunction("G", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(3), nil
		}),
	)

	outA, err := instance.Run(context.Background(), envA)
	if err != nil {
		t.Fatalf("first run failed: %v", err)
	}
	assertRuntimeValueEquals(t, outA, runtime.NewInt(1))

	outB, err := instance.Run(context.Background(), envB)
	if err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	assertRuntimeValueEquals(t, outB, runtime.NewInt(2))

	if !instance.cache.HostFunctions[hostID].Bound {
		t.Fatal("expected host call slot to be rebound")
	}
}

func TestWarmupPreservesFunctionsRefAcrossFailedAndRecoveredRuns(t *testing.T) {
	program := newHostCallProgram(hostCallSpec{name: "F"})
	instance := mustNewVM(t, program)

	hostID := instance.plan.hostCallDescriptors[0].ID
	validEnv := mustNewEnvironment(t, WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
			return runtime.NewInt(7), nil
		})
	}))
	missingEnv := NewDefaultEnvironment()
	recoveredEnv := mustNewEnvironment(t, WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
			return runtime.NewInt(9), nil
		})
	}))

	out, err := instance.Run(context.Background(), validEnv)
	if err != nil {
		t.Fatalf("initial run failed: %v", err)
	}
	assertRuntimeValueEquals(t, out, runtime.NewInt(7))

	if instance.cache.FunctionsRef != validEnv.Functions {
		t.Fatal("expected FunctionsRef to track the successful environment")
	}

	assertUnresolvedFunctionError(t, func() error {
		_, err := instance.Run(context.Background(), missingEnv)
		return err
	}())
	if instance.cache.FunctionsRef != validEnv.Functions {
		t.Fatal("expected failed warmup to preserve the last successful FunctionsRef")
	}
	if instance.cache.HostFunctions[hostID].Bound {
		t.Fatal("expected failed binding to clear the unresolved host slot")
	}

	out, err = instance.Run(context.Background(), recoveredEnv)
	if err != nil {
		t.Fatalf("recovery run failed: %v", err)
	}
	assertRuntimeValueEquals(t, out, runtime.NewInt(9))

	if instance.cache.FunctionsRef != recoveredEnv.Functions {
		t.Fatal("expected recovery run to install the recovered FunctionsRef")
	}
	if !instance.cache.HostFunctions[hostID].Bound {
		t.Fatal("expected host slot to be rebound after recovery")
	}
}

func TestWarmupRebindReplacesCachedFunctionShapeForCurrentEnvironment(t *testing.T) {
	program := newHostCallProgram(hostCallSpec{name: "F", args: []runtime.Value{runtime.NewInt(1), runtime.NewInt(2)}})
	instance := mustNewVM(t, program)

	hostID := instance.plan.hostCallDescriptors[0].ID
	envA2 := mustNewEnvironment(t, WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.A2().Add("F", func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(12), nil
		})
	}))
	envVar := mustNewEnvironment(t, WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.NewInt(102), nil
	}))
	envA2Again := mustNewEnvironment(t, WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.A2().Add("F", func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(22), nil
		})
	}))

	out, err := instance.Run(context.Background(), envA2)
	if err != nil {
		t.Fatalf("fixed-arity run failed: %v", err)
	}
	assertRuntimeValueEquals(t, out, runtime.NewInt(12))

	slot := instance.cache.HostFunctions[hostID]
	if slot.Fn2 == nil || slot.FnV != nil {
		t.Fatalf("expected host slot to be specialized as Fn2 after fixed-arity bind: %+v", slot)
	}

	out, err = instance.Run(context.Background(), envVar)
	if err != nil {
		t.Fatalf("vararg run failed: %v", err)
	}
	assertRuntimeValueEquals(t, out, runtime.NewInt(102))

	slot = instance.cache.HostFunctions[hostID]
	if slot.Fn2 != nil || slot.FnV == nil {
		t.Fatalf("expected host slot to be rebound as FnV after vararg bind: %+v", slot)
	}

	out, err = instance.Run(context.Background(), envA2Again)
	if err != nil {
		t.Fatalf("second fixed-arity run failed: %v", err)
	}
	assertRuntimeValueEquals(t, out, runtime.NewInt(22))

	slot = instance.cache.HostFunctions[hostID]
	if slot.Fn2 == nil || slot.FnV != nil {
		t.Fatalf("expected host slot to switch back to Fn2 after rebind: %+v", slot)
	}
}

func TestWarmupPartialFailureLeavesMissingSlotsUnboundUntilRecovery(t *testing.T) {
	program := newHostCallProgram(
		hostCallSpec{name: "F"},
		hostCallSpec{name: "G"},
		hostCallSpec{name: "F"},
	)
	instance := mustNewVM(t, program)

	envAll := mustNewEnvironment(t, WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
			return runtime.NewInt(1), nil
		})
		fns.A0().Add("G", func(context.Context) (runtime.Value, error) {
			return runtime.NewInt(2), nil
		})
	}))
	missingEnv := mustNewEnvironment(t, WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
			return runtime.NewInt(20), nil
		})
	}))
	recoveredEnv := mustNewEnvironment(t, WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
			return runtime.NewInt(30), nil
		})
		fns.A0().Add("G", func(context.Context) (runtime.Value, error) {
			return runtime.NewInt(40), nil
		})
	}))

	if _, err := instance.Run(context.Background(), envAll); err != nil {
		t.Fatalf("initial multi-callsite run failed: %v", err)
	}
	if instance.cache.FunctionsRef != envAll.Functions {
		t.Fatal("expected FunctionsRef to track the successful multi-callsite environment")
	}

	assertUnresolvedFunctionError(t, func() error {
		_, err := instance.Run(context.Background(), missingEnv)
		return err
	}())
	if instance.cache.FunctionsRef != envAll.Functions {
		t.Fatal("expected failed multi-callsite warmup to preserve the last successful FunctionsRef")
	}

	for _, descriptor := range instance.plan.hostCallDescriptors {
		cached := instance.cache.HostFunctions[descriptor.ID]
		switch descriptor.DisplayName {
		case "F":
			if !cached.Bound {
				t.Fatalf("expected F slot %d to remain bound during partial warmup failure", descriptor.ID)
			}
		case "G":
			if cached.Bound {
				t.Fatalf("expected missing G slot %d to remain unbound after failed warmup", descriptor.ID)
			}
		}
	}

	if _, err := instance.Run(context.Background(), recoveredEnv); err != nil {
		t.Fatalf("recovery multi-callsite run failed: %v", err)
	}
	if instance.cache.FunctionsRef != recoveredEnv.Functions {
		t.Fatal("expected recovery run to install the recovered multi-callsite FunctionsRef")
	}

	for _, descriptor := range instance.plan.hostCallDescriptors {
		if !instance.cache.HostFunctions[descriptor.ID].Bound {
			t.Fatalf("expected host slot %d (%s) to be rebound after recovery", descriptor.ID, descriptor.DisplayName)
		}
	}
}
