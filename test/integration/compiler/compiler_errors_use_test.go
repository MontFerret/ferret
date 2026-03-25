package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestUseErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`
			USE IO::FS AS F
			USE IO::NET AS F

			RETURN F::READ("file.txt")
		`,
			E{
				Kind:    parserd.NameError,
				Message: "USE alias 'F' is already defined",
			},
			"Duplicate USE alias"),
	})
}
