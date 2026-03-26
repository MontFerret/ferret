package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestMatchErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`RETURN MATCH { a: 1, b: 2 } ( { a: v, b: v } => v, _ => 0, )`,
			E{
				Kind:    parserd.NameError,
				Message: "duplicate binding 'v' in MATCH pattern",
			},
		),
	})
}
