package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestPathNamespaceContracts(t *testing.T) {
	specs := []spec.Spec{
		S(`RETURN path::base("a/b.txt")`, "b.txt", "namespaced base"),
		S(`RETURN path::clean("a//b/../c")`, "a/c", "namespaced clean"),
		S(`RETURN path::dir("a/b.txt")`, "a", "namespaced dir"),
		S(`RETURN path::ext("a/b.txt")`, ".txt", "namespaced ext"),
		S(`RETURN path::is_abs("/a")`, true, "namespaced is_abs"),
		Array(`RETURN path::separate("a/b.txt")`, []any{"a/", "b.txt"}, "namespaced separate"),
		S(`RETURN path::match("*.txt", "b.txt")`, true, "namespaced match"),
		S(`RETURN path::join("a", "b.txt")`, "a/b.txt", "namespaced variadic join"),
		S(`RETURN PATH::JOIN(["a", "b.txt"])`, "a/b.txt", "namespaced array join"),
		S(`RETURN PATH::BASE(PATH::JOIN("a", "b.txt"))`, "b.txt", "case-insensitive path namespace"),
		S(`RETURN JOIN(["a", "b"], "|")`, "a|b", "global string join"),
	}
	for _, query := range []string{
		`RETURN BASE("a/b")`, `RETURN CLEAN("a/b")`, `RETURN DIR("a/b")`,
		`RETURN EXT("a.txt")`, `RETURN IS_ABS("/a")`, `RETURN SEPARATE("a/b")`,
		`RETURN MATCH("*", "a")`,
	} {
		specs = append(specs, spec.NewSpec(query, query).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "unresolved function"}))
	}

	for _, query := range []string{
		`RETURN JOIN()`, `RETURN JOIN(["a", "b"])`, `RETURN JOIN("a", "b", "c")`,
		`RETURN PATH::BASE()`, `RETURN PATH::BASE("a", "b")`,
		`RETURN PATH::MATCH("*")`, `RETURN PATH::MATCH("*", "a", "b")`,
	} {
		specs = append(specs, spec.NewSpec(query, query).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "invalid number of arguments"}))
	}

	for _, query := range []string{
		`RETURN JOIN("a", "b")`, `RETURN PATH::BASE(1)`, `RETURN PATH::JOIN("a", 1)`,
	} {
		specs = append(specs, spec.NewSpec(query, query).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "invalid argument type"}))
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		RunSpecsWith(t, level.String(), c, specs)
	}
}
