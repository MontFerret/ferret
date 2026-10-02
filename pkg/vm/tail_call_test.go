package vm

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestTailCallEliminationExecutionAcrossLevels(t *testing.T) {
	cases := []struct {
		want  runtime.Value
		name  string
		query string
	}{
		{name: "chain", query: `func triple(value) => value * 3
func double(value) => triple(value * 2)
return double(3)`, want: runtime.NewInt(18)},
		{name: "immutable captures", query: `let base = 3
func outer(value) {
let offset = value + 1
func target() => base + offset
func forward() => target()
return forward()
}
return outer(2)`, want: runtime.NewInt(6)},
		{name: "argument order and borrowed captures", query: `var total = 0
func next() { total += 1 return total }
func target(a,b) => [a,b,total]
func forward() => target(next(),next())
return [forward(),total]`, want: runtime.NewArrayWith(runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(2)), runtime.NewInt(2))},
		{name: "local cell", query: `func outer() {
var value = 1
func inner() { value += 1 return value }
return inner()
}
return outer()`, want: runtime.NewInt(2)},
		{name: "borrowed cell", query: `var value = 1
func inner() { value += 1 return value }
func forward() => inner()
return [forward(),value]`, want: runtime.NewArrayWith(runtime.NewInt(2), runtime.NewInt(2))},
		{name: "large argument range", query: `func target(a,b,c,d,e,f,g,h,i,j) => a+b+c+d+e+f+g+h+i+j
func forward() => target(1,2,3,4,5,6,7,8,9,10)
return forward()`, want: runtime.NewInt(55)},
		{name: "distinct", query: `func target() => [1,1,2]
func forward() { return distinct target() }
return forward()`, want: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))},
	}

	for _, tc := range cases {
		for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
			t.Run(tc.name+"/"+level.String(), func(t *testing.T) {
				program, err := mustNewCompiler(t, compiler.WithOptimizationLevel(level)).Compile(t.Context(), source.NewAnonymous(tc.query))
				if err != nil {
					t.Fatal(err)
				}

				instance := mustNewVM(t, program)
				t.Cleanup(func() { _ = instance.Close() })
				result := mustRunResult(t, instance, NewDefaultEnvironment())
				t.Cleanup(func() { _ = result.Close() })

				equal, err := runtime.EqualValues(t.Context(), result.Root(), tc.want)
				if err != nil || !equal {
					t.Fatalf("result = %v, want %v: %v", result.Root(), tc.want, err)
				}
			})
		}
	}
}

func TestTailCallEliminationRecoveryAcrossLevels(t *testing.T) {
	stop := errors.New("host failure")
	cases := []struct {
		want      runtime.Value
		name      string
		policy    string
		failures  int
		calls     int
		wantError bool
	}{
		{name: "propagate", failures: 1, calls: 1, want: runtime.None, wantError: true},
		{name: "explicit fail", policy: " on error fail", failures: 1, calls: 1, want: runtime.None, wantError: true},
		{name: "suppress", policy: "?", failures: 1, calls: 1, want: runtime.None},
		{name: "fallback", policy: " on error return 42", failures: 1, calls: 1, want: runtime.NewInt(42)},
		{name: "retry success", policy: " on error retry 2 or return 42", failures: 2, calls: 3, want: runtime.NewInt(9)},
		{name: "retry fallback", policy: " on error retry 2 or return 42", failures: 3, calls: 3, want: runtime.NewInt(42)},
		{name: "retry fail", policy: " on error retry 2 or fail", failures: 3, calls: 3, want: runtime.None, wantError: true},
	}

	for _, tc := range cases {
		for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
			t.Run(tc.name+"/"+level.String(), func(t *testing.T) {
				// Argument recovery finishes before the eligible UDF call. Its work
				// must still run when that call later replaces the caller's frame.
				query := "func target(value) => value\nfunc forward() => target(STEP()" + tc.policy + ")\nreturn forward()"
				program, err := mustNewCompiler(t, compiler.WithOptimizationLevel(level)).Compile(t.Context(), source.NewAnonymous(query))
				if err != nil {
					t.Fatal(err)
				}

				calls := 0
				env := mustNewEnvironment(t, WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					calls++
					if calls <= tc.failures {
						return runtime.None, stop
					}

					return runtime.NewInt(9), nil
				}))
				instance := mustNewVM(t, program)
				t.Cleanup(func() { _ = instance.Close() })
				result, err := instance.Run(t.Context(), env)
				if tc.wantError {
					if !errors.Is(err, stop) {
						t.Fatalf("expected host error, got %v", err)
					}
				} else {
					if err != nil {
						t.Fatal(err)
					}

					if got := mustResultRootAndClose(t, result); got != tc.want {
						t.Fatalf("result = %v, want %v", got, tc.want)
					}
				}

				if calls != tc.calls {
					t.Fatalf("calls = %d, want %d", calls, tc.calls)
				}
			})
		}
	}
}

func TestTailCallEliminationBoundsRecursionAndPreservesRecovery(t *testing.T) {
	const depth = 1024
	const query = `func recur(value) => recur(NEXT(value))
return [recur(1024) on error return none, recur(1024) on error return none]`
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		t.Run(level.String(), func(t *testing.T) {
			program, err := mustNewCompiler(t, compiler.WithOptimizationLevel(level)).Compile(t.Context(), source.NewAnonymous(query))
			if err != nil {
				t.Fatal(err)
			}

			instance := mustNewVM(t, program)
			t.Cleanup(func() { _ = instance.Close() })
			maxDepth, calls := 0, 0
			env := mustNewEnvironment(t, WithFunction("NEXT", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
				calls++
				if current := instance.state.frames.Len(); current > maxDepth {
					maxDepth = current
				}

				if instance.state.frames.NearestRecoveryBoundary() != 0 {
					t.Fatal("lost the main caller's recovery boundary")
				}

				value := args[0].(runtime.Int)
				if value == 0 {
					return runtime.None, errors.New("recursion complete")
				}

				return value - 1, nil
			}))
			result := mustRunResult(t, instance, env)
			equal, err := runtime.EqualValues(t.Context(), result.Root(), runtime.NewArrayWith(runtime.None, runtime.None))
			if err != nil || !equal {
				t.Fatalf("recovered recursion result = %v: %v", result.Root(), err)
			}

			if err := result.Close(); err != nil {
				t.Fatal(err)
			}

			wantDepth := 1
			if level == compiler.None {
				wantDepth = depth + 1
			}

			if maxDepth != wantDepth || calls != 2*(depth+1) || instance.state.frames.Len() != 0 {
				t.Fatalf("depth/calls/remaining frames = %d/%d/%d, want %d/%d/0", maxDepth, calls, instance.state.frames.Len(), wantDepth, 2*(depth+1))
			}
		})
	}
}

func TestTailCallEliminationPreservesResourceAliases(t *testing.T) {
	const query = `func leaf(a,b) { CHECK(a,b) return a }
func forward(value) { let discarded = OPEN() return leaf(value,value) }
return forward(OPEN())`
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		t.Run(level.String(), func(t *testing.T) {
			program, err := mustNewCompiler(t, compiler.WithOptimizationLevel(level)).Compile(t.Context(), source.NewAnonymous(query))
			if err != nil {
				t.Fatal(err)
			}

			var resources []*trackingCloser
			env := mustNewEnvironment(t,
				WithFunction("OPEN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					value := newTrackingCloser("resource")
					resources = append(resources, value)

					return value, nil
				}),
				WithFunction("CHECK", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
					if len(resources) != 2 || args[0] != resources[0] || args[1] != resources[0] || resources[0].closed != 0 || resources[1].closed != 0 {
						t.Fatal("resource aliases or deferred lifetime changed")
					}

					return runtime.None, nil
				}),
			)
			instance := mustNewVM(t, program)
			assertTailCallCandidate(t, instance, level, "forward", "leaf", true)
			t.Cleanup(func() { _ = instance.Close() })
			result := mustRunResult(t, instance, env)
			if result.Root() != resources[0] || resources[0].closed != 0 || resources[1].closed != 0 {
				t.Fatal("returned and deferred resources must survive until result closure")
			}

			if err := result.Close(); err != nil {
				t.Fatal(err)
			}

			if resources[0].closed != 1 || resources[1].closed != 1 {
				t.Fatal("resources must close exactly once")
			}
		})
	}
}

func TestTailCallEliminationTransfersFrameOwnedResourceAliases(t *testing.T) {
	cases := []struct {
		name        string
		grow        bool
		failCleanup bool
	}{
		{name: "reused window"},
		{name: "grown window", grow: true},
		{name: "cleanup error", failCleanup: true},
	}

	for _, tc := range cases {
		for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
			t.Run(tc.name+"/"+level.String(), func(t *testing.T) {
				fixture := newTailCallLifecycleFixture(t, level, tc.grow, "return forward(@borrowed)")
				if tc.failCleanup {
					fixture.closeErr = errors.New("resource cleanup failed")
				}

				result := mustRunResult(t, fixture.instance, fixture.env)
				t.Cleanup(func() { _ = result.Close() })
				fixture.assertProgress(1, 1, 1, 1)
				fixture.assertSettled()
				fixture.assertResult(result, 0)
				fixture.assertResourceCloses(0, 0)

				// Successful execution transfers responsibility to Result. Closing
				// the VM must not shorten either returned or discarded lifetimes.
				fixture.closeVM()
				fixture.assertResourceCloses(0, 0)

				if err := result.Close(); !errors.Is(err, fixture.closeErr) {
					t.Fatalf("result cleanup error = %v, want %v", err, fixture.closeErr)
				}

				if err := result.Close(); err != nil {
					t.Fatalf("repeated result close = %v", err)
				}

				fixture.assertResourceCloses(0, 1)
			})
		}
	}
}

func TestTailCallEliminationFrameOwnedResourcesOnFailure(t *testing.T) {
	for _, failCleanup := range []bool{false, true} {
		name := "execution error"
		if failCleanup {
			name = "execution and cleanup errors"
		}

		for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
			t.Run(name+"/"+level.String(), func(t *testing.T) {
				fixture := newTailCallLifecycleFixture(t, level, true, "return forward(@borrowed)")
				stop := errors.New("tail callee failed")
				fixture.action = func(context.Context) error { return stop }
				if failCleanup {
					fixture.closeErr = errors.New("resource cleanup failed")
				}

				result, err := fixture.instance.Run(t.Context(), fixture.env)
				if result != nil || !errors.Is(err, stop) {
					t.Fatalf("run = %v, %v; want nil result and execution error", result, err)
				}

				// Raw VM failure cleanup preserves the execution error rather than
				// adding cleanup errors or retaining them for a later VM.Close.
				if fixture.closeErr != nil && errors.Is(err, fixture.closeErr) {
					t.Fatalf("cleanup error changed execution error reporting: %v", err)
				}

				fixture.assertProgress(1, 1, 1, 0)
				fixture.assertSettled()
				fixture.assertResourceCloses(0, 1)
				fixture.closeVM()
				fixture.assertResourceCloses(0, 1)
			})
		}
	}
}

func TestTailCallEliminationFrameOwnedResourcesAcrossOuterRecovery(t *testing.T) {
	const main = `let recovered = forward(@borrowed) on error return none
CHECK_RECOVERY(recovered)
return forward(@borrowed)`
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		t.Run(level.String(), func(t *testing.T) {
			fixture := newTailCallLifecycleFixture(t, level, false, main)
			fixture.boundaries = []int{0, -1}
			stop := errors.New("first tail callee failed")
			fixture.action = func(context.Context) error {
				if fixture.actions == 1 {
					return stop
				}

				return nil
			}

			result := mustRunResult(t, fixture.instance, fixture.env)
			t.Cleanup(func() { _ = result.Close() })
			fixture.assertProgress(2, 2, 2, 1)
			if fixture.recoveries != 1 {
				t.Fatalf("recovery callbacks = %d, want 1", fixture.recoveries)
			}

			fixture.assertSettled()
			fixture.assertResult(result, 0)
			fixture.assertResourceCloses(0, 0)
			if err := result.Close(); err != nil {
				t.Fatal(err)
			}

			fixture.assertResourceCloses(0, 1)
			fixture.closeVM()
			fixture.assertResourceCloses(0, 1)
		})
	}
}

func TestTailCallEliminationFrameOwnedResourcesOnCancellationAndReuse(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		t.Run(level.String(), func(t *testing.T) {
			fixture := newTailCallLifecycleFixture(t, level, false, "return forward(@borrowed)")
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			fixture.entryAction = func() {
				if fixture.entries == 1 {
					cancel()
				}
			}

			result, err := fixture.instance.Run(ctx, fixture.env)
			if result != nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("canceled run = %v, %v", result, err)
			}

			// ENTRY cancels after proving replacement and transfer. The next
			// host call is an existing safepoint and must not run.
			fixture.assertProgress(1, 0, 0, 0)
			fixture.assertSettled()
			fixture.assertResourceCloses(0, 1)

			firstRunResources := len(fixture.resources)
			result = mustRunResult(t, fixture.instance, fixture.env)
			t.Cleanup(func() { _ = result.Close() })
			fixture.assertProgress(2, 1, 1, 1)
			fixture.assertSettled()
			fixture.assertResult(result, firstRunResources)
			fixture.assertResourceCloses(firstRunResources, 0)
			if err := result.Close(); err != nil {
				t.Fatal(err)
			}

			fixture.assertResourceCloses(0, 1)
			fixture.closeVM()
			fixture.assertResourceCloses(0, 1)
		})
	}
}

func TestTailCallEliminationReusesWindowAcrossRuns(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		t.Run(level.String(), func(t *testing.T) {
			fixture := newTailCallLifecycleFixture(t, level, false, "return forward(@borrowed)")
			var window *runtime.Value
			for run := range 3 {
				start := len(fixture.resources)
				result := mustRunResult(t, fixture.instance, fixture.env)
				t.Cleanup(func() { _ = result.Close() })
				fixture.assertSettled()
				fixture.assertResult(result, start)
				if run == 0 {
					window = fixture.window
				} else if fixture.window != window {
					t.Fatal("repeated execution allocated another forward window after tail-call reuse")
				}

				if err := result.Close(); err != nil {
					t.Fatal(err)
				}
				fixture.assertResourceCloses(0, 1)
			}
			fixture.assertProgress(3, 3, 3, 3)
			fixture.closeVM()
			fixture.assertResourceCloses(0, 1)
		})
	}
}

func tailCallLifecycleQuery(grow bool, main string) string {
	// Context-aware host calls keep all 32 arguments live, including when
	// register coalescing runs. Runtime assertions still verify the actual path.
	values := make([]string, 32)
	for i := range values {
		values[i] = "VALUE(" + strconv.Itoa(i) + ")"
	}
	wide := "PAD(" + strings.Join(values, ",") + ")\n"
	forwardWork, leafWork := wide, ""
	if grow {
		forwardWork, leafWork = "", wide
	}

	return "func leaf(a,b,borrowed) {\nENTRY()\n" + leafWork + `CHECK_SURVIVOR(b,borrowed)
STEP()
AFTER_STEP()
return b
}
func forward(borrowed) {
let resource = OPEN()
let discarded = OPEN()
` + forwardWork + `SNAPSHOT()
return leaf(resource,resource,borrowed)
}
` + main
}

func assertTailCallCandidate(t testing.TB, instance *VM, level compiler.OptimizationLevel, callerName, calleeName string, eligible bool) callDescriptor {
	t.Helper()

	callerID, calleeID := bytecode.NoFunction, bytecode.NoFunction
	for id, udf := range instance.program.Functions.UserDefined {
		if udf.DisplayName == callerName {
			callerID = bytecode.FunctionID(id)
		}

		if udf.DisplayName == calleeName {
			calleeID = bytecode.FunctionID(id)
		}
	}

	if !callerID.Valid() || !calleeID.Valid() {
		t.Fatalf("missing candidate functions %q and %q", callerName, calleeName)
	}

	start := instance.program.Functions.UserDefined[callerID].Entry
	end := len(instance.program.Bytecode)
	for _, udf := range instance.program.Functions.UserDefined {
		if udf.Entry > start && udf.Entry < end {
			end = udf.Entry
		}
	}

	want := bytecode.OpCall
	descriptors := instance.plan.udfCallDescriptors
	if eligible && level != compiler.None {
		want = bytecode.OpTailCall
		descriptors = instance.plan.udfTailCallDescriptors
	}

	var candidate callDescriptor
	count := 0
	for _, descriptor := range descriptors {
		if descriptor.ID == calleeID && descriptor.PC >= start && descriptor.PC < end {
			candidate = descriptor
			count++
		}
	}

	if count != 1 {
		t.Fatalf("%s -> %s candidate calls with opcode %v = %d, want 1", callerName, calleeName, want, count)
	}

	call := instance.program.Bytecode[candidate.PC]
	if call.Opcode != want {
		t.Fatalf("candidate opcode = %v, want %v", call.Opcode, want)
	}

	if eligible {
		if candidate.PC+1 >= end {
			t.Fatal("candidate has no retained return")
		}

		ret := instance.program.Bytecode[candidate.PC+1]
		if ret.Opcode != bytecode.OpReturn || ret.Operands[0] != candidate.Dst {
			t.Fatalf("candidate return = %v, want unchanged call result", ret)
		}
	}

	return candidate
}
