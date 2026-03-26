package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"

	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestCollectNameErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`
			LET users = [1, 2, 3]
			FOR u IN users
				LET total = 0
				COLLECT AGGREGATE total = COUNT(u)
				RETURN total
		`, E{
				Kind:    parserd.NameError,
				Message: "Variable 'total' is already defined",
			}, "COLLECT AGGREGATE should reject duplicate output variable names"),
		Failure(
			`
			LET users = [1, 2, 3]
			FOR u IN users
				LET g = "already-defined"
				COLLECT g = u
				RETURN g
		`, E{
				Kind:    parserd.NameError,
				Message: "Variable 'g' is already defined",
			}, "COLLECT grouping should reject duplicate group variable names"),
		Failure(
			`
			FOR u IN [1]
				COLLECT g = u INTO grouped KEEP missing
				RETURN grouped
		`, E{
				Kind:    parserd.NameError,
				Message: "Variable 'missing' is not defined",
			}, "COLLECT KEEP should report missing variables as diagnostics"),
		Failure(
			`
				FOR u IN [1, 2, 3]
					COLLECT WITH sum INTO total
					RETURN total
			`, E{
				Kind: parserd.SyntaxError,
			}, "COLLECT WITH should reject non-COUNT identifiers without panicking"),
	})
}
