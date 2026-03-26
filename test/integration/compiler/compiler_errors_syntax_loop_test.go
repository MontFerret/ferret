package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestForLoopSyntaxErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`
			FOR i IN [1, 2, 3]
				RETURN
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'RETURN'",
				Hint:    "Did you forget to provide a value to return?",
			}, "Missing return value in for loop"),

		Failure(
			`
			FOR i IN 
				RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'IN'",
				Hint:    "Each FOR loop must iterate over a collection or range.",
			}, "Missing iterable in FOR"),

		Failure(
			`
			FOR i [1, 2, 3]
				RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected 'IN' after loop variable",
				Hint:    "Use 'FOR x IN [iterable]' syntax.",
			}, "Missing IN in FOR"),

		Failure(
			`
			FOR i = 0 WHILE i < 5 STEP i = i + 1
				RETURN i
		`, E{
				Kind: parserd.SyntaxError,
			}, "Legacy STEP syntax fails generically"),

		Failure(
			`
			VAR i = 0
			WHILE i < 10
				i = i + 1
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Standalone WHILE loops are not supported",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "Standalone WHILE loop at top level"),

		Failure(
			`
			VAR i = 0
			DO WHILE i < 10
				i = i + 1
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Standalone DO WHILE loops are not supported",
				Hint:    "Use 'FOR DO WHILE [condition]' or 'FOR x DO WHILE [condition]' syntax.",
			}, "Standalone DO WHILE loop at top level"),

		Failure(
			`
			FOR WHILE
				RETURN 1
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected condition after 'WHILE'",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "FOR WHILE without condition"),

		Failure(
			`
			FOR DO WHILE
				RETURN 1
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected condition after 'WHILE'",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "FOR DO WHILE without condition"),

		Failure(
			`
			FOR 123 WHILE TRUE
				RETURN 1
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected identifier before 'WHILE'",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "FOR WHILE with invalid binding"),

		Failure(
			`
			FOR 123 DO WHILE TRUE
				RETURN 1
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected identifier before 'WHILE'",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "FOR DO WHILE with invalid binding"),

		Failure(
			`
			LET ok = (
				FOR WHILE FALSE
					RETURN 1
			)

			FOR 123 WHILE TRUE
				RETURN ok
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected identifier before 'WHILE'",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "Later FOR WHILE invalid binding is detected"),

		Failure(
			`
			LET ok = (
				FOR DO WHILE FALSE
					RETURN 1
			)

			FOR 123 DO WHILE TRUE
				RETURN ok
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected identifier before 'WHILE'",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "Later FOR DO WHILE invalid binding is detected"),

		Failure(
			`
			LET ok = (
				FOR DO WHILE FALSE
					RETURN 1
			)

			FOR 123 WHILE TRUE
				RETURN ok
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected identifier before 'WHILE'",
				Hint:    "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax.",
			}, "Mixed while-loop forms still find later invalid binding"),

		Failure(
			`
			FOR IN [1, 2, 3]
				RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected loop variable before 'IN'",
				Hint:    "FOR must declare a variable.",
			}, "FOR without variable"),

		Failure(
			`
			LET users = []
			FOR x IN users
				FILTER
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Incomplete FILTER clause",
				Hint:    "FILTER requires a boolean expression.",
			}, "FILTER with no expression"),

		Failure(
			`
			LET users = []
			FOR x IN users
				FILTER x =
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Incomplete FILTER clause",
				Hint:    "FILTER requires a boolean expression.",
			}, "FILTER with no expression 2"),

		Failure(
			`
			LET users = []
			FOR x IN users
				LIMIT
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected number after 'LIMIT'",
				Hint:    "LIMIT requires a numeric value.",
			}, "LIMIT with no value"),

		Failure(
			`
			LET users = []
			FOR x IN users
				LIMIT 1, 2, 3
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Too many arguments provided to LIMIT clause",
				Hint:    "LIMIT accepts at most two arguments: offset and count.",
			}, "LIMIT with too many values"),

		Failure(
			`
			LET users = []
			FOR x IN users
				LIMIT 1, 2,
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Too many arguments provided to LIMIT clause",
				Hint:    "LIMIT accepts at most two arguments: offset and count.",
			}, "LIMIT unexpected comma"),

		Failure(
			`
			LET users = []
			FOR x IN users
				LIMIT 1,
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Dangling comma in LIMIT clause",
				Hint:    "LIMIT accepts one or two arguments. Did you forget to add a value?",
			}, "LIMIT unexpected comma 2"),

		Failure(
			`
			LET users = []
			FOR x IN users
				LIMIT ,
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Dangling comma in LIMIT clause",
				Hint:    "LIMIT accepts one or two arguments. Did you forget to add a value?",
			}, "LIMIT unexpected comma 3"),

		Failure(
			`
			LET users = []
			FOR x IN users
				LIMIT,
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Dangling comma in LIMIT clause",
				Hint:    "LIMIT accepts one or two arguments. Did you forget to add a value?",
			}, "LIMIT unexpected comma 4"),

		Failure(
			`
			LET users = []
			FOR x IN users
				COLLECT =
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable before '=' in COLLECT",
				Hint:    "COLLECT must group by a variable.",
			}, "COLLECT with no variable"),

		Failure(
			`
			LET users = []
			FOR x IN users
				COLLECT
				RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Incomplete COLLECT clause",
				Hint:    "COLLECT must specify a grouping key, an AGGREGATE clause, or WITH COUNT.",
			}, "COLLECT with no variables"),

		Failure(
			`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender INTO
				RETURN {
					gender,
					values
				}`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name after INTO",
				Hint:    "Provide a variable name to store grouped values, e.g. INTO groups.",
			}, "COLLECT INTO with no variable"),

		Failure(
			`
			LET users = []
			FOR x IN users
				COLLECT AGGREGATE total = 
				RETURN total
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'total'",
				Hint:    "Did you forget to provide a value?",
			}, "COLLECT AGGREGATE without expression"),

		Failure(
			`
			LET users = []
			FOR x IN users
				COLLECT AGGREGATE 
				RETURN total
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable assignment after AGGREGATE",
				Hint:    "Provide at least one variable assignment, e.g. AGGREGATE total = COUNT(x).",
			}, "COLLECT AGGREGATE without expression 2"),
	})
}
