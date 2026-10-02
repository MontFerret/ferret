package vm

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type (
	tailCallLifecycleResource struct {
		closeErr error
		name     string
		closed   int
	}

	tailCallLifecycleFixture struct {
		t            *testing.T
		instance     *VM
		env          *Environment
		borrowed     *trackingCloser
		resources    []*tailCallLifecycleResource
		action       func(context.Context) error
		entryAction  func()
		closeErr     error
		window       *runtime.Value
		boundaries   []int
		candidate    callDescriptor
		forwardID    bytecode.FunctionID
		level        compiler.OptimizationLevel
		capacity     int
		depth        int
		snapshots    int
		entries      int
		survivors    int
		actions      int
		afterActions int
		recoveries   int
		grow         bool
	}
)

func newTailCallLifecycleFixture(t *testing.T, level compiler.OptimizationLevel, grow bool, main string) *tailCallLifecycleFixture {
	t.Helper()

	program, err := mustNewCompiler(t, compiler.WithOptimizationLevel(level)).Compile(t.Context(), source.NewAnonymous(tailCallLifecycleQuery(grow, main)))
	if err != nil {
		t.Fatal(err)
	}

	fixture := &tailCallLifecycleFixture{
		t:        t,
		instance: mustNewVM(t, program),
		borrowed: newTrackingCloser("borrowed"),
		level:    level,
		grow:     grow,
	}
	fixture.candidate = assertTailCallCandidate(t, fixture.instance, level, "forward", "leaf", true)
	for id, udf := range program.Functions.UserDefined {
		if udf.DisplayName == "forward" {
			fixture.forwardID = bytecode.FunctionID(id)
		}
	}

	fixture.env = mustNewEnvironment(t,
		WithParam("borrowed", fixture.borrowed),
		WithFunction("OPEN", fixture.open),
		WithFunction("SNAPSHOT", fixture.snapshot),
		WithFunction("ENTRY", fixture.entry),
		WithFunction("CHECK_SURVIVOR", fixture.checkSurvivor),
		WithFunction("STEP", fixture.step),
		WithFunction("AFTER_STEP", fixture.afterStep),
		WithFunction("CHECK_RECOVERY", fixture.checkRecovery),
		WithFunction("VALUE", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
			if len(args) != 1 {
				t.Fatalf("VALUE argument count = %d, want 1", len(args))
			}

			return args[0], nil
		}),
		WithFunction("PAD", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
			if len(args) != 32 {
				t.Fatalf("wide host argument count = %d, want 32", len(args))
			}

			for i, value := range args {
				if value != runtime.NewInt(i) {
					t.Fatalf("wide host argument %d = %v, want %d", i, value, i)
				}
			}

			return runtime.None, nil
		}),
	)
	t.Cleanup(func() { _ = fixture.instance.Close() })

	return fixture
}

func (r *tailCallLifecycleResource) Close() error {
	r.closed++

	return r.closeErr
}

func (r *tailCallLifecycleResource) String() string {
	return r.name
}

func (r *tailCallLifecycleResource) Hash() uint64 {
	return 0
}

func (r *tailCallLifecycleResource) Copy() runtime.Value {
	return r
}

func (f *tailCallLifecycleFixture) open(context.Context, ...runtime.Value) (runtime.Value, error) {
	resource := &tailCallLifecycleResource{name: fmt.Sprintf("resource-%d", len(f.resources))}
	if len(f.resources)%2 == 0 {
		resource.closeErr = f.closeErr
	}

	f.resources = append(f.resources, resource)

	return resource, nil
}

func (f *tailCallLifecycleFixture) snapshot(context.Context, ...runtime.Value) (runtime.Value, error) {
	f.snapshots++
	if len(f.resources) != 2*f.snapshots {
		f.t.Fatalf("resources before candidate = %d, want %d", len(f.resources), 2*f.snapshots)
	}

	state := &f.instance.state
	resource := f.resources[len(f.resources)-2]
	if state.frames.CurrentFunctionID() != f.forwardID || !state.owned.Owns(resource) {
		f.t.Fatal("forward must own its produced resource before the candidate call")
	}

	if state.owned.Owns(f.borrowed) || f.hasDeferred(f.borrowed) {
		f.t.Fatal("forward adopted the borrowed host parameter")
	}

	f.window = &state.registers[0]
	f.capacity = cap(state.registers)
	f.depth = state.frames.Len()
	leafRegisters := f.instance.program.Functions.UserDefined[f.candidate.ID].Registers
	if (f.capacity < leafRegisters) != f.grow {
		f.t.Fatalf("candidate window capacity = %d, leaf registers = %d; growth = %v, want %v", f.capacity, leafRegisters, f.capacity < leafRegisters, f.grow)
	}

	return runtime.None, nil
}

func (f *tailCallLifecycleFixture) entry(context.Context, ...runtime.Value) (runtime.Value, error) {
	f.entries++
	if f.snapshots != f.entries {
		f.t.Fatal("leaf entry did not execute the snapshotted candidate")
	}

	state := &f.instance.state
	resource, discarded := f.resources[len(f.resources)-2], f.resources[len(f.resources)-1]
	wantDepth := f.depth
	if f.level == compiler.None {
		wantDepth++
	}

	if state.frames.Len() != wantDepth || state.frames.CurrentFunctionID() != f.candidate.ID || state.frames.Top().CallSitePC != f.candidate.CallSitePC {
		f.t.Fatal("leaf entry did not preserve the expected candidate frame structure")
	}

	wantBoundary := -1
	if len(f.boundaries) > 0 {
		if f.entries > len(f.boundaries) {
			f.t.Fatal("unexpected leaf entry after recovery")
		}

		wantBoundary = f.boundaries[f.entries-1]
	}

	if got := state.frames.NearestRecoveryBoundary(); got != wantBoundary {
		f.t.Fatalf("leaf recovery boundary = %d, want %d", got, wantBoundary)
	}

	if state.registers[1] != resource || state.registers[2] != resource || state.registers[3] != f.borrowed {
		f.t.Fatal("leaf argument identities changed")
	}

	if resource.closed != 0 || discarded.closed != 0 || f.borrowed.closed != 0 {
		f.t.Fatal("leaf arguments or discarded resource closed before their owner settled")
	}

	key, _, ok := mem.ResourceKeyOf(resource)
	if !ok {
		f.t.Fatal("resource fixture is not trackable")
	}

	// ENTRY takes no arguments, so no host-argument moves have added aliases
	// before these checks. UDF parameter registers remain pinned under Full.
	if got := f.aliasCount(state.registers, key); got != 2 {
		f.t.Fatalf("direct leaf argument aliases = %d, want 2", got)
	}

	if f.level == compiler.None {
		if state.owned.Owns(resource) || state.aliases.Count(key) != 0 || !state.frames.Top().OwnedResources.Owns(resource) || f.hasDeferred(resource) {
			f.t.Fatal("ordinary leaf call must borrow the resource from its live forward frame")
		}
	} else {
		if (&state.registers[0] != f.window) != f.grow {
			f.t.Fatalf("tail replacement installed a new register window = %v, want %v", &state.registers[0] != f.window, f.grow)
		}

		if !state.owned.Owns(resource) || state.aliases.Count(key) != 2 || f.hasDeferred(resource) {
			f.t.Fatal("tail replacement must transfer ownership and both argument aliases into leaf")
		}

		if state.owned.Owns(discarded) || !f.hasDeferred(discarded) {
			f.t.Fatal("discarded forward resource must follow deferred cleanup after replacement")
		}
	}

	if state.owned.Owns(f.borrowed) || f.hasDeferred(f.borrowed) {
		f.t.Fatal("tail replacement adopted the borrowed host parameter")
	}

	// The first parameter has no later source-level use. Clear its actual
	// register through the normal ownership path, then let FQL use and return b.
	state.clearRegister(bytecode.NewRegister(1))
	if state.registers[2] != resource || resource.closed != 0 || f.aliasCount(state.registers, key) != 1 {
		f.t.Fatal("clearing one argument alias invalidated the surviving alias")
	}

	if f.level != compiler.None && (!state.owned.Owns(resource) || state.aliases.Count(key) != 1 || f.hasDeferred(resource)) {
		f.t.Fatal("leaf lost ownership while its second argument alias remained live")
	}

	if f.entryAction != nil {
		f.entryAction()
	}

	return runtime.None, nil
}

func (f *tailCallLifecycleFixture) checkSurvivor(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	f.survivors++
	resource := f.resources[len(f.resources)-2]
	if len(args) != 2 || args[0] != resource || args[1] != f.borrowed || resource.closed != 0 || f.borrowed.closed != 0 {
		f.t.Fatal("FQL could not use the surviving resource and borrowed parameter")
	}

	if f.level != compiler.None {
		state := &f.instance.state
		key, _, _ := mem.ResourceKeyOf(resource)
		if !state.owned.Owns(resource) || state.aliases.Count(key) != f.aliasCount(state.registers, key) || f.hasDeferred(resource) {
			f.t.Fatal("surviving resource escaped active leaf ownership")
		}
	}

	return runtime.None, nil
}

func (f *tailCallLifecycleFixture) step(ctx context.Context, _ ...runtime.Value) (runtime.Value, error) {
	f.actions++
	if f.action != nil {
		return runtime.None, f.action(ctx)
	}

	return runtime.None, nil
}

func (f *tailCallLifecycleFixture) afterStep(context.Context, ...runtime.Value) (runtime.Value, error) {
	f.afterActions++

	return runtime.None, nil
}

func (f *tailCallLifecycleFixture) checkRecovery(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	f.recoveries++
	state := &f.instance.state
	if len(args) != 1 || args[0] != runtime.None || f.entries != 1 || state.frames.Len() != 0 || state.frames.NearestRecoveryBoundary() != -1 || state.hasFail {
		f.t.Fatal("outer recovery did not resume in main with its fallback and cleared boundary")
	}

	for _, resource := range f.resources {
		if resource.closed != 0 || !f.hasDeferred(resource) {
			f.t.Fatal("recovered frame resources must remain deferred until result closure")
		}
	}

	return runtime.None, nil
}

func (f *tailCallLifecycleFixture) assertProgress(entries, survivors, actions, afterActions int) {
	f.t.Helper()

	if f.snapshots != entries || f.entries != entries || f.survivors != survivors || f.actions != actions || f.afterActions != afterActions {
		f.t.Fatalf("snapshot/entry/survivor/action/after calls = %d/%d/%d/%d/%d, want %d/%d/%d/%d/%d", f.snapshots, f.entries, f.survivors, f.actions, f.afterActions, entries, entries, survivors, actions, afterActions)
	}
}

func (f *tailCallLifecycleFixture) assertSettled() {
	f.t.Helper()

	state := &f.instance.state
	if state.frames.Len() != 0 || state.frames.NearestRecoveryBoundary() != -1 || !state.owned.Empty() || !state.deferred.Empty() || state.env != nil || state.hasFail || state.hasCaught || len(state.cellHandles) != 0 {
		f.t.Fatal("run left stale frames, ownership, recovery, or execution state")
	}

	for _, value := range state.registers {
		if value != runtime.None {
			f.t.Fatal("settled register window retained an execution value")
		}
	}

	for _, resource := range f.resources {
		key, _, _ := mem.ResourceKeyOf(resource)
		if state.aliases.Count(key) != 0 {
			f.t.Fatal("run left stale resource alias counts")
		}
	}
}

func (f *tailCallLifecycleFixture) assertResult(result *Result, start int) {
	f.t.Helper()

	if result.Root() != f.resources[len(f.resources)-2] || result.set.Len() != len(f.resources)-start {
		f.t.Fatal("result did not acquire the returned and deferred resources")
	}

	for _, resource := range f.resources[start:] {
		found := false
		result.set.ForEach(func(closer io.Closer) {
			found = found || closer == resource
			if closer == f.borrowed {
				f.t.Fatal("result adopted the borrowed host parameter")
			}
		})
		if !found {
			f.t.Fatalf("result omitted %s", resource)
		}
	}
}

func (f *tailCallLifecycleFixture) assertResourceCloses(start, want int) {
	f.t.Helper()

	for _, resource := range f.resources[start:] {
		if resource.closed != want {
			f.t.Fatalf("%s close count = %d, want %d", resource, resource.closed, want)
		}
	}

	if f.borrowed.closed != 0 {
		f.t.Fatalf("borrowed host parameter closed %d times", f.borrowed.closed)
	}
}

func (f *tailCallLifecycleFixture) closeVM() {
	f.t.Helper()

	for range 2 {
		if err := f.instance.Close(); err != nil {
			f.t.Fatalf("VM close = %v", err)
		}
	}
}

func (f *tailCallLifecycleFixture) hasDeferred(value io.Closer) bool {
	found := false
	f.instance.state.deferred.ForEach(func(closer io.Closer) {
		found = found || closer == value
	})

	return found
}

func (f *tailCallLifecycleFixture) aliasCount(registers []runtime.Value, key mem.ResourceKey) int {
	count := 0
	for _, value := range registers {
		if aliasKey, _, ok := mem.ResourceKeyOf(value); ok && aliasKey == key {
			count++
		}
	}

	return count
}
