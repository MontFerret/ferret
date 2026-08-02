package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

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
		"LET delay = 1.5S",
		"TIMEOUT delay + 500ms",
		"EVERY 1e2MS",
		"DELAY (0ms OR 1ms) OR RETURN NONE",
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
	input := `RETURN MATCH 5s(5000MS=>true,_=>false)`
	var first bytes.Buffer
	fmt := New()

	if err := fmt.Format(&first, source.NewAnonymous(input)); err != nil {
		t.Fatalf("format failed: %v", err)
	}
	formatted := first.String()
	if !strings.Contains(formatted, "MATCH 5s") || !strings.Contains(formatted, "5000MS => TRUE") {
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
	expected := strings.TrimSpace(`RETURN [
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
	target := `DISPATCH "input" IN (QUERY ONE "#query" IN page USING css)`
	if targetIdx, withIdx := strings.Index(out, target), strings.Index(out, "WITH {"); targetIdx < 0 || withIdx < targetIdx+len(target) {
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
	input := `FUNC read( value )(
LET brand=value.product.brand
VAR price=value["prices"]["current"]
price=value.prices["sale"]
value.metadata.lastSeen
RETURN [ brand,price ]
)
RETURN read(@product)`
	src := source.NewAnonymous(input)
	var buf bytes.Buffer
	fmt := New()

	if err := fmt.Format(&buf, src); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	out := buf.String()
	expected := `FUNC read(value) (
    LET brand = value.product.brand
    VAR price = value["prices"]["current"]
    price = value.prices["sale"]
    value.metadata.lastSeen
    RETURN [brand, price]
)
RETURN read(@product)`
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
	if !strings.Contains(out, "WHEN .type == \"match\"") {
		t.Fatalf("expected WAITFOR event filter to use WHEN; got:\n%s", out)
	}
	if !strings.Contains(out, "WHEN .visible") {
		t.Fatalf("expected WAITFOR event filter to preserve repeated WHEN; got:\n%s", out)
	}
	if strings.Contains(out, "FILTER .type == \"match\"") {
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
	whenIdx := strings.Index(out, "WHEN .type == \"match\"")
	triggerIdx := strings.Index(out, "TRIGGER (")
	timeoutIdx := strings.Index(out, "TIMEOUT 1ms")
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
	whenIdx := strings.Index(out, "WHEN .type == \"match\"")
	triggerIdx := strings.Index(out, "TRIGGER button <- \"click\"")
	timeoutIdx := strings.Index(out, "TIMEOUT 1ms")
	if whenIdx < 0 || triggerIdx < 0 || timeoutIdx < 0 {
		t.Fatalf("expected inline WAITFOR trigger clauses in formatted output; got:\n%s", out)
	}
	if !(whenIdx < triggerIdx && triggerIdx < timeoutIdx) {
		t.Fatalf("expected WHEN -> TRIGGER -> TIMEOUT order; got:\n%s", out)
	}
	if strings.Contains(out, "TRIGGER (") {
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
	if !strings.Contains(out, "WHEN .ok WHEN .ok == TRUE") {
		t.Fatalf("expected WAITFOR predicate repeated WHEN clauses; got:\n%s", out)
	}

	var roundTrip bytes.Buffer
	if err := fmt.Format(&roundTrip, source.NewAnonymous(out)); err != nil {
		t.Fatalf("formatted output must remain parseable: %v\nformatted:\n%s", err, out)
	}
}
