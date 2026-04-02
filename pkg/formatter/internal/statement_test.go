package internal

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestStatementFormatter_DispatchEventNameString(t *testing.T) {
	input := "DISPATCH \"evt\" IN target"
	program := parseProgram(t, input)
	eventName := mustFirst[*fql.DispatchEventNameContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, DefaultOptions())

	e.statement.formatDispatchEventName(eventName)
	if got := buf.String(); got != "\"evt\"" {
		t.Fatalf("unexpected dispatch event name formatting: %q", got)
	}
}

func TestStatementFormatter_DispatchExpressionShorthand(t *testing.T) {
	input := `"click"->target`
	program := parseProgram(t, input+"\nRETURN 1")
	dispatchExpr := mustFirst[*fql.DispatchExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, DefaultOptions())

	e.statement.formatDispatchExpression(dispatchExpr)
	if got := buf.String(); got != `"click" -> target` {
		t.Fatalf("unexpected shorthand dispatch formatting: %q", got)
	}
}

func TestStatementFormatter_DispatchExpressionErrorPolicyTail(t *testing.T) {
	input := `DISPATCH "evt" IN target ON ERROR RETURN NONE`
	program := parseProgram(t, input+"\nRETURN 1")
	dispatchExpr := mustFirst[*fql.DispatchExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, DefaultOptions())

	e.statement.formatDispatchExpression(dispatchExpr)
	if got := buf.String(); got != `DISPATCH "evt" IN target ON ERROR RETURN NONE` {
		t.Fatalf("unexpected dispatch error policy formatting: %q", got)
	}
}

func TestStatementFormatter_WaitForExpressionErrorPolicyTail(t *testing.T) {
	input := `WAITFOR VALUE ready ON ERROR FAIL`
	program := parseProgram(t, input+"\nRETURN 1")
	waitExpr := mustFirst[*fql.WaitForExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, DefaultOptions())

	e.statement.formatWaitForExpression(waitExpr)
	if got := buf.String(); got != `WAITFOR VALUE ready ON ERROR FAIL` {
		t.Fatalf("unexpected waitfor error policy formatting: %q", got)
	}
}

func TestStatementFormatter_WaitForExpressionRecoveryTailCanonicalOrder(t *testing.T) {
	input := `WAITFOR VALUE ready TIMEOUT 1 ON ERROR FAIL ON TIMEOUT RETURN NONE`
	program := parseProgram(t, input+"\nRETURN 1")
	waitExpr := mustFirst[*fql.WaitForExpressionContext](t, program)

	var buf bytes.Buffer
	e := newEngine(source.NewAnonymous(input), &buf, DefaultOptions())

	e.statement.formatWaitForExpression(waitExpr)
	if got := buf.String(); got != `WAITFOR VALUE ready TIMEOUT 1 ON TIMEOUT RETURN NONE ON ERROR FAIL` {
		t.Fatalf("unexpected waitfor recovery formatting: %q", got)
	}
}
