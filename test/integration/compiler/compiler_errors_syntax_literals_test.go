package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
	specexec "github.com/MontFerret/ferret/v2/test/spec/exec"
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
				Hint:    "Add ']' to close the array.",
			}, "Incomplete array literal"),

		Failure(
			`
			LET i = [1
			RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Unclosed array literal",
				Hint:    "Add ']' to close the array.",
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

func TestIncompleteArrayDiagnosticRendering(t *testing.T) {
	for _, test := range []struct {
		query   string
		snippet string
		span    source.Span
		line    int
		column  int
	}{
		{
			query: "let arr = [1, 2, 3", span: source.Span{Start: 18, End: 18}, line: 1, column: 19,
			snippet: "1 | let arr = [1, 2, 3\n  |                   ^ expected ']'\n",
		},
		{
			query: "let arr = [1, 2,", span: source.Span{Start: 16, End: 16}, line: 1, column: 17,
			snippet: "1 | let arr = [1, 2,\n  |                 ^ expected ']'\n",
		},
		{
			query: "return [\n  \"é🙂\"", span: source.Span{Start: 19, End: 19}, line: 2, column: 11,
			snippet: "2 |   \"é🙂\"\n  |       ^ expected ']'\n",
		},
	} {
		t.Run(test.query, func(t *testing.T) {
			src := source.New(".tmp/errors/arr.fql", test.query)
			_, err := mustNewCompiler(t).Compile(t.Context(), src)
			diagnostic := firstCompilationError(err)
			if diagnostic == nil {
				t.Fatalf("expected diagnostic, got %v", err)
			}

			if diagnostic.Kind != parserd.SyntaxError || diagnostic.Message != "Unclosed array literal" || diagnostic.Hint != "Add ']' to close the array." {
				t.Fatalf("unexpected diagnostic: %+v", diagnostic)
			}

			want := pkgdiagnostics.NewMainErrorSpan(test.span, "expected ']'")
			if len(diagnostic.Spans) != 1 || diagnostic.Spans[0] != want {
				t.Fatalf("primary span = %+v, want %+v", diagnostic.Spans, want)
			}

			formatted := pkgdiagnostics.Format(err)
			location := fmt.Sprintf(" --> .tmp/errors/arr.fql:%d:%d\n", test.line, test.column)
			for _, required := range []string{"SyntaxError: Unclosed array literal\n", location, test.snippet, "Hint: Add ']' to close the array.\n"} {
				if !strings.Contains(formatted, required) {
					t.Errorf("missing %q in:\n%s", required, formatted)
				}
			}

			for _, forbidden := range []string{"Expected expression after ','", "no viable alternative", "Syntax error:"} {
				if strings.Contains(formatted, forbidden) {
					t.Errorf("unexpected %q in:\n%s", forbidden, formatted)
				}
			}

			if err != diagnostic || strings.Count(formatted, "SyntaxError:") != 1 {
				t.Fatalf("expected one actionable diagnostic, got:\n%s", formatted)
			}
		})
	}
}

func TestIncompleteArrayLiteralBoundaries(t *testing.T) {
	for _, test := range []struct {
		name   string
		query  string
		before string
	}{
		{name: "empty", query: "return ["},
		{name: "one entry", query: "return [1"},
		{name: "trailing comma", query: "return [1,"},
		{name: "final newline", query: "return [1,\n", before: "\n"},
		{name: "trailing whitespace", query: "return [1,  \t", before: "  \t"},
		{name: "line comment", query: "return [1, // ] , [", before: " // ] , ["},
		{name: "block comment", query: "return [1 /* ] , [ */", before: " /* ] , [ */"},
		{name: "comment between entries", query: "return [1, /* ], [ */ 2"},
		{name: "multiline CRLF", query: "return [\r\n  1,\r\n  2, // ] ,\r\n", before: " // ] ,\r\n"},
		{name: "Unicode before insertion", query: "let text = 'é🙂'\nreturn [text,"},
		{name: "delimiter strings", query: "return [\"[,]}\", '(',"},
		{name: "complete nested array", query: "return [1, [2, 3]"},
		{name: "two unfinished arrays", query: "return [1, [2, 3"},
		{name: "three unfinished arrays", query: "return [1, [2, [3,"},
		{name: "empty inner array", query: "return [1, ["},
		{name: "complete object", query: "return [1, { value: 2 }"},
		{name: "complete call", query: "func value() => 2\nreturn [1, value()"},
		{name: "keyword-named call", query: "func return() => 2\nreturn [1, return()"},
		{name: "keyword property", query: "let doc = { return: 1 }\nreturn [doc.return"},
		{name: "complete computed key", query: "let key = 'value'\nreturn [{ [key]: 1 }"},
		{name: "complete arithmetic", query: "return [1, 2 + 3"},
		{name: "complete spread", query: "return [1, ...[2, 3]"},
		{name: "complete template", query: "return [1, `[, ] ${2}`"},
		{name: "long array", query: "return [" + strings.Repeat("1, ", 256) + "2"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := mustNewCompiler(t).Compile(t.Context(), source.New("array.fql", test.query))
			diagnostic := firstCompilationError(err)
			if diagnostic == nil || diagnostic.Kind != parserd.SyntaxError || diagnostic.Message != "Unclosed array literal" || diagnostic.Hint != "Add ']' to close the array." {
				t.Fatalf("unexpected diagnostic: %v", err)
			}

			offset := len(test.query) - len(test.before)
			want := pkgdiagnostics.NewMainErrorSpan(source.Span{Start: offset, End: offset}, "expected ']'")
			if len(diagnostic.Spans) != 1 || diagnostic.Spans[0] != want {
				t.Fatalf("spans = %+v, want %+v", diagnostic.Spans, want)
			}

			if err != diagnostic {
				t.Fatalf("unexpected diagnostic cascade: %v", err)
			}
		})
	}
}

func TestIncompleteArrayPreservesOtherSyntaxErrors(t *testing.T) {
	for _, test := range []struct {
		query   string
		message string
	}{
		{query: "return [1,,2]", message: "Expected expression after ','"},
		{query: "return [1,,2", message: "Expected expression after ','"},
		{query: "return [1 2]", message: "Expected ',' between array items"},
		{query: "return [1 2", message: "Expected ',' between array items"},
		{query: "return [1 2, [3", message: "Expected ',' between array items"},
		{query: "return [1,, [2"},
		{query: "return [1, return 2]"},
		{query: "return [1, 2 +", message: "Expected right-hand expression after '+'"},
		{query: "return [1, ...", message: "Expected expression after '...' in array spread"},
		{query: "return [1, 2)", message: "Unexpected closing delimiter ')'"},
		{query: "return [1, (2 + 3", message: "Unclosed parenthesized expression"},
		{query: "return [1, { value: 2"},
		{query: "func value(x) => x\nreturn [1, value(2"},
		{query: "return [1, 'unfinished"},
		{query: "return [1, `unfinished"},
		{query: "return [1, 2 ? 3"},
		{query: "let arr = [1, 2]\nreturn arr[1"},
		{query: "return [1][0"},
		{query: "let arr = [1, 2]\nreturn [arr[1"},
		{query: "let key = 'value'\nreturn { [key"},
		{query: "let key = 'value'\nreturn { first: 1, [key"},
		{query: "let [first, second"},
		{query: "let [[first"},
		{query: "for [first, second"},
	} {
		t.Run(test.query, func(t *testing.T) {
			_, err := mustNewCompiler(t).Compile(t.Context(), source.NewAnonymous(test.query))
			diagnostic := firstCompilationError(err)
			if diagnostic == nil || diagnostic.Kind != parserd.SyntaxError {
				t.Fatalf("expected syntax diagnostic, got %v", err)
			}

			if diagnostic.Message == "Unclosed array literal" {
				t.Fatalf("more-specific error hidden: %v", err)
			}

			if test.message != "" && diagnostic.Message != test.message {
				t.Fatalf("message = %q, want %q", diagnostic.Message, test.message)
			}

			for _, annotation := range diagnostic.Spans {
				span := annotation.Span
				if span.Start < 0 || span.End < span.Start || span.End > len(test.query) {
					t.Fatalf("out-of-bounds span: %+v", span)
				}
			}
		})
	}
}

func TestArrayOptionalContentsRemainValid(t *testing.T) {
	specexec.RunSpecs(t, []spec.Spec{
		specexec.S("return []", []any{}),
		specexec.S("return [1, 2,]", []any{float64(1), float64(2)}),
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

	_, err := mustNewCompiler(t).Compile(t.Context(), source.NewAnonymous(query))
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

	position := diag.Source.PositionAt(diag.Spans[0].Span)
	if position.Line != 7 || position.Column != 47 {
		t.Fatalf("unexpected span location: got %d:%d, want 7:47", position.Line, position.Column)
	}

	if err != diag {
		t.Fatalf("expected one native diagnostic, got %T: %v", err, err)
	}

	offset := strings.Index(query, "129.99 }") + len("129.99 }")
	if got := diag.Spans[0].Span; got != (source.Span{Start: offset, End: offset}) {
		t.Fatalf("span = %+v, want insertion after previous array item at %d", got, offset)
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
