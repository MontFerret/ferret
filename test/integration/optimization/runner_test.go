package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func RunUseCases(t *testing.T, level compiler.OptimizationLevel, useCases []spec.Spec) {
	spec.
		NewRunner("optimization", compiler.WithOptimizationLevel(level)).
		Run(t, useCases)
}
