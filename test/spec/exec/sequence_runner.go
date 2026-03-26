package exec

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func RunSequencesWith(t *testing.T, level compiler.OptimizationLevel, sequences []spec.Sequence, opts ...vm.EnvironmentOption) {
	t.Helper()

	runner := spec.SequenceRunner{
		Name:     fmt.Sprintf("VM/O%d", level),
		Compiler: compiler.New(compiler.WithOptimizationLevel(level)),
		Env:      opts,
	}

	runner.Run(t, sequences)
}

func RunSequences(t *testing.T, sequences []spec.Sequence, opts ...vm.EnvironmentOption) {
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		RunSequencesWith(t, level, sequences, opts...)
	}
}

func RunSequenceFactory(t *testing.T, factory func() []spec.Sequence, opts ...vm.EnvironmentOption) {
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		RunSequencesWith(t, level, factory(), opts...)
	}
}
