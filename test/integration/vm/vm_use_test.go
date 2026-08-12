package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestUse(t *testing.T) {
	ns := runtime.NewNamespace("Foo")
	ns.Function().A0().Add("Test_FN", func(_ context.Context) (runtime.Value, error) {
		return runtime.True, nil
	})

	RunSpecs(t, []spec.Spec{
		S(`
USE Foo AS F

RETURN F::Test_FN()`, true, "Should compile and resolve alias to the namespaced function using the namespace alias"),
		S(`
USE Foo AS F
FUNC f() => F::Test_FN()
RETURN f()`, true, "Should resolve namespace alias host call inside UDF body"),
		S(`
USE Foo AS F
FUNC f() => true
RETURN f()`, true, "Should not rewrite bare UDF call through namespace alias"),
		S(`
USE Foo::Test_FN AS Fn

RETURN Fn()`, true, "Should compile and resolve alias to the namespaced function using the function alias"),
		ErrorStr(`
USE Foo AS Alias

RETURN alias::Test_FN()`, "unresolved function", "Namespace alias names remain case-sensitive"),
		ErrorStr(`
USE Foo::Test_FN AS Fn

RETURN FN()`, "unresolved function", "Function alias resolution is case-sensitive"),
		S(`
USE Foo AS Alias

RETURN Alias::tEsT_fN()`, true, "Qualified host lookup is case-insensitive after exact namespace alias expansion"),
		S(`
USE fOO AS Alias

RETURN Alias::TEST_FN()`, true, "Namespace alias target casing is canonicalized"),
	}, vm.WithNamespace(ns))
}
