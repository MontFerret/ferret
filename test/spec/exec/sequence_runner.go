package exec

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func RunSequencesWith(t *testing.T, name string, c *compiler.Compiler, sequences []spec.Sequence, opts ...vm.EnvironmentOption) {
	t.Helper()

	runner := spec.SequenceRunner{
		Name:     name,
		Compiler: c,
		Env:      opts,
	}

	runner.Run(t, sequences)
}

func RunSequenceFactory(t *testing.T, factory func() []spec.Sequence, opts ...vm.EnvironmentOption) {
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		RunSequencesWith(t, fmt.Sprintf("VM/O%d", level), compiler.New(compiler.WithOptimizationLevel(level)), factory(), opts...)
	}
}

func RunSequences(t *testing.T, sequences []spec.Sequence, opts ...vm.EnvironmentOption) {
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		RunSequencesWith(t, "VM/O%d", compiler.New(compiler.WithOptimizationLevel(level)), sequences, opts...)
	}
}
