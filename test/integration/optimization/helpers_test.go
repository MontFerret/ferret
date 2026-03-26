package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func runProgramSpecs(t *testing.T, specs []spec.Spec) {
	t.Helper()

	spec.NewRunner("optimization/program", compiler.WithOptimizationLevel(compiler.O0)).Run(t, specs)
}
