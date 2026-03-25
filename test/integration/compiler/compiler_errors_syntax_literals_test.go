package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestLiteralsSyntaxErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`
			LET i = "foo
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string literal (closing quote missing)"),

		Failure(
			`
			LET i = "foo bar
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing)"),

		Failure(
			`
			LET i = foo"
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string literal (opening quote missing)"),

		Failure(
			`
			LET i = foo bar"
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete multi-string literal  (opening quote missing)"),

		Failure(
			`
			LET i = 'foo
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (closing quote missing) 2"),

		Failure(
			`
			LET i = 'foo bar
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing) 2"),

		Failure(
			`
			LET i = foo'
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (opening quote missing) 2"),

		Failure(
			`
			LET i = foo bar'
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete multi-string literal  (opening quote missing) 2"),

		Failure(
			"LET i = `foo "+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete string literal (closing quote missing) 3"),

		Failure(
			"LET i = `foo bar"+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing) 3"),

		Failure(
			"LET i = foo` "+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete string literal (opening quote missing) 3"),

		Failure(
			"LET i = foo bar` "+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '`' to close the string.",
			}, "Incomplete multi-string literal  (opening quote missing) 3"),

		Failure(
			`
			LET i = { "foo: }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string literal (closing quote missing) 4"),

		Failure(
			`
			LET i = { "foo bar: }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing) 4"),

		Failure(
			`
			LET i = { foo": }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string.",
			}, "Incomplete string literal (opening quote missing) 4"),

		Failure(
			`
			LET i = { foo bar": }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching '\"' to close the string",
			}, "Incomplete multi-string literal  (opening quote missing) 4").Skip(),

		Failure(
			`
			LET i = { 'foo: }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (closing quote missing) 5"),

		Failure(
			`
			LET i = { foo': }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (opening quote missing) 5"),

		Failure(
			`
			LET i = { foo bar': }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete multi-string literal  (opening quote missing) 5").Skip(),

		Failure(
			"LET i = { 'foo: }"+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete string literal (closing quote missing) 6"),

		Failure(
			"LET i = { 'foo bar: }"+
				"RETURN i", E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed string literal",
				Hint:    "Add a matching \"'\" to close the string.",
			}, "Incomplete multi-string literal  (closing quote missing) 6"),

		Failure(
			`
			LET o = { foo: "bar" }
			LET i = o.
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete member access").Skip(),

		Failure(
			`
			LET o = { foo: "bar" }
			LET i = o.
			FN(i)
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete member access 2").Skip(),

		Failure(
			`
			LET i = [
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed array literal",
				Hint:    "Add a closing ']' to complete the array.",
			}, "Incomplete array literal"),

		Failure(
			`
			LET i = [1
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed array literal",
				Hint:    "Add a closing ']' to complete the array.",
			}, "Incomplete array literal 2"),

		Failure(
			`
			LET i = [,]
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected a valid list of values",
				Hint:    "Did you forget to provide a value?",
			}, "Incomplete array literal 3"),

		Failure(
			`
			LET i = {
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed object literal",
				Hint:    "Add a closing '}' to complete the object.",
			}, "Incomplete object literal"),

		Failure(
			`
			LET i = { foo: }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected value after object property name",
				Hint:    "Provide a value for the property, e.g. { foo: 123 }.",
			}, "Incomplete object literal 2"),

		Failure(
			`
			LET i = { : }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected property name before ':'",
				Hint:    "Object properties must have a name before the colon, e.g. { property: 123 }.",
			}, "Incomplete object literal 3"),

		Failure(
			`
			LET i = { a 123 }
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected property name before ':'",
				Hint:    "Object properties must have a name before the colon, e.g. { property: 123 }.",
			}, "Incomplete object literal 4").Skip(),

		Failure(
			`
			LET arr = [1, 2, 3]
			LET v = arr[1
			RETURN v
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed computed property expression",
				Hint:    "Add a closing ']' to complete the computed property expression.",
			}, "Unclosed computed property expression"),

		Failure(
			`
			LET arr = [1, 2, 3]
			LET v = arr[]
			RETURN v
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed computed property expression",
				Hint:    "Add a closing ']' to complete the computed property expression.",
			}, "Invalid computed property expression"),
	})
}
