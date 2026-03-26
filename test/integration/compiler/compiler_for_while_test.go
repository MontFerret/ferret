package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestForWhile(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		spec.NewSpec(
			`
FOR WHILE UNTIL(5)
	RETURN 1
`,
			"anonymous loop variable",
		),
		spec.NewSpec(
			`
FOR i WHILE UNTIL(5)
	RETURN i
`,
			"named loop variable",
		),
		spec.NewSpec(
			`
FOR _ WHILE UNTIL(2)
	RETURN 1
`,
			"discard loop variable",
		),
	}, compiler.O0, compiler.O1)
}
