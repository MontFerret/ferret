package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
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
				Kind:    compiler.NameError,
				Message: "USE alias 'F' is already defined",
			},
			"Duplicate USE alias"),
	})
}
