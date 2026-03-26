package exec

import (
	"fmt"
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
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		RunSpecsWith(t, fmt.Sprintf("VM/O%d", level), compiler.New(compiler.WithOptimizationLevel(level)), specs, opts...)
	}
}

func RunSpecFactory(t *testing.T, factory func() []spec.Spec, opts ...vm.EnvironmentOption) {
	t.Helper()

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, level := range levels {
		RunSpecsWith(t, fmt.Sprintf("VM/O%d", level), compiler.New(compiler.WithOptimizationLevel(level)), factory(), opts...)
	}
}

func RunProgramSpecs(t *testing.T, specs []spec.Spec, opts ...vm.EnvironmentOption) {
	t.Helper()

	RunSpecsWith(t, "VM/Program", compiler.New(compiler.WithOptimizationLevel(compiler.O0)), specs, opts...)
}
