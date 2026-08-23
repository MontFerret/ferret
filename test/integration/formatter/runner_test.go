package formatter_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/test/spec/format"
)

func RunSpecs(t *testing.T, specs []format.Spec) {
	RunSpecsWith(t, mustNewFormatter(t), specs)
}

func RunSpecsWith(t *testing.T, f *formatter.Formatter, specs []format.Spec) {
	runner := &format.Runner{
		Name:      "Formatter",
		Formatter: f,
	}

	runner.Run(t, specs)
}

func mustNewFormatter(t testing.TB, setters ...formatter.Option) *formatter.Formatter {
	t.Helper()

	formatterInstance, err := formatter.New(setters...)
	if err != nil {
		t.Fatalf("formatter.New() error = %v", err)
	}

	return formatterInstance
}
