package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestFormatterSimplifiesOnlyRedundantParentheses(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "nested statement grouping", input: "((PRINT(1)))\nRETURN NONE", want: "PRINT(1)\nreturn none"},
		{name: "nested value grouping", input: "RETURN ((value))", want: "return value"},
		{name: "higher precedence operand", input: "RETURN 1 + (2 * 3)", want: "return 1 + 2 * 3"},
		{name: "left associative left operand", input: "RETURN (1 - 2) - 3", want: "return 1 - 2 - 3"},
		{name: "right associative right operand", input: "RETURN 1 ?? (2 ?? 3)", want: "return 1 ?? 2 ?? 3"},
		{name: "grammar unary precedence", input: "RETURN -(1 + 2)", want: "return -1 + 2"},
		{name: "simple member source", input: "RETURN (value).member", want: "return value.member"},
		{name: "nested simple member source", input: "RETURN ((value)).member", want: "return value.member"},
		{name: "multiplicative grouping", input: "RETURN (1 + 2) * 3", want: "return (1 + 2) * 3"},
		{name: "left associative right operand", input: "RETURN 1 - (2 - 3)", want: "return 1 - (2 - 3)"},
		{name: "right associative left operand", input: "RETURN (1 ?? 2) ?? 3", want: "return (1 ?? 2) ?? 3"},
		{name: "minus token boundary", input: "RETURN -(-1)", want: "return -(-1)"},
		{name: "plus token boundary", input: "RETURN +(+1)", want: "return +(+1)"},
		{name: "ternary condition associativity", input: "RETURN (a ? b : c) ? d : e", want: "return a ? b : c ? d : e"},
		{name: "ternary true branch", input: "RETURN a ? (b ? c : d) : e", want: "return a ? b ? c : d : e"},
		{name: "ternary false branch", input: "RETURN a ? b : (c ? d : e)", want: "return a ? b : (c ? d : e)"},
		{name: "compound member source", input: "RETURN (1 + 2).member", want: "return (1 + 2).member"},
		{
			name:  "query payload boundary",
			input: "RETURN QUERY (prefix + selector) IN page",
			want:  "return query (prefix + selector) in page",
		},
		{
			name:  "simple query payload",
			input: "RETURN QUERY ((selector)) IN page",
			want:  "return query selector in page",
		},
		{
			name:  "simple dispatch target",
			input: "DISPATCH \"ready\" IN ((target))\nRETURN NONE",
			want:  "dispatch \"ready\" in target\nreturn none",
		},
		{
			name:  "recovery tail ownership",
			input: "RETURN (FAIL() + 1) ON ERROR RETURN NONE",
			want:  "return (FAIL() + 1) on error return none",
		},
		{
			name:  "retry delay grammar boundary",
			input: "RETURN FAIL() ON ERROR RETRY 1 DELAY (0MS OR 1MS) OR RETURN NONE",
			want:  "return FAIL() on error retry 1 delay (0MS or 1MS) or return none",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := formatParenthesesStable(t, test.input)
			if got != test.want {
				t.Fatalf("formatted output:\n%s\nwant:\n%s", got, test.want)
			}
		})
	}
}

func TestFormatterSimplifiesGroupedForOnlyAtStatementBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
		excludes string
	}{
		{
			name:     "standalone grouped loop",
			input:    "(FOR value IN [1] RETURN value)\nRETURN NONE",
			contains: "for value in [1]\n    return value\nreturn none",
			excludes: "(\n    for value",
		},
		{
			name:     "direct returned grouped loop",
			input:    "RETURN ((FOR value IN [1] RETURN value))",
			contains: "return for value in [1]\n    return value",
			excludes: "return (",
		},
		{
			name:     "initializer keeps grouping",
			input:    "LET values = (FOR value IN [1] RETURN value)\nRETURN values",
			contains: "let values = (\n    for value in [1]",
		},
		{
			name: "nested grouped loop statements",
			input: `(FOR outer IN [1] {
  (FOR inner IN [outer] RETURN inner)
  RETURN outer
})
RETURN NONE`,
			contains: "for outer in [1] {\n    for inner in [outer]\n        return inner\n    return outer\n}\nreturn none",
			excludes: "(for",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := formatParenthesesStable(t, test.input)
			if !strings.Contains(got, test.contains) {
				t.Fatalf("formatted output does not contain %q:\n%s", test.contains, got)
			}

			if test.excludes != "" && strings.Contains(got, test.excludes) {
				t.Fatalf("formatted output unexpectedly contains %q:\n%s", test.excludes, got)
			}
		})
	}
}

func TestFormatterPreservesCommentsInsideRequiredParentheses(t *testing.T) {
	inputs := []string{
		"RETURN (/* leading */ 1 + 2)",
		"RETURN (1 + 2 /* trailing */)",
		"RETURN (\n// leading\n1 + 2\n)",
		"RETURN (1 + 2 // trailing\n)",
		"RETURN QUERY (/* payload */ prefix + selector) IN page",
		"RETURN (/* source */ 1 + 2).member",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			got := formatParenthesesStable(t, input)
			if !strings.Contains(got, "leading") && !strings.Contains(got, "trailing") &&
				!strings.Contains(got, "payload") && !strings.Contains(got, "source") {
				t.Fatalf("parenthesized comment was lost:\n%s", got)
			}
		})
	}
}

func TestFormatterParenthesisSimplificationPreservesExecution(t *testing.T) {
	inputs := []string{
		"RETURN ((1 + 2)) * 3",
		"RETURN 1 - (2 - 3)",
		"RETURN -(1 + 2)",
		"RETURN (NONE ?? 2) ?? 3",
		"RETURN 1 ?? (2 ?? 3)",
		"RETURN (FOR value IN [1, 2] { RETURN value * 2 })",
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, input := range inputs {
			t.Run(optimizationNameForFormatter(level)+"/"+input, func(t *testing.T) {
				formatted := formatParenthesesStable(t, input)
				originalProgram, err := spec.Compile(input, level)
				if err != nil {
					t.Fatalf("compile original: %v", err)
				}

				formattedProgram, err := spec.Compile(formatted, level)
				if err != nil {
					t.Fatalf("compile formatted: %v\n%s", err, formatted)
				}

				originalResult, err := spec.Run(originalProgram)
				if err != nil {
					t.Fatalf("run original: %v", err)
				}

				formattedResult, err := spec.Run(formattedProgram)
				if err != nil {
					t.Fatalf("run formatted: %v", err)
				}

				if !bytes.Equal(originalResult, formattedResult) {
					t.Fatalf("execution changed: original %s, formatted %s\n%s", originalResult, formattedResult, formatted)
				}
			})
		}
	}
}

func formatParenthesesStable(t *testing.T, input string) string {
	t.Helper()

	format := New()
	var first bytes.Buffer
	if err := format.Format(&first, source.NewAnonymous(input)); err != nil {
		t.Fatalf("format input: %v", err)
	}

	var second bytes.Buffer
	if err := format.Format(&second, source.NewAnonymous(first.String())); err != nil {
		t.Fatalf("format first output: %v\n%s", err, first.String())
	}

	if first.String() != second.String() {
		t.Fatalf("formatter is unstable:\nfirst:\n%s\nsecond:\n%s", first.String(), second.String())
	}

	return first.String()
}

func optimizationNameForFormatter(level compiler.OptimizationLevel) string {
	return fmt.Sprintf("O%d", level)
}
