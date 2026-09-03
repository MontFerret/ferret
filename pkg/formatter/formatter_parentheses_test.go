package formatter

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/parser"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/mock"
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

func TestFormatterWaitForEventOperands(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       string
		wantName   string
		wantSource string
	}{
		{
			name:       "composed operands",
			input:      `RETURN WAITFOR EVENT (@kind ?? "message") IN (@source ?? fallback)`,
			want:       `return waitfor event @kind ?? "message" in @source ?? fallback`,
			wantName:   `@kind??"message"`,
			wantSource: `@source??fallback`,
		},
		{
			name:       "membership name removes redundant grouping",
			input:      `RETURN WAITFOR EVENT (@kind IN @names) IN @source`,
			want:       `return waitfor event @kind in @names in @source`,
			wantName:   `@kindin@names`,
			wantSource: `@source`,
		},
		{
			name:       "membership source retains delimiter boundary",
			input:      `RETURN WAITFOR EVENT "message" IN (@candidate IN @sources)`,
			want:       `return waitfor event "message" in (@candidate in @sources)`,
			wantName:   `"message"`,
			wantSource: `(@candidatein@sources)`,
		},
		{
			name:       "negated membership source retains delimiter boundary",
			input:      `RETURN WAITFOR EVENT "message" IN (@candidate NOT IN @sources)`,
			want:       `return waitfor event "message" in (@candidate not in @sources)`,
			wantName:   `"message"`,
			wantSource: `(@candidatenotin@sources)`,
		},
		{
			name:       "function argument membership is bounded",
			input:      `RETURN WAITFOR EVENT "message" IN (SOURCE(@candidate IN @sources))`,
			want:       `return waitfor event "message" in SOURCE(@candidate in @sources)`,
			wantName:   `"message"`,
			wantSource: `SOURCE(@candidatein@sources)`,
		},
		{
			name:       "array entry membership is bounded",
			input:      `RETURN WAITFOR EVENT "message" IN ([@candidate IN @sources])`,
			want:       `return waitfor event "message" in [@candidate in @sources]`,
			wantName:   `"message"`,
			wantSource: `[@candidatein@sources]`,
		},
		{
			name:       "object property membership is bounded",
			input:      `RETURN WAITFOR EVENT "message" IN ({ candidate: @candidate IN @sources })`,
			want:       `return waitfor event "message" in { candidate: @candidate in @sources }`,
			wantName:   `"message"`,
			wantSource: `{candidate:@candidatein@sources}`,
		},
		{
			name:       "ternary true branch membership is bounded",
			input:      `RETURN WAITFOR EVENT "message" IN (TRUE ? @candidate IN @sources : FALSE)`,
			want:       `return waitfor event "message" in true ? @candidate in @sources : false`,
			wantName:   `"message"`,
			wantSource: `true?@candidatein@sources:false`,
		},
		{
			name:       "precedence grouping bounds nested membership",
			input:      `RETURN WAITFOR EVENT "message" IN ((@candidate IN @sources) + 1)`,
			want:       `return waitfor event "message" in (@candidate in @sources) + 1`,
			wantName:   `"message"`,
			wantSource: `(@candidatein@sources)+1`,
		},
		{
			name:       "removable inner grouping leaves membership exposed",
			input:      `RETURN WAITFOR EVENT "message" IN (NOT (@candidate IN @sources))`,
			want:       `return waitfor event "message" in (not @candidate in @sources)`,
			wantName:   `"message"`,
			wantSource: `(not@candidatein@sources)`,
		},
		{
			name:       "ternary condition membership remains exposed",
			input:      `RETURN WAITFOR EVENT "message" IN (@candidate IN @sources ? TRUE : FALSE)`,
			want:       `return waitfor event "message" in (@candidate in @sources ? true : false)`,
			wantName:   `"message"`,
			wantSource: `(@candidatein@sources?true:false)`,
		},
		{
			name:       "ternary false branch membership remains exposed",
			input:      `RETURN WAITFOR EVENT "message" IN (TRUE ? FALSE : @candidate IN @sources)`,
			want:       `return waitfor event "message" in (true ? false : @candidate in @sources)`,
			wantName:   `"message"`,
			wantSource: `(true?false:@candidatein@sources)`,
		},
		{
			name:       "membership operands use distinct grouping",
			input:      `RETURN WAITFOR EVENT (@kind IN @names) IN (@candidate IN @sources)`,
			want:       `return waitfor event @kind in @names in (@candidate in @sources)`,
			wantName:   `@kindin@names`,
			wantSource: `(@candidatein@sources)`,
		},
		{
			name:       "array membership source retains delimiter boundary",
			input:      `RETURN WAITFOR EVENT "message" IN (@candidate ANY IN @sources)`,
			want:       `return waitfor event "message" in (@candidate any in @sources)`,
			wantName:   `"message"`,
			wantSource: `(@candidateanyin@sources)`,
		},
		{
			name:       "query source",
			input:      `RETURN WAITFOR EVENT GET_NAME(@id) IN QUERY ONE ".source" IN @registry`,
			want:       `return waitfor event GET_NAME(@id) in query one ".source" in @registry`,
			wantName:   `GET_NAME(@id)`,
			wantSource: `queryone".source"in@registry`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := formatParenthesesStable(t, test.input)
			if got != test.want {
				t.Fatalf("formatted output:\n%s\nwant:\n%s", got, test.want)
			}

			assertWaitForEventOperandBoundaries(t, got, test.wantName, test.wantSource)
		})
	}
}

func assertWaitForEventOperandBoundaries(t *testing.T, input, wantName, wantSource string) {
	t.Helper()

	p := parser.New(input)
	program := p.Program()
	if !p.AtEOF() {
		t.Fatalf("formatted output did not parse completely:\n%s", input)
	}

	event := findFirstWaitForEvent(program)
	if event == nil {
		t.Fatalf("formatted output has no WAITFOR EVENT expression:\n%s", input)
	}

	name := event.WaitForEventName().(*fql.WaitForEventNameContext).Expression()
	if got := name.GetText(); got != wantName {
		t.Fatalf("reparsed event-name expression = %q, want %q", got, wantName)
	}

	source := event.WaitForEventSource().(*fql.WaitForEventSourceContext).Expression()
	if got := source.GetText(); got != wantSource {
		t.Fatalf("reparsed source expression = %q, want %q", got, wantSource)
	}
}

func findFirstWaitForEvent(tree antlr.Tree) *fql.WaitForEventExpressionContext {
	if tree == nil {
		return nil
	}

	if event, ok := tree.(*fql.WaitForEventExpressionContext); ok {
		return event
	}

	for i := 0; i < tree.GetChildCount(); i++ {
		if event := findFirstWaitForEvent(tree.GetChild(i)); event != nil {
			return event
		}
	}

	return nil
}

func TestFormatterWaitForEventOperandsPreserveExecution(t *testing.T) {
	tests := []struct {
		newEnvironment func() []vm.EnvironmentOption
		name           string
		input          string
	}{
		{
			name: "composed operands",
			input: `RETURN WAITFOR EVENT (@eventName ?? "message") IN (@source ?? @fallback)
WHEN .type == "match"`,
			newEnvironment: func() []vm.EnvironmentOption {
				source := mock.NewObservable([]runtime.Value{
					mock.NewTestEventType("ignored"),
					mock.NewTestEventType("match"),
				})

				return []vm.EnvironmentOption{
					vm.WithParams(map[string]runtime.Value{
						"eventName": runtime.None,
						"fallback":  runtime.None,
						"source":    source,
					}),
				}
			},
		},
		{
			name: "bounded membership source",
			input: `RETURN WAITFOR EVENT "message" IN (SOURCE(@candidate IN @sources))
WHEN .type == "match"`,
			newEnvironment: func() []vm.EnvironmentOption {
				source := mock.NewObservable([]runtime.Value{
					mock.NewTestEventType("ignored"),
					mock.NewTestEventType("match"),
				})

				return []vm.EnvironmentOption{
					vm.WithParams(map[string]runtime.Value{
						"candidate": runtime.NewString("candidate"),
						"sources":   runtime.NewArrayWith(runtime.NewString("candidate")),
					}),
					vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
						return source, nil
					}),
				}
			},
		},
	}

	for _, test := range tests {
		formatted := formatParenthesesStable(t, test.input)

		for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Full} {
			t.Run(test.name+"/"+optimizationNameForFormatter(level), func(t *testing.T) {
				originalProgram, err := spec.Compile(test.input, level)
				if err != nil {
					t.Fatalf("compile original: %v", err)
				}

				formattedProgram, err := spec.Compile(formatted, level)
				if err != nil {
					t.Fatalf("compile formatted: %v\n%s", err, formatted)
				}

				originalResult, err := spec.Run(originalProgram, test.newEnvironment()...)
				if err != nil {
					t.Fatalf("run original: %v", err)
				}

				formattedResult, err := spec.Run(formattedProgram, test.newEnvironment()...)
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
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "nested arithmetic grouping",
			input: "RETURN ((1 + 2)) * 3",
			want:  "return (1 + 2) * 3",
		},
		{
			name:  "right subtraction operand",
			input: "RETURN 1 - (2 - 3)",
			want:  "return 1 - (2 - 3)",
		},
		{
			name:  "nested conditional condition",
			input: "RETURN (TRUE ? FALSE : TRUE) ? 1 : 2",
			want:  "return true ? false : true ? 1 : 2",
		},
		{
			name:  "nested conditional true branch",
			input: "RETURN TRUE ? (FALSE ? 1 : 2) : 3",
			want:  "return true ? false ? 1 : 2 : 3",
		},
		{
			name:  "nested conditional false branch",
			input: "RETURN FALSE ? 1 : (TRUE ? 2 : 3)",
			want:  "return false ? 1 : (true ? 2 : 3)",
		},
		{
			name:  "left comparison grouping",
			input: "RETURN (1 < 2) == TRUE",
			want:  "return 1 < 2 == true",
		},
		{
			name:  "right comparison grouping",
			input: "RETURN TRUE == (1 < 2)",
			want:  "return true == (1 < 2)",
		},
		{
			name:  "in conditional",
			input: `RETURN (1 IN [1]) ? "yes" : "no"`,
			want:  `return 1 in [1] ? "yes" : "no"`,
		},
		{
			name:  "in equality operand",
			input: "RETURN TRUE == (1 IN [1])",
			want:  "return true == (1 in [1])",
		},
		{
			name:  "like conditional",
			input: `RETURN ("foo" LIKE "f*") ? "yes" : "no"`,
			want:  `return "foo" like "f*" ? "yes" : "no"`,
		},
		{
			name:  "like equality operand",
			input: `RETURN TRUE == ("foo" LIKE "f*")`,
			want:  `return true == ("foo" like "f*")`,
		},
		{
			name:  "regex match equality",
			input: `RETURN ("abc" =~ "^a") == TRUE`,
			want:  `return "abc" =~ "^a" == true`,
		},
		{
			name:  "regex non-match conditional",
			input: `RETURN ("abc" !~ "^z") ? "yes" : "no"`,
			want:  `return "abc" !~ "^z" ? "yes" : "no"`,
		},
		{
			name:  "left coalesce grouping",
			input: "RETURN (NONE ?? 2) ?? 3",
			want:  "return (none ?? 2) ?? 3",
		},
		{
			name:  "right coalesce grouping",
			input: "RETURN 1 ?? (2 ?? 3)",
			want:  "return 1 ?? 2 ?? 3",
		},
		{
			name:  "coalesce conditional condition",
			input: "RETURN (NONE ?? FALSE) ? 1 : 2",
			want:  "return none ?? false ? 1 : 2",
		},
		{
			name:  "conditional coalesce operand",
			input: "RETURN TRUE ?? (FALSE ? 1 : 2)",
			want:  "return true ?? (false ? 1 : 2)",
		},
		{
			name:  "unary not",
			input: "RETURN NOT TRUE",
			want:  "return not true",
		},
		{
			name:  "unary additive grouping",
			input: "RETURN -(1 + 2)",
			want:  "return -1 + 2",
		},
		{
			name:  "unary additive operand",
			input: "RETURN (-1) + 2",
			want:  "return (-1) + 2",
		},
		{
			name:  "adjacent unary token boundary",
			input: "RETURN -(-1)",
			want:  "return -(-1)",
		},
		{
			name:  "mixed grammar precedence",
			input: `RETURN (NONE ?? "foo") LIKE "f*" ? (1 IN [1]) : ("abc" !~ "^z")`,
			want:  `return (none ?? "foo") like "f*" ? 1 in [1] : "abc" !~ "^z"`,
		},
		{
			name:  "grouped collecting loop",
			input: "RETURN (FOR value IN [1, 2] { RETURN value * 2 })",
			want:  "return for value in [1, 2] {\n    return value * 2\n}",
		},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Full} {
		for _, test := range tests {
			t.Run(optimizationNameForFormatter(level)+"/"+test.name, func(t *testing.T) {
				formatted := formatParenthesesStable(t, test.input)
				if formatted != test.want {
					t.Fatalf("formatted output:\n%s\nwant:\n%s", formatted, test.want)
				}

				originalProgram, err := spec.Compile(test.input, level)
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

	format := mustNewFormatter(t)
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
	return level.String()
}
