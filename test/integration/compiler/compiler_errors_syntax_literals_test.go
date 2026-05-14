package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
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
			`RETURN [1 2]`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected ',' between array items",
				Hint:    "Separate array items with commas, e.g. [1, 2, 3].",
			},
			"Missing comma between array items",
		),

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

func TestArrayMissingCommaDiagnosticSpanDoesNotCascade(t *testing.T) {
	query := `LET products = [
    { name: "Widget", price: 19.99 },
    { name: "Gadget", price: 149.99 },
    { name: "Thingamajig", price: 49.99 },
    { name: "Doodad", price: 9.99 },
    { name: "Doohickey", price: 199.99 },
    { name: "Whatchamacallit", price: 129.99 }
    { name: "Contraption", price: 89.99 },
    { name: "Gizmo", price: 179.99 }
]

RETURN products[*
    FILTER TO_FLOAT(.price) > 100
    LIMIT 3
    RETURN {
        title: .name,
        price: TO_FLOAT(.price)
    }
]`

	_, err := compiler.New().Compile(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diag := firstCompilationError(err)
	if diag == nil {
		t.Fatalf("expected diagnostic, got %T", err)
	}

	if diag.Kind != parserd.SyntaxError {
		t.Fatalf("unexpected diagnostic kind: %s", diag.Kind)
	}

	if diag.Message != "Expected ',' between array items" {
		t.Fatalf("unexpected diagnostic message: %q", diag.Message)
	}

	if diag.Hint != "Separate array items with commas, e.g. [1, 2, 3]." {
		t.Fatalf("unexpected diagnostic hint: %q", diag.Hint)
	}

	if len(diag.Spans) == 0 {
		t.Fatal("expected diagnostic span")
	}

	if diag.Spans[0].Label != "missing comma" {
		t.Fatalf("unexpected span label: %q", diag.Spans[0].Label)
	}

	line, col := diag.Source.LocationAt(diag.Spans[0].Span)
	if line != 7 || col != 47 {
		t.Fatalf("unexpected span location: got %d:%d, want 7:47", line, col)
	}

	formatted := pkgdiagnostics.Format(err)
	if got := strings.Count(formatted, "SyntaxError:"); got != 1 {
		t.Fatalf("expected one syntax diagnostic, got %d:\n%s", got, formatted)
	}

	if !strings.Contains(formatted, "7 |     { name: \"Whatchamacallit\", price: 129.99 }\n  |                                               ^ missing comma\n8 |     { name: \"Contraption\", price: 89.99 },") {
		t.Fatalf("diagnostic should point after previous array item, got:\n%s", formatted)
	}

	for _, unexpected := range []string{
		"no viable alternative at input",
		"mismatched input ']'",
	} {
		if strings.Contains(formatted, unexpected) {
			t.Fatalf("formatted diagnostic contains cascade %q:\n%s", unexpected, formatted)
		}
	}
}
