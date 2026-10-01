package vm

import (
	"context"
	"errors"
	"testing"

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
