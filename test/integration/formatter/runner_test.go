package formatter_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/test/spec/format"
)

func RunSpecs(t *testing.T, specs []format.Spec) {
	RunSpecsWith(t, formatter.New(), specs)
}

func RunSpecsWith(t *testing.T, f *formatter.Formatter, specs []format.Spec) {
	runner := &format.Runner{
		Name:      "Formatter",
		Formatter: f,
	}

	runner.Run(t, specs)
}
