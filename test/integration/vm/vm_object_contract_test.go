package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestObjectContracts(t *testing.T) {
	specs := []spec.Spec{
		S(`return sorted(object::keys({b: 2, a: 1})) == ["a", "b"]`, true),
		S(`return object::keys({}) == [] and object::values({}) == [] and object::entries({}) == []`, true),
		S(`return sorted(object::values({b: 2, a: 1})) == [1, 2]`, true),
		S(`return object::has_key({a: none}, "a") and not object::has_key({}, "a")`, true),
		S(`return object::keep_keys({a: 1, b: {x: 2}}, "b", "missing", "b") == {b: {x: 2}}`, true),
		S(`return object::keep_keys({a: 1}, []) == {} and object::keep_keys({a: 1}, ["a"]) == {a: 1}`, true),
		S(`return object::omit_keys({a: 1, b: 2}, "a", "missing", "a") == {b: 2}`, true),
		S(`return object::omit_keys({a: 1}, []) == {a: 1} and object::omit_keys({a: 1}, ["a"]) == {}`, true),
		S(`return object::merge([]) == {} and object::merge_deep([]) == {}`, true),
		S(`return object::merge({a: {x: 1}}, {a: {y: 2}}) == {a: {y: 2}}`, true),
		S(`return object::merge_deep({a: {x: 1}}, {a: {y: 2}}) == {a: {x: 1, y: 2}}`, true),
		S(`return object::merge([{a: 1}, {a: 2}]) == {a: 2}`, true),
		S(`return object::merge_deep([{a: {x: 1}}, {a: {y: 2}}]) == {a: {x: 1, y: 2}}`, true),
		S(`return object::merge_deep({a: [1], b: {x: 1}, c: 1}, {a: [2], b: none, c: {y: 2}}) == {a: [2], b: none, c: {y: 2}}`, true),
		S(`return object::zip(["a", "b", "a"], [1, 2, 3]) == {a: 3, b: 2}`, true),
		S(`return object::from_entries([["a", 1], ["b", 2], ["a", 3]]) == {a: 3, b: 2}`, true),
		S(`return object::zip([], []) == {} and object::from_entries([]) == {}`, true),
		S(`let value = {a: [1, none], b: {c: 2}} return object::from_entries(object::entries(value)) == value`, true),
		S(`let value = {a: 1, b: "two", c: none} return (for pair in object::entries(value) return value[pair[0]] == pair[1]) == [true, true, true]`, true),
		S(`let base = {a: {left: 1}, b: [1]} let result = object::merge_deep(base, {a: {right: 2}}) return base == {a: {left: 1}, b: [1]} and result.a == {left: 1, right: 2}`, true),
	}

	for _, call := range []string{
		"object::keys(1)", "object::values(1)", "object::entries(1)", "object::has_key(1, \"a\")", "object::has_key({}, 1)",
		"object::keep_keys(1, \"a\")", "object::keep_keys({}, 1)", "object::omit_keys({}, [1])",
		"object::merge(1)", "object::merge([{}, 1])", "object::merge_deep(1)", "object::merge_deep([{}, 1])",
		"object::zip(1, [])", "object::zip([], 1)", "object::zip([1], [1])", "object::zip([\"a\"], [])",
		"object::from_entries(1)", "object::from_entries([1])", "object::from_entries([[]])",
		"object::from_entries([[\"a\"]])", "object::from_entries([[\"a\", 1, 2]])", "object::from_entries([[1, 2]])",
	} {
		specs = append(specs, S("return "+call+` on error return "caught"`, "caught", call))
	}

	for _, call := range []string{"object::keys()", "object::keys({}, true)", "object::values({}, false)", "object::entries()", "object::has_key({})", "object::keep_keys({})", "object::omit_keys({})", "object::merge()", "object::merge_deep()", "object::zip([])", "object::from_entries()"} {
		specs = append(specs, spec.NewSpec("return "+call, call).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "invalid number of arguments"}))
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		RunSpecsWith(t, level.String(), c, specs)
	}
}
