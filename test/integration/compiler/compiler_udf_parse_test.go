package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestUdfNestedLetReturnParses(t *testing.T) {
	src := `
FUNC outer(a) (
  FUNC inner(b) (
    RETURN b
  )
  LET v = inner(1)
  RETURN v
)
RETURN outer(2)
`
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	if _, err := c.Compile(file.NewSource("udf_nested_let_return", src)); err != nil {
		t.Fatalf("compile failed: %v", err)
	}
}
