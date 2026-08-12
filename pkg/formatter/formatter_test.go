package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestFormatter_DefaultKeywordCase(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"uppercase": {
			input: `RETURN FOR x IN items { RETURN x }`,
			want:  "return for x in items {\n    return x\n}",
		},
		"lowercase": {
			input: "return for x in items {\n    return x\n}",
			want:  "return for x in items {\n    return x\n}",
		},
		"mixed case preserves symbols and contents": {
			input: `ReTuRn FoR Item In Items {
    FiLtEr Item.Status == "RETURN"
    LeT DISTINCT = Item.Status
    LeT Rows = QuErY "SELECT RETURN" In Item UsInG PgSQL
    // FOR remains uppercase in comments
    ReTuRn { TRIGGER: DISTINCT, Value: DB::POSTGRES::Map(Item), Rows: Rows }
}`,
			want: `return for Item in Items {
    filter Item.Status == "RETURN"
    let DISTINCT = Item.Status
    let Rows = query "SELECT RETURN" in Item using PgSQL
    // FOR remains uppercase in comments
    return { TRIGGER: DISTINCT, Value: DB::POSTGRES::Map(Item), Rows: Rows }
}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var first bytes.Buffer
			format := New()

			if err := format.Format(&first, source.NewAnonymous(test.input)); err != nil {
				t.Fatalf("format failed: %v", err)
			}

			if first.String() != test.want {
				t.Fatalf("unexpected default formatting:\nwant:\n%s\ngot:\n%s", test.want, first.String())
			}

			var second bytes.Buffer
			if err := format.Format(&second, source.NewAnonymous(first.String())); err != nil {
				t.Fatalf("second format failed: %v", err)
			}

			if second.String() != first.String() {
				t.Fatalf("default formatting must be stable:\nfirst:\n%s\nsecond:\n%s", first.String(), second.String())
			}
		})
	}
}

func TestFormatterPreservesFunctionIdentifierCase(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if err := New().Format(&output, source.NewAnonymous(`RETURN TO_STRING(value)`)); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	if got, want := output.String(), "return TO_STRING(value)"; got != want {
		t.Fatalf("formatter changed semantic function identifier: got %q, want %q", got, want)
	}
}

func TestFormatter_LiteralSpread(t *testing.T) {
	tests := []struct {
		name                   string
		input                  string
		contains               []string
		printWidth             int
		allowTriviaAfterSpread bool
		noBracketSpacing       bool
	}{
		{
			name:     "inline array and object",
			input:    `RETURN [... [1,2], {... {a:1}, b:2},]`,
			contains: []string{`return [`, `...[1, 2]`, `{ ...{ a: 1 }, b: 2 }`},
		},
		{
			name:                   "comments around spread entries",
			allowTriviaAfterSpread: true,
			input: `RETURN [
// array source
... /* spread source */ values, // copied values
1
]`,
			contains: []string{"// array source", "/* spread source */", "// copied values"},
		},
		{
			name:       "spread entries wrap at print width",
			input:      `RETURN { ...firstObject, explicitProperty: 1, ...secondObject }`,
			contains:   []string{"\n", "...firstObject", "...secondObject"},
			printWidth: 20,
		},
		{
			name:             "object spread without bracket spacing",
			input:            `RETURN { ...source }`,
			contains:         []string{`return {...source}`},
			noBracketSpacing: true,
		},
		{
			name:                   "comment between object spread and operand",
			input:                  `RETURN { ... /* defaults */ defaults, value: 1 }`,
			contains:               []string{"/* defaults */", "defaults"},
			allowTriviaAfterSpread: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options := make([]Option, 0, 2)
			if test.printWidth > 0 {
				options = append(options, WithPrintWidth(uint64(test.printWidth)))
			}

			if test.noBracketSpacing {
				options = append(options, WithBracketSpacing(false))
			}

			format := New(options...)
			var first bytes.Buffer

			if err := format.Format(&first, source.NewAnonymous(test.input)); err != nil {
				t.Fatalf("format failed: %v", err)
			}

			formatted := first.String()
			if !test.allowTriviaAfterSpread && strings.Contains(formatted, "... ") {
				t.Fatalf("spread operator must stay adjacent to its operand:\n%s", formatted)
			}

			for _, fragment := range test.contains {
				if !strings.Contains(formatted, fragment) {
					t.Fatalf("formatted spread literal missing %q:\n%s", fragment, formatted)
				}
			}

			var second bytes.Buffer
			if err := format.Format(&second, source.NewAnonymous(formatted)); err != nil {
				t.Fatalf("second format failed: %v", err)
			}

			if second.String() != formatted {
				t.Fatalf("spread formatting must be stable:\nfirst:\n%s\nsecond:\n%s", formatted, second.String())
			}
		})
	}
}

func TestFormatter_ExplicitUppercaseKeywordCase(t *testing.T) {
	var out bytes.Buffer
	format := New(WithCaseMode(CaseModeUpper))

	if err := format.Format(&out, source.NewAnonymous(`return for x in items { return x }`)); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	want := "RETURN FOR x IN items {\n    RETURN x\n}"
	if out.String() != want {
		t.Fatalf("unexpected uppercase formatting:\nwant:\n%s\ngot:\n%s", want, out.String())
	}
}

func TestFormatter_TemplateLiteralDoesNotIndentInterpolation(t *testing.T) {
	input := "RETURN { foo: `line1\n${1}`, veryLongPropertyNameThatForcesMultilineFormatting: 1 }"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New(WithPrintWidth(10))

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "line1\n${1}") {
		t.Fatalf("expected interpolation to start immediately after newline; got:\n%s", out)
	}
	if strings.Contains(out, "line1\n    ${1}") {
		t.Fatalf("unexpected indentation injected before interpolation; got:\n%s", out)
	}
}

func TestFormatter_DurationExpressionsRoundTrip(t *testing.T) {
	input := `LET delay=1.5S
RETURN WAITFOR FALSE TIMEOUT delay+500ms EVERY 1e2MS ON ERROR RETRY 1 DELAY (0ms OR 1ms) OR RETURN NONE`
	src := source.NewAnonymous(input)
	var first bytes.Buffer
	fmt := New()

	if err := fmt.Format(&first, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	formatted := first.String()
	for _, expected := range []string{
		"let delay = 1.5S",
		"timeout delay + 500ms",
		"every 1e2MS",
		"delay (0ms or 1ms) or return none",
	} {
		if !strings.Contains(formatted, expected) {
			t.Fatalf("formatted duration expression missing %q:\n%s", expected, formatted)
		}
	}

	var second bytes.Buffer
	if err := fmt.Format(&second, source.NewAnonymous(formatted)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v", err)
	}
	if second.String() != formatted {
		t.Fatalf("duration formatting must be stable:\nfirst:\n%s\nsecond:\n%s", formatted, second.String())
	}
}

func TestFormatter_DurationMatchLiteralRoundTrip(t *testing.T) {
	input := `RETURN MATCH 5s{5000MS=>true,_=>false}`
	var first bytes.Buffer
	fmt := New()

	if err := fmt.Format(&first, source.NewAnonymous(input)); err != nil {
		t.Fatalf("format failed: %v", err)
	}
	formatted := first.String()
	if !strings.Contains(formatted, "match 5s") || !strings.Contains(formatted, "5000MS => true") {
		t.Fatalf("duration match literals were not preserved:\n%s", formatted)
	}

	var second bytes.Buffer
	if err := fmt.Format(&second, source.NewAnonymous(formatted)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v", err)
	}
	if second.String() != formatted {
		t.Fatalf("duration match formatting must be stable:\nfirst:\n%s\nsecond:\n%s", formatted, second.String())
	}
}

func TestFormatter_ArrayTemplateLiteralNewlineForcesMultiline(t *testing.T) {
	input := "RETURN [`line1\n${1}`]"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New(WithPrintWidth(200))

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "line1\n${1}") {
		t.Fatalf("expected newline in template literal; got:\n%s", out)
	}
	if strings.Contains(out, "line1 ${1}") {
		t.Fatalf("unexpected newline collapse in template literal; got:\n%s", out)
	}
}

func TestFormatter_NestedObjectRespectsPrintWidthAtLineStart(t *testing.T) {
	input := "RETURN [{ a: 1, bb: 2 }]"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New(WithPrintWidth(18))

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	expected := strings.TrimSpace(`return [
    {
        a: 1,
        bb: 2
    }
]`)

	if out != expected {
		t.Fatalf("unexpected nested object formatting:\n%s", out)
	}

	for _, line := range strings.Split(out, "\n") {
		if len(line) > 18 {
			t.Fatalf("line exceeds print width 18 (%d): %q", len(line), line)
		}
	}
}

func TestFormatter_BlockCommentPreservesLeadingSpace(t *testing.T) {
	input := "LET x = 1\n/*\n * a\n * b\n */\nRETURN 2"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\n * a\n") {
		t.Fatalf("expected leading space in block comment line; got:\n%s", out)
	}
	if !strings.Contains(out, "\n * b\n") {
		t.Fatalf("expected leading space in block comment line; got:\n%s", out)
	}
}

func TestFormatter_DispatchGroupedQueryTargetRemainsParseable(t *testing.T) {
	input := "DISPATCH \"input\" IN (QUERY ONE \"#query\" IN page USING css) WITH { value: \"ferret\" }\nRETURN 1"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	target := `dispatch "input" in (query one "#query" in page using css)`
	if targetIdx, withIdx := strings.Index(out, target), strings.Index(out, "with {"); targetIdx < 0 || withIdx < targetIdx+len(target) {
		t.Fatalf("expected grouped query target and dispatch payload to remain distinct; got:\n%s", out)
	}

	var roundTrip bytes.Buffer
	if err := fmt.Format(&roundTrip, source.NewAnonymous(out)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v\nformatted:\n%s", err, out)
	}

	if roundTrip.String() != out {
		t.Fatalf("formatted output must be stable:\nfirst:\n%s\nsecond:\n%s", out, roundTrip.String())
	}
}

func TestFormatter_UdfMemberStatementsRemainUnparenthesizedAndParseable(t *testing.T) {
	input := `FUNC read( value ){
LET brand=value.product.brand
VAR price=value["prices"]["current"]
price=value.prices["sale"]
value.metadata.lastSeen
RETURN [ brand,price ]
}
RETURN read(@product)`
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	expected := `func read(value) {
    let brand = value.product.brand
    var price = value["prices"]["current"]
    price = value.prices["sale"]
    value.metadata.lastSeen
    return [brand, price]
}
return read(@product)`
	if out != expected {
		t.Fatalf("unexpected UDF member statement formatting:\nexpected:\n%s\nactual:\n%s", expected, out)
	}

	var roundTrip bytes.Buffer
	if err := fmt.Format(&roundTrip, source.NewAnonymous(out)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v\nformatted:\n%s", err, out)
	}

	if roundTrip.String() != out {
		t.Fatalf("formatted output must be stable:\nfirst:\n%s\nsecond:\n%s", out, roundTrip.String())
	}
}

func TestFormatter_WaitForEventFilterUsesWhenAndRemainsParseable(t *testing.T) {
	input := "LET obs = []\nWAITFOR EVENT \"test\" IN obs WHEN .type == \"match\" WHEN .visible\nRETURN 1"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "when .type == \"match\"") {
		t.Fatalf("expected WAITFOR event filter to use WHEN; got:\n%s", out)
	}
	if !strings.Contains(out, "when .visible") {
		t.Fatalf("expected WAITFOR event filter to preserve repeated WHEN; got:\n%s", out)
	}
	if strings.Contains(out, "filter .type == \"match\"") {
		t.Fatalf("unexpected legacy FILTER in WAITFOR event filter; got:\n%s", out)
	}

	var roundTrip bytes.Buffer
	if err := fmt.Format(&roundTrip, source.NewAnonymous(out)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v\nformatted:\n%s", err, out)
	}
}

func TestFormatter_WaitForEventTriggerRemainsParseable(t *testing.T) {
	input := "LET obs = []\nLET button = @button\nWAITFOR EVENT \"test\" IN obs WHEN .type == \"match\" TRIGGER (button <- \"click\") TIMEOUT 1ms\nRETURN 1"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	whenIdx := strings.Index(out, "when .type == \"match\"")
	triggerIdx := strings.Index(out, "trigger (")
	timeoutIdx := strings.Index(out, "timeout 1ms")
	if whenIdx < 0 || triggerIdx < 0 || timeoutIdx < 0 {
		t.Fatalf("expected WAITFOR trigger clauses in formatted output; got:\n%s", out)
	}
	if !(whenIdx < triggerIdx && triggerIdx < timeoutIdx) {
		t.Fatalf("expected WHEN -> TRIGGER -> TIMEOUT order; got:\n%s", out)
	}
	if !strings.Contains(out, "\n    button <- \"click\"\n") {
		t.Fatalf("expected trigger body to be formatted as a block; got:\n%s", out)
	}

	var roundTrip bytes.Buffer
	if err := fmt.Format(&roundTrip, source.NewAnonymous(out)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v\nformatted:\n%s", err, out)
	}
}

func TestFormatter_WaitForEventInlineTriggerRemainsInline(t *testing.T) {
	input := "LET obs = []\nLET button = @button\nWAITFOR EVENT \"test\" IN obs WHEN .type == \"match\" TRIGGER button <- \"click\" TIMEOUT 1ms\nRETURN 1"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	whenIdx := strings.Index(out, "when .type == \"match\"")
	triggerIdx := strings.Index(out, "trigger button <- \"click\"")
	timeoutIdx := strings.Index(out, "timeout 1ms")
	if whenIdx < 0 || triggerIdx < 0 || timeoutIdx < 0 {
		t.Fatalf("expected inline WAITFOR trigger clauses in formatted output; got:\n%s", out)
	}
	if !(whenIdx < triggerIdx && triggerIdx < timeoutIdx) {
		t.Fatalf("expected WHEN -> TRIGGER -> TIMEOUT order; got:\n%s", out)
	}
	if strings.Contains(out, "trigger (") {
		t.Fatalf("expected trigger shorthand to remain inline; got:\n%s", out)
	}

	var roundTrip bytes.Buffer
	if err := fmt.Format(&roundTrip, source.NewAnonymous(out)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v\nformatted:\n%s", err, out)
	}
}

func TestFormatter_WaitForPredicateRepeatedWhenRemainsParseable(t *testing.T) {
	input := "LET value = WAITFOR VALUE { ok: true } WHEN .ok WHEN .ok == true TIMEOUT 1ms\nRETURN value"
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "when .ok when .ok == true") {
		t.Fatalf("expected WAITFOR predicate repeated WHEN clauses; got:\n%s", out)
	}

	var roundTrip bytes.Buffer
	if err := fmt.Format(&roundTrip, source.NewAnonymous(out)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v\nformatted:\n%s", err, out)
	}
}
