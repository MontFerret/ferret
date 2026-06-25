package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestForInSourceExpressionsCompile(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
			LET items = [1]
			FOR item IN items
				RETURN item
		`, noCompilerError, "identifier source"),
		ProgramCheck(`
			LET response = { items: [1] }
			FOR item IN response.items
				RETURN item
		`, noCompilerError, "member source"),
		ProgramCheck(`
			FOR item IN [1, 2, 3]
				RETURN item
		`, noCompilerError, "array literal source"),
		ProgramCheck(`
			FUNC GET_ITEMS() => [1]
			FOR item IN GET_ITEMS()
				RETURN item
		`, noCompilerError, "function call source"),
		ProgramCheck(`
			FOR order IN QUERY "/products/index.json" IN @api WITH {
				query: {
					status: "open",
					limit: 20
				}
			}
				RETURN order
		`, noCompilerError, "query source"),
		ProgramCheck(`
			FOR order IN (
				QUERY "/products/index.json" IN @api WITH {
					query: {
						status: "open",
						limit: 20
					}
				}
			)
				RETURN order
		`, noCompilerError, "parenthesized query source"),
		ProgramCheck(`
			FOR order IN QUERY "/orders" IN @api
				RETURN order
		`, noCompilerError, "query source with nested IN"),
		ProgramCheck(`
			LET items = [1]
			FOR item IN WAITFOR VALUE items
				RETURN item
		`, noCompilerError, "waitfor source"),
		ProgramCheck(`
			FOR item IN MATCH "a" ("a" => [1], _ => [2])
				RETURN item
		`, noCompilerError, "match source"),
	})
}
