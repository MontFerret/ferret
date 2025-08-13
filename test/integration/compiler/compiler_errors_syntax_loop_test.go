package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestForLoopSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
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
				FILTER
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Incomplete FILTER clause",
				Hint:    "FILTER requires a boolean expression.",
			}, "FILTER with no expression"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				FILTER x =
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Incomplete FILTER clause",
				Hint:    "FILTER requires a boolean expression.",
			}, "FILTER with no expression 2"),

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
				Message: "Dangling comma in LIMIT clause",
				Hint:    "LIMIT accepts one or two arguments. Did you forget to add a value?",
			}, "LIMIT unexpected comma 2"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				LIMIT ,
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Dangling comma in LIMIT clause",
				Hint:    "LIMIT accepts one or two arguments. Did you forget to add a value?",
			}, "LIMIT unexpected comma 3"),

		ErrorCase(
			`
			LET users = []
			FOR x IN users
				LIMIT,
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Dangling comma in LIMIT clause",
				Hint:    "LIMIT accepts one or two arguments. Did you forget to add a value?",
			}, "LIMIT unexpected comma 4"),

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
				COLLECT
				RETURN x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Incomplete COLLECT clause",
				Hint:    "COLLECT must specify a grouping key, an AGGREGATE clause, or WITH COUNT.",
			}, "COLLECT with no variables"),

		ErrorCase(
			`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender INTO
				RETURN {
					gender,
					values
				}`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected variable name after INTO",
				Hint:    "Provide a variable name to store grouped values, e.g. INTO groups.",
			}, "COLLECT INTO with no variable"),

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
				COLLECT AGGREGATE 
				RETURN total
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected variable assignment after AGGREGATE",
				Hint:    "Provide at least one variable assignment, e.g. AGGREGATE total = COUNT(x).",
			}, "COLLECT AGGREGATE without expression 2"),
	})
}
