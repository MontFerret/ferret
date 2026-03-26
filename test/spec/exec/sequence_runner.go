package exec

import (
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

func RunSequences(t *testing.T, sequences []spec.Sequence, opts ...vm.EnvironmentOption) {
	t.Helper()

	RunSequencesWith(t, "VM/O0", compiler.New(compiler.WithOptimizationLevel(compiler.O0)), sequences, opts...)
	RunSequencesWith(t, "VM/O1", compiler.New(compiler.WithOptimizationLevel(compiler.O1)), sequences, opts...)
}
