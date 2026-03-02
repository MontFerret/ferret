package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			LET
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name"),

		ErrorCase(
			`
			LET
			RETURN 5
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name 2"),

		ErrorCase(
			`
			LET = 1
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name 3"),

		ErrorCase(
			`
			LET i = NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected a RETURN or FOR clause at end of query",
				Hint:    "All queries must return a value. Add a RETURN statement to complete the query.",
			}, "Missing return statement"),

		ErrorCase(
			`
			LET i = NONE
			RETURN
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'RETURN'",
				Hint:    "Did you forget to provide a value to return?",
			}, "Missing return value"),
		ErrorCase(
			`
			FUNC f(x)
			  RETURN x
			RETURN f(1)
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected '=>' or '(' after function declaration",
				Hint:    "Use 'FUNC f(x) => expr' or 'FUNC f(x) ( ... RETURN expr )'.",
			}, "Undelimited function body"),
		ErrorCase(
			`
			FUNC f() => RETURN 1
			RETURN f()
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=>'",
				Hint:    "Provide an expression, e.g. FUNC f() => x + 1",
			}, "Missing arrow expression"),

		ErrorCase(
			`
			LET a = 1
			LET b = 2
			LET i = (a ||
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected right-hand expression after '||'",
				Hint:    "Provide an expression after the logical operator, e.g. (a || b).",
			}, "Incomplete logical expression"),

		ErrorCase(
			`
			LET a = 1
			LET b = 2
			LET i = (a OR
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected right-hand expression after 'OR'",
				Hint:    "Provide an expression after the logical operator, e.g. (a OR b).",
			}, "Incomplete logical expression 2"),

		ErrorCase(
			`
			LET a = 1
			LET b = 2
			LET i = (a &&
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected right-hand expression after '&&'",
				Hint:    "Provide an expression after the logical operator, e.g. (a && b).",
			}, "Incomplete logical expression 3"),

		ErrorCase(
			`
			LET a = 1
			LET b = 2
			LET i = (a AND
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected right-hand expression after 'AND'",
				Hint:    "Provide an expression after the logical operator, e.g. (a AND b).",
			}, "Incomplete logical expression 4"),

		ErrorCase(
			`
			LET a = 1
			LET b = 2
			LET i = b > 1 ? a :
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after ':' in ternary operator",
				Hint:    "Provide an expression after the colon to complete the ternary operation.",
			}, "Incomplete ternary expression"),

		ErrorCase(
			`
			LET a = 1
			LET b = 2
			LET i = b > 1 ? 1 + 1 + 1 :
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after ':' in ternary operator",
				Hint:    "Provide an expression after the colon to complete the ternary operation.",
			}, "Incomplete ternary expression 2"),

		ErrorCase(
			`
			LET a = 1
			LET b = 2
			LET i = b > 1 ?
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '?' in ternary operator",
				Hint:    "Provide an expression after the question mark to complete the ternary operation.",
			}, "Incomplete ternary expression 3"),

		SkipErrorCase(
			`
			LET i = NONE
			RETURN i,
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "--",
				Hint:    "--",
			}, "Dangling comma in return"),

		ErrorCase(
			`
			LET a = 1
			LET b = 2
			LET i = (a AND b
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed parenthesized expression",
				Hint:    "Add a closing ')' to complete the expression.",
			}, "Unclosed grouping 2"),

		ErrorCase(
			`
			LET i =
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value"),

		ErrorCase(
			`
			LET i =
			LET j = 5
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value 2"),

		ErrorCase(
			`
			LET i =
			FOR j IN [1, 2, 3] RETURN j
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value 3"),

		ErrorCase(
			`
			FN(1,
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after ','",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete function call"),

		ErrorCase(
			`
			FN(,)
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected a valid list of arguments",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete function call 2"),

		ErrorCase(
			`
			FN(
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed function call",
				Hint:    "Add a closing ')' to complete the function call.",
			}, "Incomplete function call 3"),

		ErrorCase(
			`
			FN(1
			RETURN NONE
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed function call",
				Hint:    "Add a closing ')' to complete the function call.",
			}, "Incomplete function call 4"),

		ErrorCase(
			`
			LET r = 0..
			RETURN r
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected end value after '..' in range expression",
				Hint:    "Provide an end value to complete the range, e.g. 0..10.",
			}, "Incomplete range"),

		SkipErrorCase(
			`
			LET r = ..0
			RETURN r
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected end value before '..' in range expression",
				Hint:    "Object properties must have a name before the colon, e.g. { property: 123 }.",
			}, "Incomplete range 2"),
	})
}
