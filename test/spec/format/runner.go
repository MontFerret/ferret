package format

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/test/spec"
)

type Runner struct {
	Name      string
	Formatter *formatter.Formatter
}

func (r *Runner) Run(t *testing.T, specs []Spec) {
	t.Helper()

	for _, s := range specs {
		suiteName := s.SuiteName(r.Name)

		t.Run(suiteName, func(t *testing.T) {
			if s.Base.SkipInfo.Active {
				t.Skip(s.Base.SkipInfo.Reason)
			}

			var out strings.Builder

			err := r.Formatter.Format(&out, file.NewSource("Test case", s.Base.Expression))

			if err != nil {
				if s.Base.DebugOutput {
					spec.PrintError(t, err)
				}

				if s.Output.Error.Defined() {
					s.Output.Error.Assert(t, err)
					return
				}

				t.Fatalf("unexpected formatting error:\n%s", diagnostics.Format(err))
			}

			if s.Base.DebugOutput {
				PrintDebug(t, suiteName, s.Base.Expression, out)
			}

			if s.Output.Result.Defined() {
				s.Output.Result.Assert(t, out.String())
			}
		})
	}
}
