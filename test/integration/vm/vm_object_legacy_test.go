package vm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestObjectLegacyAliases(t *testing.T) {
	var specs []spec.Spec
	for _, test := range []struct {
		legacy    string
		canonical string
		args      string
		sort      bool
	}{
		{"keys", "object::keys", `{b: 2, a: 1}`, true},
		{"values", "object::values", `{b: 2, a: 1}`, true},
		{"has", "object::has_key", `{a: none}, "a"`, false},
		{"has", "object::has_key", `{}, "missing"`, false},
		{"keep_keys", "object::keep_keys", `{a: 1, b: {x: 2}}, "b", "missing", "b"`, false},
		{"keep_keys", "object::keep_keys", `{a: 1, b: 2}, ["b"]`, false},
		{"keep_keys", "object::keep_keys", `{a: 1}, []`, false},
		{"merge", "object::merge", `{a: {x: 1}}, {a: {y: 2}}`, false},
		{"merge", "object::merge", `[{a: 1}, {a: 2}]`, false},
		{"merge", "object::merge", `[]`, false},
		{"merge_recursive", "object::merge_deep", `{a: {x: 1}}, {a: {y: 2}}`, false},
		{"merge_recursive", "object::merge_deep", `[{a: {x: 1}}, {a: {y: 2}}]`, false},
		{"merge_recursive", "object::merge_deep", `[]`, false},
		{"zip", "object::zip", `["a", "b", "a"], [1, 2, 3]`, false},
		{"zip", "object::zip", `[], []`, false},
	} {
		for _, name := range []string{test.legacy, strings.ToUpper(test.legacy)} {
			legacy := fmt.Sprintf("%s(%s)", name, test.args)
			canonical := fmt.Sprintf("%s(%s)", test.canonical, test.args)
			if test.sort {
				legacy = "arrays::sorted(" + legacy + ")"
				canonical = "arrays::sorted(" + canonical + ")"
			}

			specs = append(specs, S("return "+legacy+" == "+canonical, true))
		}
	}

	specs = append(specs,
		S(`return zip(["a", "b", "a"], [1, 2, 3]) == {a: 3, b: 2}`, true),
		S(`let value = {a: {x: 1}, b: 2} let result = keep_keys(value, "a") return value == {a: {x: 1}, b: 2} and result == {a: {x: 1}}`, true),
		S(`let value = {a: {x: 1}} let result = merge(value, {a: {y: 2}}) return value == {a: {x: 1}} and result == {a: {y: 2}}`, true),
		S(`let value = {a: {x: 1}} let result = merge_recursive(value, {a: {y: 2}}) return value == {a: {x: 1}} and result == {a: {x: 1, y: 2}}`, true),
		S(`let value = {x: 1} let result = zip(["a"], [value]) let changed = object::mut::merge(result.a, {y: 2}) return value == {x: 1} and changed == {x: 1, y: 2}`, true),
	)

	for _, call := range []string{
		`keys(1)`, `values(1)`, `has({}, 1)`, `keep_keys({}, 1)`,
		`merge([{}, 1])`, `merge_recursive([{}, 1])`, `zip(["a"], [])`,
	} {
		specs = append(specs, S("return "+call+` on error return "caught"`, "caught"))
	}

	for _, call := range []string{"keys()", "keys({}, true)", "values({}, false)", "has({})", "keep_keys({})", "merge()", "merge_recursive()", "zip([])"} {
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
