package compile

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func RunSpecsWith(t *testing.T, level compiler.OptimizationLevel, specs []spec.Spec) {
	t.Helper()

	runner := &spec.Runner{
		Name:     fmt.Sprintf("Compiler/O%d", level),
		Compiler: mustNewCompiler(t, compiler.WithOptimizationLevel(level)),
	}

	runner.Run(t, specs)
}

func RunSpecs(t *testing.T, useCases []spec.Spec) {
	t.Helper()

	RunSpecsWith(t, compiler.OptimizationNone, useCases)
}

func RunSpecsLevels(t *testing.T, specs []spec.Spec, levels ...compiler.OptimizationLevel) {
	t.Helper()

	if len(levels) == 0 {
		levels = []compiler.OptimizationLevel{compiler.OptimizationNone}
	}

	for _, level := range levels {
		RunSpecsWith(
			t,
			level,
			specs,
		)
	}
}
