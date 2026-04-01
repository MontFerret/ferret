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
