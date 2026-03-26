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
USE Foo AS F

RETURN f::Test_FN()`, "Unresolved function", "Namespace alias resolution is case-sensitive"),
		ErrorStr(`
USE Foo::Test_FN AS Fn

RETURN FN()`, "Unresolved function", "Function alias resolution is case-sensitive"),
		ErrorStr(`
USE Foo AS F

RETURN F::test_fn()`, "Unresolved function", "Host lookup remains case-sensitive after alias expansion"),
	}, vm.WithNamespace(ns))
}
