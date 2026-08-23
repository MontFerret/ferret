package internal

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestStatementFormatter_DispatchEventNameString(t *testing.T) {
	input := "dispatch \"evt\" in target"
	program := parseProgram(t, input)
	eventName := mustFirst[*fql.DispatchEventNameContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.statement.formatDispatchEventName(eventName)
	if got := buf.String(); got != "\"evt\"" {
		t.Fatalf("unexpected dispatch event name formatting: %q", got)
	}
}

func TestStatementFormatter_DispatchExpressionShorthand(t *testing.T) {
	input := `target<-"click"`
	program := parseProgram(t, input+"\nRETURN 1")
	dispatchExpr := mustFirst[*fql.DispatchExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.statement.formatDispatchExpression(dispatchExpr)
	if got := buf.String(); got != `target <- "click"` {
		t.Fatalf("unexpected shorthand dispatch formatting: %q", got)
	}
}

func TestStatementFormatter_DispatchExpressionErrorPolicyTail(t *testing.T) {
	input := `dispatch "evt" in target on error return none`
	program := parseProgram(t, input+"\nRETURN 1")
	dispatchExpr := mustFirst[*fql.DispatchExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.statement.formatDispatchExpression(dispatchExpr)
	if got := buf.String(); got != `dispatch "evt" in target on error return none` {
		t.Fatalf("unexpected dispatch error policy formatting: %q", got)
	}
}

func TestStatementFormatter_WaitForExpressionErrorPolicyTail(t *testing.T) {
	input := `waitfor value ready on error fail`
	program := parseProgram(t, input+"\nRETURN 1")
	waitExpr := mustFirst[*fql.WaitForExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.statement.formatWaitForExpression(waitExpr)
	if got := buf.String(); got != `waitfor value ready on error fail` {
		t.Fatalf("unexpected waitfor error policy formatting: %q", got)
	}
}

func TestStatementFormatter_WaitForExpressionRecoveryTailCanonicalOrder(t *testing.T) {
	input := `waitfor value ready timeout 1MS on error fail on timeout return none`
	program := parseProgram(t, input+"\nRETURN 1")
	waitExpr := mustFirst[*fql.WaitForExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.statement.formatWaitForExpression(waitExpr)
	if got := buf.String(); got != `waitfor value ready timeout 1MS on timeout return none on error fail` {
		t.Fatalf("unexpected waitfor recovery formatting: %q", got)
	}
}

func TestStatementFormatter_WaitForExpressionPredicateWhenClause(t *testing.T) {
	input := `waitfor exists rows when LENGTH(.) >= 10 when LENGTH(.) > 0 timeout 1MS every 10MS`
	program := parseProgram(t, input+"\nRETURN 1")
	waitExpr := mustFirst[*fql.WaitForExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.statement.formatWaitForExpression(waitExpr)
	if got := buf.String(); got != `waitfor exists rows when LENGTH(.) >= 10 when LENGTH(.) > 0 timeout 1MS every 10MS` {
		t.Fatalf("unexpected waitfor predicate when formatting: %q", got)
	}
}

func TestStatementFormatter_WaitForExpressionRetryTailCanonicalOrder(t *testing.T) {
	input := `waitfor value ready timeout 1MS on error retry 3 delay 10MS or return none on timeout return false`
	program := parseProgram(t, input+"\nRETURN 1")
	waitExpr := mustFirst[*fql.WaitForExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.statement.formatWaitForExpression(waitExpr)
	if got := buf.String(); got != `waitfor value ready timeout 1MS on timeout return false on error retry 3 delay 10MS or return none` {
		t.Fatalf("unexpected waitfor retry recovery formatting: %q", got)
	}
}
