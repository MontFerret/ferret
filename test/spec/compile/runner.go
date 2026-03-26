package compile

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func RunSpecsWith(t *testing.T, name string, c *compiler.Compiler, specs []spec.Spec) {
	t.Helper()

	runner := &spec.Runner{
		Name:     name,
		Compiler: c,
	}

	runner.Run(t, specs)
}

func RunSpecs(t *testing.T, useCases []spec.Spec) {
	t.Helper()
	RunSpecsWith(t, "Compiler/O0", compiler.New(compiler.WithOptimizationLevel(compiler.O0)), useCases)
}

func RunSpecsLevels(t *testing.T, specs []spec.Spec, levels ...compiler.OptimizationLevel) {
	t.Helper()

	if len(levels) == 0 {
		levels = []compiler.OptimizationLevel{compiler.O0}
	}

	for _, level := range levels {
		RunSpecsWith(
			t,
			fmt.Sprintf("Compiler/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			specs,
		)
	}
}
