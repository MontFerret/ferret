package vm_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestCollectionContracts(t *testing.T) {
	cases := []spec.Spec{
		S(`return {ascending: collections::count(1..3), descending: collections::count(3..1), distinct: collections::count_distinct({first: 7, second: 7.0}), present: collections::includes("123", 123), reversed: collections::reverse([1,2,3])} == {ascending: 3, descending: 3, distinct: 1, present: true, reversed: [3,2,1]}`, true),
		S(`return collections::count(1..3) == 3 and collections::count(3..1) == 3`, true),
		S(`return collections::count_distinct(1..3) == 3 and collections::count_distinct(3..1) == 3`, true),
		S(`return collections::count(@source) == 3 and collections::count_distinct(@source) == 2`, true),
		S(`return collections::includes(@source, 2) and not collections::includes(@source, 3)`, true),
		S(`return collections::count({a: 1, b: 1.0}) == 2 and collections::count_distinct({a: 1, b: 1.0}) == 1`, true),
		S(`return collections::includes({key: 7}, 7) and not collections::includes({key: 7}, "key")`, true),
		S(`return collections::includes("123", 123)`, true),
		S(`return collections::reverse("a狐犬") == "犬狐a" and collections::reverse([1,2,3]) == [3,2,1]`, true),
		S(`return collections::count([]) == 0 and collections::count_distinct({}) == 0 and collections::reverse([]) == []`, true),
		S(`let source = [1,2,3] let result = collections::reverse(source) return source == [1,2,3] and result == [3,2,1]`, true),
	}

	for _, call := range []string{`collections::count("abc")`, `collections::count_distinct("abc")`, `collections::count(1)`, `collections::count_distinct(true)`, `collections::reverse(1..3)`, `collections::includes(1, 1)`} {
		cases = append(cases, S("return "+call+` on error return "caught"`, "caught", call))
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		source := &collectionIterable{values: runtime.NewArrayWith(runtime.Int(1), runtime.Int(2), runtime.Float(2))}
		RunSpecsWith(t, level.String(), c, cases, vm.WithParam("source", source))
	}
}

func TestCollectionCompatibility(t *testing.T) {
	var cases []spec.Spec
	for _, prefix := range []string{"collections::", ""} {
		for _, tc := range []struct {
			want any
			call string
		}{
			{3, "count(1..3)"},
			{1, "count_distinct([7, 7.0])"},
			{true, `includes("123", 123)`},
			{"犬狐a", `reverse("a狐犬")`},
		} {
			cases = append(cases, S("return "+prefix+tc.call, tc.want))
		}

		for _, call := range []string{"count()", "count([], [])", "count_distinct()", "count_distinct([], [])", "includes([])", "includes([], 1, 2)", "reverse()", "reverse([], [])"} {
			cases = append(cases, spec.NewSpec("return "+prefix+call).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "invalid number of arguments"}))
		}
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		RunSpecsWith(t, level.String(), c, cases)
	}
}

func TestMeasuredCollectionCountContracts(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		for _, length := range []runtime.Int{-1, math.MinInt64, 0, 42} {
			t.Run(fmt.Sprintf("%s/length=%d", level, length), func(t *testing.T) {
				c, err := compiler.New(compiler.WithOptimizationLevel(level))
				if err != nil {
					t.Fatal(err)
				}

				source := &measuredCollectionIterable{collectionIterable: &collectionIterable{values: runtime.NewArray(0)}, length: length}
				test := S(`return collections::count(@source)`, float64(length))
				if length < 0 {
					test = spec.NewSpec(`return collections::count(@source)`).Expect().ExecError(func(t *testing.T, actual any, _ ...any) {
						t.Helper()

						err, ok := actual.(error)
						if !ok || !errors.Is(err, runtime.ErrInvalidOperation) {
							t.Fatalf("negative measurement error lost: %v", actual)
						}
					})
				}

				RunSpecsWith(t, level.String(), c, []spec.Spec{test}, vm.WithParam("source", source))
				if source.measurements != 1 || source.iterations != 0 {
					t.Fatalf("unexpected host work: %+v", source)
				}
			})
		}
	}
}

func TestImmutableMergeHostResultLifecycle(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		for _, function := range []string{"object::merge", "object::merge_deep"} {
			for _, fail := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%s/fail=%v", level, function, fail), func(t *testing.T) {
					c, err := compiler.New(compiler.WithOptimizationLevel(level))
					if err != nil {
						t.Fatal(err)
					}

					primary, cleanup := errors.New("host population failed"), errors.New("host cleanup failed")
					source := &ownedCollectionMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.Int(1)}), backend: "host storage"}
					query := "return " + function + "(@source)"
					test := S(query, map[string]any{"key": float64(1)})
					if fail {
						source.writeErr, source.closeErr = primary, cleanup
						test = spec.NewSpec(query).Expect().ExecError(func(t *testing.T, actual any, _ ...any) {
							t.Helper()

							err, ok := actual.(error)
							if !ok || !errors.Is(err, primary) || !errors.Is(err, cleanup) {
								t.Fatalf("host error identities lost: %v (primary=%v, cleanup=%v)", actual, errors.Is(err, primary), errors.Is(err, cleanup))
							}
						})
					}

					RunSpecsWith(t, level.String(), c, []spec.Spec{test}, vm.WithParam("source", source))
					if source.closes != 0 || source.created == nil || source.created.closes != 1 || source.created.backend != source.backend || source.String() != `{"key":1}` {
						t.Fatalf("result ownership: source closes %d, destination %+v", source.closes, source.created)
					}
				})
			}
		}
	}
}

func TestReverseHostResultLifecycle(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		source := &ownedCollectionList{Array: runtime.NewArrayWith(runtime.Int(1), runtime.Int(2), runtime.Int(3)), backend: "host storage"}
		RunSpecsWith(t, level.String(), c, []spec.Spec{
			S(`return collections::reverse(@source)`, []any{float64(3), float64(2), float64(1)}),
		}, vm.WithParam("source", source))
		if source.closes != 0 || source.created == nil || source.created.backend != source.backend || source.created.closes != 1 {
			t.Fatalf("result ownership: source closes %d, destination %+v", source.closes, source.created)
		}
	}
}
