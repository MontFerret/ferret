package exec

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func RunSpecsWith(t *testing.T, name string, c *compiler.Compiler, specs []Spec, baseSpecs []spec.Spec, opts ...vm.EnvironmentOption) {
	runner := spec.Runner{
		Name:     name,
		Compiler: c,
		EnvOpts:  opts,
		SpecEnvOpts: func(i int) []vm.EnvironmentOption {
			return specs[i].EnvOptions
		},
	}

	runner.Run(t, baseSpecs)
}

func RunSpecs(t *testing.T, specs []Spec, opts ...vm.EnvironmentOption) {
	baseSpecs := make([]spec.Spec, len(specs))

	for i, s := range specs {
		baseSpecs[i] = s.Base
	}

	RunSpecsWith(t, "VM O0:", compiler.New(compiler.WithOptimizationLevel(compiler.O0)), specs, baseSpecs, opts...)
	RunSpecsWith(t, "VM O1:", compiler.New(compiler.WithOptimizationLevel(compiler.O1)), specs, baseSpecs, opts...)
}
