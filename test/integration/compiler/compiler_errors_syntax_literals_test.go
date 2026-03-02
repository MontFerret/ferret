package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestLiteralsSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			LET i = "foo
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string literal (closing quote missing)"),

		ErrorCase(
			`
			LET i = "foo bar
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing)"),

		ErrorCase(
			`
			LET i = foo"
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string literal (opening quote missing)"),

		ErrorCase(
			`
			LET i = foo bar"
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete multi-string literal  (opening quote missing)"),

		ErrorCase(
			`
			LET i = 'foo
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (closing quote missing) 2"),

		ErrorCase(
			`
			LET i = 'foo bar
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing) 2"),

		ErrorCase(
			`
			LET i = foo'
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (opening quote missing) 2"),

		ErrorCase(
			`
			LET i = foo bar'
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete multi-string literal  (opening quote missing) 2"),

		ErrorCase(
			"LET i = `foo "+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete string literal (closing quote missing) 3"),

		ErrorCase(
			"LET i = `foo bar"+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing) 3"),

		ErrorCase(
			"LET i = foo` "+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete string literal (opening quote missing) 3"),

		ErrorCase(
			"LET i = foo bar` "+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete multi-string literal  (opening quote missing) 3"),

		ErrorCase(
			`
			LET i = { "foo: }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string literal (closing quote missing) 4"),

		ErrorCase(
			`
			LET i = { "foo bar: }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing) 4"),

		ErrorCase(
			`
			LET i = { foo": }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string literal (opening quote missing) 4"),

		SkipErrorCase(
			`
			LET i = { foo bar": }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string",
			}, "Incomplete multi-string literal  (opening quote missing) 4"),

		ErrorCase(
			`
			LET i = { 'foo: }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (closing quote missing) 5"),

		ErrorCase(
			`
			LET i = { foo': }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (opening quote missing) 5"),

		SkipErrorCase(
			`
			LET i = { foo bar': }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete multi-string literal  (opening quote missing) 5"),

		ErrorCase(
			"LET i = { 'foo: }"+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (closing quote missing) 6"),

		ErrorCase(
			"LET i = { 'foo bar: }"+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing) 6"),

		SkipErrorCase(
			`
			LET o = { foo: "bar" }
			LET i = o.
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete member access"),

		SkipErrorCase(
			`
			LET o = { foo: "bar" }
			LET i = o.
			FN(i)
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete member access 2"),

		ErrorCase(
			`
			LET i = [
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed array literal",
				Hint:    "Add a closing ']' to complete the array.",
			}, "Incomplete array literal"),

		ErrorCase(
			`
			LET i = [1
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed array literal",
				Hint:    "Add a closing ']' to complete the array.",
			}, "Incomplete array literal 2"),

		ErrorCase(
			`
			LET i = [,]
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected a valid list of values",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete array literal 3"),

		ErrorCase(
			`
			LET i = {
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed object literal",
				Hint:    "Add a closing '}' to complete the object.",
			}, "Incomplete object literal"),

		ErrorCase(
			`
			LET i = { foo: }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected value after object property name",
				Hint:    "Provide a value for the property, e.g. { foo: 123 }.",
			}, "Incomplete object literal 2"),

		ErrorCase(
			`
			LET i = { : }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected property name before ':'",
				Hint:    "Object properties must have a name before the colon, e.g. { property: 123 }.",
			}, "Incomplete object literal 3"),

		SkipErrorCase(
			`
			LET i = { a 123 }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected property name before ':'",
				Hint:    "Object properties must have a name before the colon, e.g. { property: 123 }.",
			}, "Incomplete object literal 4"),

		ErrorCase(
			`
			LET arr = [1, 2, 3]
			LET v = arr[1
			RETURN v
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected a RETURN or FOR clause at end of query",
				Hint:    "All queries must return a value. Add a RETURN statement to complete the query.",
			}, "Unclosed computed property expression"),

		ErrorCase(
			`
			LET arr = [1, 2, 3]
			LET v = arr[]
			RETURN v
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected a RETURN or FOR clause at end of query",
				Hint:    "All queries must return a value. Add a RETURN statement to complete the query.",
			}, "Invalid computed property expression"),
	})
}
