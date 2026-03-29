package vm_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestTernaryOperator(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN 1 < 2 ? 3 : 4", 3),
		S("RETURN 1 > 2 ? 3 : 4", 4),
		S("RETURN 2 ? : 4", 2),
		S("LET foo = TRUE RETURN foo ? TRUE : FALSE", true),
		S("LET foo = FALSE RETURN foo ? TRUE : FALSE", false),
		Array("FOR i IN [1, 2, 3, 4, 5, 6] RETURN i < 3 ? i * 3 : i * 2", []any{3, 6, 6, 8, 10, 12}),
		Array(`FOR i IN [NONE, 2, 3, 4, 5, 6] RETURN i ? : i`, []any{nil, 2, 3, 4, 5, 6}),
		S(`RETURN 0 && true ? "1" : "some"`, "some", "Should return 'some' when first operand is 0"),
		S(`RETURN length([]) > 0 && true ? "1" : "some"`, "some", "Should return 'some' when first operand is an empty array"),
	})

	t.Run("Should compile ternary operator with default values", func(t *testing.T) {
		vals := []string{
			"0",
			"0.0",
			"''",
			"NONE",
			"FALSE",
		}

		c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))

		for _, val := range vals {
			val := val

			t.Run(val, func(t *testing.T) {
				p, err := c.Compile(source.NewAnonymous(fmt.Sprintf(`
			FOR i IN [%s, 1, 2, 3]
				RETURN i ? i * 2 : 'no value'
		`, val)))

				if err != nil {
					t.Fatalf("compile failed: %v", err)
				}

				out, err := spec.Run(p)

				if err != nil {
					t.Fatalf("run failed: %v", err)
				}

				if string(out) != `["no value",2,4,6]` {
					t.Fatalf("unexpected output: got %s", string(out))
				}
			})
		}
	})
}
