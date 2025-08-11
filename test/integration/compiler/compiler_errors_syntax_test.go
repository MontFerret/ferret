package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			LET
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name"),

		ErrorCase(
			`
			LET
			RETURN 5
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name 2"),

		ErrorCase(
			`
			LET = 1
			RETURN NONE
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "Missing variable name 3"),

		ErrorCase(
			`
			LET i = "foo
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string"),

		ErrorCase(
			`
			LET i = 'foo
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string 2"),

		ErrorCase(
			"LET i = `foo "+
				"RETURN i", E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete string 3"),

		ErrorCase(
			`
			LET i = { "foo: }
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string 4"),

		ErrorCase(
			`
			LET i = { 'foo: }
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string 5"),

		ErrorCase(
			"LET i = { 'foo: }"+
				"RETURN i", E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string 6"),

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
			LET a = 1
			LET b = 2
			LET i = (a ||
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
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
				Kind:    compiler.SyntaxError,
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
				Kind:    compiler.SyntaxError,
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
				Kind:    compiler.SyntaxError,
				Message: "Expected right-hand expression after 'AND'",
				Hint:    "Provide an expression after the logical operator, e.g. (a AND b).",
			}, "Incomplete logical expression 4"),

		SkipErrorCase(
			`
			LET i = NONE
			RETURN i,
		`, E{
				Kind:    compiler.SyntaxError,
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
				Kind:    compiler.SyntaxError,
				Message: "Unclosed parenthesized expression",
				Hint:    "Add a closing ')' to complete the expression.",
			}, "Unclosed grouping 2"),

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
			LET i =
			LET j = 5
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value 2"),

		ErrorCase(
			`
			LET i =
			FOR j IN [1, 2, 3] RETURN j
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value 3"),

		SkipErrorCase(
			`
			LET o = { foo: "bar" }
			LET i = o.
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete member access"),

		SkipErrorCase(
			`
			LET o = { foo: "bar" }
			LET i = o.
			FUNC(i)
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete member access 2"),

		ErrorCase(
			`
			FUNC(1,
			RETURN NONE
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after ','",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete function call"),

		ErrorCase(
			`
			FUNC(,)
			RETURN NONE
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected a valid list of arguments",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete function call 2"),

		ErrorCase(
			`
			FUNC(
			RETURN NONE
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed function call",
				Hint:    "Add a closing ')' to complete the function call.",
			}, "Incomplete function call 3"),

		ErrorCase(
			`
			FUNC(1
			RETURN NONE
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed function call",
				Hint:    "Add a closing ')' to complete the function call.",
			}, "Incomplete function call 4"),

		ErrorCase(
			`
			LET i = [
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed array literal",
				Hint:    "Add a closing ']' to complete the array.",
			}, "Incomplete array literal"),

		ErrorCase(
			`
			LET i = [1
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed array literal",
				Hint:    "Add a closing ']' to complete the array.",
			}, "Incomplete array literal 2"),

		ErrorCase(
			`
			LET i = [,]
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected a valid list of values",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete array literal 3"),

		ErrorCase(
			`
			LET i = {
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Unclosed object literal",
				Hint:    "Add a closing '}' to complete the object.",
			}, "Incomplete object literal"),

		ErrorCase(
			`
			LET i = { foo: }
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected value after object property name",
				Hint:    "Provide a value for the property, e.g. { foo: 123 }.",
			}, "Incomplete object literal 2"),

		ErrorCase(
			`
			LET i = { : }
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected property name before ':'",
				Hint:    "Object properties must have a name before the colon, e.g. { property: 123 }.",
			}, "Incomplete object literal 3"),

		SkipErrorCase(
			`
			LET i = { a 123 }
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected property name before ':'",
				Hint:    "Object properties must have a name before the colon, e.g. { property: 123 }.",
			}, "Incomplete object literal 4"),

		ErrorCase(
			`
			LET r = 0..
			RETURN r
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected end value after '..' in range expression",
				Hint:    "Provide an end value to complete the range, e.g. 0..10.",
			}, "Incomplete range"),

		SkipErrorCase(
			`
			LET r = ..0
			RETURN r
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected end value before '..' in range expression",
				Hint:    "Object properties must have a name before the colon, e.g. { property: 123 }.",
			}, "Incomplete range 2"),
	})
}
