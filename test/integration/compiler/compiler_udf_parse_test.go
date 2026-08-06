package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestUdfNestedLetReturnParses(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		spec.NewSpec(
			`
FUNC outer(a) {
  FUNC inner(b) {
    RETURN b
  }
  LET v = inner(1)
  RETURN v
}
RETURN outer(2)
`,
			"nested udf block with let/return",
		),
	})
}

func TestUdfMemberStatementsParse(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		spec.NewSpec(
			`
FUNC PARSE_PRICE(product) {
  LET priceNode = QUERY ONE ".product-price" IN product USING css
  LET priceText = priceNode.attributes["data-price"]
  LET price = TO_FLOAT(SUBSTITUTE(priceText, "$", ""))
  RETURN price
}
RETURN PARSE_PRICE({})
`,
			"reported price parser accepts an unparenthesized mixed member initializer",
		),
		spec.NewSpec(
			`
FUNC read(value) {
  LET dot = value.foo
  VAR computed = value["fallback"]
  computed = value.nested["answer"]
  RETURN [dot, computed]
}
RETURN read({ foo: 1, fallback: 2, nested: { answer: 3 } })
`,
			"UDF LET, VAR, and reassignment accept member expressions",
		),
	}, compiler.O0, compiler.O1)
}
