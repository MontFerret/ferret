package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			LET i = NONE
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected a RETURN or FOR clause at end of query",
				Hint:    "All queries must return a value. Add a RETURN statement to complete the query.",
			}, "Missing return statement"),

		ErrorCase(
			`
			LET i = NONE
			RETURN
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after 'RETURN'",
				Hint:    "Did you forget to provide a value to return?",
			}, "Missing return value"),

		ErrorCase(
			`
			FOR i IN [1, 2, 3]
				RETURN
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after 'RETURN'",
				Hint:    "Did you forget to provide a value to return?",
			}, "Missing return value in for loop"),

		ErrorCase(
			`
			LET i =
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value"),

		ErrorCase(
			`
			FOR i IN 
				RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after 'IN'",
				Hint:    "Each FOR loop must iterate over a collection or range.",
			}, "Missing iterable in FOR"),

		ErrorCase(
			`
			FOR i [1, 2, 3]
				RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected 'IN' after loop variable",
				Hint:    "Use 'FOR x IN [iterable]' syntax.",
			}, "Missing IN in FOR"),

		ErrorCase(
			`
			FOR IN [1, 2, 3]
				RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected loop variable before 'IN'",
				Hint:    "FOR must declare a variable.",
			}, "FOR without variable"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				COLLECT =
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected variable before '=' in COLLECT",
				Hint:    "COLLECT must group by a variable.",
			}, "COLLECT with no variable"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				COLLECT AGGREGATE total = 
				RETURN total
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'total'",
				Hint:    "Did you forget to provide a value?",
			}, "COLLECT AGGREGATE without expression"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				FILTER
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected condition after 'FILTER'",
				Hint:    "FILTER requires a boolean expression.",
			}, "FILTER with no expression"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				LIMIT
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected number after 'LIMIT'",
				Hint:    "LIMIT requires a numeric value.",
			}, "LIMIT with no value"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				LIMIT 1, 2, 3
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Too many arguments provided to LIMIT clause",
				Hint:    "LIMIT accepts at most two arguments: offset and count.",
			}, "LIMIT with too many values"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				LIMIT 1, 2,
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Too many arguments provided to LIMIT clause",
				Hint:    "LIMIT accepts at most two arguments: offset and count.",
			}, "LIMIT unexpected comma"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				LIMIT 1,
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "---",
				Hint:    "FILTER requires a boolean expression.",
			}, "LIMIT unexpected comma 2"),
	})
}
