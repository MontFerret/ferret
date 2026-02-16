package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestUseErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
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
