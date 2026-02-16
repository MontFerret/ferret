package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestCollectNameErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
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
		ErrorCase(
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
		ErrorCase(
			`
			FOR u IN [1]
				COLLECT g = u INTO grouped KEEP missing
				RETURN grouped
		`, E{
				Kind:    parserd.NameError,
				Message: "Variable 'missing' is not defined",
			}, "COLLECT KEEP should report missing variables as diagnostics"),
		ErrorCase(
			`
			FOR u IN [1, 2, 3]
				COLLECT WITH sum INTO total
				RETURN total
		`, E{
				Kind:    diagnostics.Kind("SemanticError"),
				Message: "Invalid count projection",
				Hint:    "Use WITH COUNT INTO <variable>.",
			}, "COLLECT WITH should reject non-COUNT identifiers without panicking"),
	})
}
