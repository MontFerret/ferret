package internal

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestExpressionFormatter_UnaryNot(t *testing.T) {
	input := "return not a"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "not a" {
		t.Fatalf("unexpected unary operator formatting: %q", got)
	}
}

func TestExpressionFormatter_ImplicitMemberExpression(t *testing.T) {
	input := "return [1][* return .name]"
	program := parseProgram(t, input)
	inlineRet := mustFirst[*fql.InlineReturnContext](t, program)
	expr := inlineRet.Expression().(*fql.ExpressionContext)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != ".name" {
		t.Fatalf("unexpected implicit member formatting: %q", got)
	}
}

func TestExpressionFormatter_ImplicitMemberExpressionOptional(t *testing.T) {
	input := "return [1][* return ?.name]"
	program := parseProgram(t, input)
	inlineRet := mustFirst[*fql.InlineReturnContext](t, program)
	expr := inlineRet.Expression().(*fql.ExpressionContext)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "?.name" {
		t.Fatalf("unexpected implicit optional member formatting: %q", got)
	}
}

func TestExpressionFormatter_ImplicitCurrentExpression(t *testing.T) {
	input := "return [1][* return .]"
	program := parseProgram(t, input)
	inlineRet := mustFirst[*fql.InlineReturnContext](t, program)
	expr := inlineRet.Expression().(*fql.ExpressionContext)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "." {
		t.Fatalf("unexpected implicit current formatting: %q", got)
	}
}

func TestExpressionFormatter_RangeOperandImplicitCurrentExpression(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "return [1][* return . .. 10]",
			want:  "...10",
		},
		{
			input: "return [1][* return 1 .. .]",
			want:  "1...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			inlineRet := mustFirst[*fql.InlineReturnContext](t, program)
			expr := inlineRet.Expression().(*fql.ExpressionContext)

			var buf bytes.Buffer
			e := newEngine(source.NewAnonymous(tt.input), &buf, defaultTestConfig())

			e.expression.formatExpression(expr)
			if got := buf.String(); got != tt.want {
				t.Fatalf("unexpected range formatting: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExpressionFormatter_QueryExpressionInline(t *testing.T) {
	input := "return query `.items` in doc using css with { limit: 10 }"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query `.items` in doc using css with { limit: 10 }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionWithoutUsingInline(t *testing.T) {
	input := "return query `.items` in doc with { limit: 10 }"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query `.items` in doc with { limit: 10 }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionOptionsInline(t *testing.T) {
	input := "return query `.items` in doc using css options { timeout: 5000 }"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query `.items` in doc using css options { timeout: 5000 }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionWithoutUsingMultiline(t *testing.T) {
	input := "return query one `.items` in doc with { limit: 10, timeout: 5, extra: 1 } options { retry: 2, delay: 50 }"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	opts := defaultTestConfig()
	opts.PrintWidth = 20
	e := newEngine(source.NewAnonymous(input), &buf, opts)

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query one `.items` in doc\n    with {\n        limit: 10,\n        timeout: 5,\n        extra: 1\n    }\n    options {\n        retry: 2,\n        delay: 50\n    }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionParamPayload(t *testing.T) {
	input := "return query @q in doc using css"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query @q in doc using css" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionMemberPayload(t *testing.T) {
	input := "return query one email.body in model using summarize"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query one email.body in model using summarize" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionAtomicPayloads(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "return query selectors[index] in page using css",
			want:  "query selectors[index] in page using css",
		},
		{
			input: "return query GET_SELECTOR() in page using css",
			want:  "query GET_SELECTOR() in page using css",
		},
		{
			input: "return query factory().selector in page using css",
			want:  "query factory().selector in page using css",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expr := mustFirst[*fql.ExpressionContext](t, program)

			var buf bytes.Buffer
			e := newEngine(source.NewAnonymous(tt.input), &buf, defaultTestConfig())

			e.expression.formatExpression(expr)
			if got := buf.String(); got != tt.want {
				t.Fatalf("unexpected query expression formatting: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExpressionFormatter_QueryExpressionComputedPayload(t *testing.T) {
	input := `return query (prefix+selector) in page using css with { visible: true } options { timeout: 1000 }`
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query (prefix + selector) in page using css\n    with { visible: true }\n    options { timeout: 1000 }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionComputedPayloadMultilineClauses(t *testing.T) {
	input := `return query (prefix+selector) in page with { visible: true, enabled: true } options { timeout: 1000, retries: 2 }`
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	opts := defaultTestConfig()
	opts.PrintWidth = 20
	e := newEngine(source.NewAnonymous(input), &buf, opts)

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query (prefix + selector) in page\n    with {\n        visible: true,\n        enabled: true\n    }\n    options {\n        timeout: 1000,\n        retries: 2\n    }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionImplicitCurrentSource(t *testing.T) {
	input := `return [1][* return (query "a" in . using css)]`
	program := parseProgram(t, input)
	query := mustFirst[*fql.QueryExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatQueryExpression(query)
	if got := buf.String(); got != `query "a" in . using css` {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionCountModifier(t *testing.T) {
	input := "return query count `.items` in doc using css"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query count `.items` in doc using css" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryShorthand(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: `return doc[~"h1"]`,
			want:  `doc[~ "h1"]`,
		},
		{
			input: `return doc[~?"h1"]`,
			want:  `doc[~? "h1"]`,
		},
		{
			input: "return doc[~css`h1`]",
			want:  "doc[~ css`h1`]",
		},
		{
			input: "return doc[~?css`h1`]",
			want:  "doc[~? css`h1`]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			expr := mustFirst[*fql.ExpressionContext](t, program)

			var buf bytes.Buffer
			e := newEngine(source.NewAnonymous(tt.input), &buf, defaultTestConfig())

			e.expression.formatExpression(expr)
			if got := buf.String(); got != tt.want {
				t.Fatalf("unexpected query shorthand formatting: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExpressionFormatter_FunctionCallErrorPolicyTail(t *testing.T) {
	input := "return FAIL() on error return none"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "FAIL() on error return none" {
		t.Fatalf("unexpected function call error policy formatting: %q", got)
	}
}

func TestExpressionFormatter_FunctionCallRetryPolicyTail(t *testing.T) {
	input := "return FAIL() on error retry 3 delay -100MS backoff EXPONENTIAL or return none"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "FAIL() on error retry 3 delay -100MS backoff EXPONENTIAL or return none" {
		t.Fatalf("unexpected function call retry formatting: %q", got)
	}
}

func TestExpressionFormatter_ParenthesizedErrorPolicyTail(t *testing.T) {
	input := "return (FAIL() + 1) on error return none"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "(FAIL() + 1) on error return none" {
		t.Fatalf("unexpected grouped error policy formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionErrorPolicyTail(t *testing.T) {
	input := "return query `.items` in doc using css options { timeout: 5000 } on error return none"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query `.items` in doc using css options { timeout: 5000 } on error return none" {
		t.Fatalf("unexpected query error policy formatting: %q", got)
	}
}

func TestExpressionFormatter_QueryExpressionOneModifierWithMultiline(t *testing.T) {
	input := "return query one `.items` in doc using css with { limit: 10, timeout: 5 } options { retry: 2, delay: 50 }"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	opts := defaultTestConfig()
	opts.PrintWidth = 20
	e := newEngine(source.NewAnonymous(input), &buf, opts)

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "query one `.items` in doc using css\n    with {\n        limit: 10,\n        timeout: 5\n    }\n    options {\n        retry: 2,\n        delay: 50\n    }" {
		t.Fatalf("unexpected query expression formatting: %q", got)
	}
}

func TestExpressionFormatter_MatchExpressionInline(t *testing.T) {
	input := "return match x{1=>10,_=>0}"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "match x { 1 => 10, _ => 0 }" {
		t.Fatalf("unexpected match inline formatting: %q", got)
	}
}

func TestExpressionFormatter_MatchExpressionGuardMultiline(t *testing.T) {
	input := "return match{when a>0=>a,when a<0=>-a,_=>0}"
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	opts := defaultTestConfig()
	opts.PrintWidth = 10
	e := newEngine(source.NewAnonymous(input), &buf, opts)

	e.expression.formatExpression(expr)
	if got := buf.String(); got != "match {\n    when a > 0 => a,\n    when a < 0 => -a,\n    _ => 0,\n}" {
		t.Fatalf("unexpected match guard multiline formatting: %q", got)
	}
}

func TestExpressionFormatter_MatchExpressionObjectPattern(t *testing.T) {
	input := `return match obj{{ "a": 1, b: v }=>v,_=>0}`
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != `match obj { { "a": 1, b: v } => v, _ => 0 }` {
		t.Fatalf("unexpected match object pattern formatting: %q", got)
	}
}

func TestExpressionFormatter_MatchExpressionTriggerObjectPattern(t *testing.T) {
	input := `return match obj{{ TRIGGER: v }=>v,_=>0}`
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != `match obj { { TRIGGER: v } => v, _ => 0 }` {
		t.Fatalf("unexpected match object pattern formatting: %q", got)
	}
}

func TestExpressionFormatter_MatchExpressionDispatchShorthand(t *testing.T) {
	input := `return match kind{"click"=>btn<-"click",_=>input<-"focus"}`
	program := parseProgram(t, input)
	expr := mustFirst[*fql.ExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.expression.formatExpression(expr)
	if got := buf.String(); got != `match kind { "click" => btn <- "click", _ => input <- "focus" }` {
		t.Fatalf("unexpected match dispatch shorthand formatting: %q", got)
	}
}
