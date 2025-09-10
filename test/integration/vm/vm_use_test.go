package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestUse(t *testing.T) {
	ns := runtime.NewNamespace("FOO")
	ns.Functions().Set0("TEST_FN", func(_ context.Context) (runtime.Value, error) {
		return runtime.True, nil
	})

	RunUseCases(t, []UseCase{
		SkipCase(`
USE FOO

RETURN TEST_FN()`, true, "Should compile but return an error during execution because the object does not implement the interface"),
	}, vm.WithNamespace(ns))
}
