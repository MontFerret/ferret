package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestUdfNestedLetReturnParses(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		spec.NewSpec(
			`
FUNC outer(a) (
  FUNC inner(b) (
    RETURN b
  )
  LET v = inner(1)
  RETURN v
)
RETURN outer(2)
`,
			"nested udf block with let/return",
		),
	})
}
