package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestUse(t *testing.T) {
	ns := runtime.NewNamespace("FOO")
	ns.Functions().Set0("TEST_FN", func(_ context.Context) (runtime.Value, error) {
		return runtime.True, nil
	})

	RunUseCases(t, []UseCase{
		Case(`
USE FOO AS F

RETURN F::TEST_FN()`, true, "Should compile and resolve alias to the namespaced function using the namespace alias"),
		Case(`
USE FOO::TEST_FN AS FN

RETURN FN()`, true, "Should compile and resolve alias to the namespaced function using the function alias"),
	}, vm.WithNamespace(ns))
}
