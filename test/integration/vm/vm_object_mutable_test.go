package vm_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestMutableObjectContracts(t *testing.T) {
	specs := []spec.Spec{
		S(`let obj = {a: 1} let alias = obj let result = object::mut::merge(obj, {b: 2})
		let again = object::mut::merge(result, {c: 3})
		return obj == {a: 1, b: 2, c: 3} and alias == obj and again == obj`, true),
		S(`let obj = {a: 1, internal: true} let alias = obj
		let result = object::mut::omit_keys(object::mut::merge(obj, {b: 2}), "internal")
		return obj == {a: 1, b: 2} and alias == obj and result == obj`, true),
		S(`let obj = {a: 1, b: 2, c: 3} let alias = obj
		let result = object::mut::keep_keys(obj, "a", "c", "missing", "a")
		return obj == {a: 1, c: 3} and result == obj and alias == obj`, true),
		S(`let obj = {a: 1, b: 2} let alias = obj
		let result = object::mut::omit_keys(obj, ["a", "missing", "a"])
		return obj == {b: 2} and result == obj and alias == obj`, true),
		S(`let child = {leaf: {left: 1}} let obj = {change: child, retain: child} let alias = obj
		let source = {change: {leaf: {right: 2}}, borrowed: child}
		let result = object::mut::merge_deep(obj, source)
		return obj.change == {leaf: {left: 1, right: 2}} and child == {leaf: {left: 1}}
		and obj.retain == child and source == {change: {leaf: {right: 2}}, borrowed: {leaf: {left: 1}}}
		and alias == obj and result == obj`, true),
		S(`let obj = {a: {left: 1}, b: [1], c: 1}
		return object::mut::merge_deep(obj, [{a: {right: 2}}, {b: [2], c: none}]) == {a: {left: 1, right: 2}, b: [2], c: none}`, true),
		S(`let obj = {a: 1} return object::mut::merge(obj, [{a: 2}, {a: 3}]) == {a: 3}`, true),
		S(`let obj = {a: 1} let copy = object::omit_keys(obj, "a")
		let result = object::mut::omit_keys(obj, "a") return obj == {} and copy == {} and result == {}`, true),
		S(`let obj = {a: 1} return object::mut::keep_keys(obj, []) == {} and obj == {}`, true),
		S(`let obj = {a: 1} return object::mut::omit_keys(obj, []) == {a: 1} and obj == {a: 1}`, true),
	}

	for _, name := range []string{"merge", "merge_deep"} {
		for _, sources := range []string{"", ", []", ", {}"} {
			specs = append(specs, S(fmt.Sprintf(`let obj = {a: 1} let result = object::mut::%s(obj%s)
			let changed = object::mut::merge(result, {b: 2}) return obj == {a: 1, b: 2} and changed == obj`, name, sources), true))
		}

		for _, source := range []string{`{nested: {left: 1}}`, `[{nested: {left: 1}}]`} {
			specs = append(specs, S(fmt.Sprintf(`let source = %s let obj = {}
			let result = object::mut::%s(obj, source)
			let nested = object::mut::merge(result.nested, {right: 2})
			return source == %s and nested == {left: 1, right: 2}`, source, name, source), true))
		}
	}

	for _, call := range []string{
		`object::mut::merge(1)`, `object::mut::merge([])`, `object::mut::merge({}, [{}, 1])`,
		`object::mut::merge_deep(1)`, `object::mut::merge_deep({}, {}, 1)`,
		`object::mut::keep_keys({}, ["a", 1])`, `object::mut::keep_keys(1, "a")`,
		`object::mut::omit_keys({}, "a", 1)`, `object::mut::omit_keys(1, "a")`,
	} {
		specs = append(specs, S("return "+call+` on error return "caught"`, "caught", call))
	}

	for _, call := range []string{"object::mut::merge()", "object::mut::merge_deep()", "object::mut::keep_keys({})", "object::mut::omit_keys({})"} {
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
