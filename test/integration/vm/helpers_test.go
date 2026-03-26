package vm_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	specexec "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func runSpecsLevels(t *testing.T, factory func() []spec.Spec) {
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		specexec.RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			factory(),
		)
	}
}

func runSequencesLevels(t *testing.T, factory func() []spec.Sequence) {
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		specexec.RunSequencesWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			factory(),
		)
	}
}

func runProgramSpecs(t *testing.T, specs []spec.Spec) {
	t.Helper()

	spec.NewRunner("VM/Program", compiler.WithOptimizationLevel(compiler.O0)).Run(t, specs)
}
