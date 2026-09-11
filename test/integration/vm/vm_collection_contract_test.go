package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestCollectionContracts(t *testing.T) {
	cases := []spec.Spec{
		S(`return {ascending: count(1..3), descending: count(3..1), distinct: count_distinct({first: 7, second: 7.0}), present: includes("123", 123), reversed: reverse([1,2,3])} == {ascending: 3, descending: 3, distinct: 1, present: true, reversed: [3,2,1]}`, true),
		S(`return count(1..3) == 3 and count(3..1) == 3`, true),
		S(`return count_distinct(1..3) == 3 and count_distinct(3..1) == 3`, true),
		S(`return count(@source) == 3 and count_distinct(@source) == 2`, true),
		S(`return includes(@source, 2) and not includes(@source, 3)`, true),
		S(`return count({a: 1, b: 1.0}) == 2 and count_distinct({a: 1, b: 1.0}) == 1`, true),
		S(`return includes({key: 7}, 7) and not includes({key: 7}, "key")`, true),
		S(`return includes("123", 123)`, true),
		S(`return reverse("a狐犬") == "犬狐a" and reverse([1,2,3]) == [3,2,1]`, true),
		S(`return count([]) == 0 and count_distinct({}) == 0 and reverse([]) == []`, true),
		S(`let source = [1,2,3] let result = reverse(source) return source == [1,2,3] and result == [3,2,1]`, true),
	}

	for _, call := range []string{`count("abc")`, `count_distinct("abc")`, `count(1)`, `count_distinct(true)`, `reverse(1..3)`, `includes(1, 1)`} {
		cases = append(cases, S("return "+call+` on error return "caught"`, "caught", call))
	}

	for _, call := range []string{"count()", "count([], [])", "count_distinct()", "count_distinct([], [])", "includes([])", "includes([], 1, 2)", "reverse()", "reverse([], [])"} {
		cases = append(cases, spec.NewSpec("return "+call, call).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "invalid number of arguments"}))
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

func TestReverseHostResultLifecycle(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		source := &ownedCollectionList{Array: runtime.NewArrayWith(runtime.Int(1), runtime.Int(2), runtime.Int(3)), backend: "host storage"}
		RunSpecsWith(t, level.String(), c, []spec.Spec{
			S(`return reverse(@source)`, []any{float64(3), float64(2), float64(1)}),
		}, vm.WithParam("source", source))
		if source.closes != 0 || source.created == nil || source.created.backend != source.backend || source.created.closes != 1 {
			t.Fatalf("result ownership: source closes %d, destination %+v", source.closes, source.created)
		}
	}
}
