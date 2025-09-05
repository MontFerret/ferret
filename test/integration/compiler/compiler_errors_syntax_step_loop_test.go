package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestStepLoopSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			FOR i = 0 WHILE STEP i = i + 1 RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after 'WHILE'",
				Hint:    "STEP loops require a condition after WHILE, e.g., 'FOR i = 0 WHILE i < 10 STEP i = i + 1'.",
			}, "Missing WHILE condition in STEP loop"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Syntax error: no viable alternative at input 'FOR i = 0 WHILE i < 5 STEP RETURN'",
				Hint:    "Check your syntax. Did you forget to write something?",
			}, "STEP followed by RETURN (parsed as incomplete query)"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP i RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Syntax error: no viable alternative at input 'FOR i = 0 WHILE i < 5 STEP i RETURN'",
				Hint:    "Check your syntax. Did you forget to write something?",
			}, "Missing '=' in STEP assignment"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Syntax error: no viable alternative at input 'FOR i = 0 WHILE i < 5 STEP\\n\\t\\t'",
				Hint:    "Check your syntax. Did you forget to write something?",
			}, "Incomplete STEP clause at end"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Syntax error: no viable alternative at input 'FOR i = 0 WHILE i < 5 STEP x\\n\\t\\t'",
				Hint:    "Check your syntax. Did you forget to write something?",
			}, "Incomplete STEP clause with variable only"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP i = RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing expression after '=' in STEP"),

		// Additional test cases for better coverage
		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP j = j + 1 RETURN i
		`, E{
				Kind:    compiler.SemanticError,
				Message: "step variable missmatch: expected 'i' but got 'j'",
			}, "Variable mismatch between FOR and STEP clauses"),

		ErrorCase(
			`
			FOR i WHILE i < 5 STEP i = i + 1 RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Syntax error: missing '(' at 'i'",
			}, "Missing initial assignment in FOR clause"),

		ErrorCase(
			`
			FOR = 0 WHILE i < 5 STEP i = i + 1 RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected loop variable before 'IN'",
			}, "Missing variable name in FOR clause"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 i = i + 1 RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Syntax error: no viable alternative at input 'FOR i = 0 WHILE i < 5 i'",
				Hint:    "Check your syntax. Did you forget to write something?",
			}, "Missing STEP keyword"),

		ErrorCase(
			`
			FOR i = WHILE i < 5 STEP i = i + 1 RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
			}, "Missing initial value in FOR assignment"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP x = x + 1 RETURN i
		`, E{
				Kind:    compiler.SemanticError,
				Message: "step variable missmatch: expected 'i' but got 'x'",
			}, "Different variable in STEP clause not defined"),
	})
}
