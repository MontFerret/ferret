package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestForDoWhileCompiles(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		spec.New(
			`
FOR DO WHILE false
	RETURN 1
`,
			"anonymous loop variable",
		),
		spec.New(
			`
FOR i DO WHILE false
	RETURN i
`,
			"named loop variable",
		),
		spec.New(
			`
FOR _ DO WHILE false
	RETURN 1
`,
			"discard loop variable",
		),
	}, compiler.O0, compiler.O1)
}
