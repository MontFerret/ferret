package exec

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func RunSpecsWith(t *testing.T, name string, c *compiler.Compiler, specs []spec.Spec, opts ...vm.EnvironmentOption) {
	runner := spec.Runner{
		Name:     name,
		Compiler: c,
		Env:      opts,
	}

	runner.Run(t, specs)
}

func RunSpecs(t *testing.T, specs []spec.Spec, opts ...vm.EnvironmentOption) {
	RunSpecsWith(t, "VM/O0", compiler.New(compiler.WithOptimizationLevel(compiler.O0)), specs, opts...)
	RunSpecsWith(t, "VM/O1", compiler.New(compiler.WithOptimizationLevel(compiler.O1)), specs, opts...)
}
