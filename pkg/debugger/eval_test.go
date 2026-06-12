package debugger

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type hostileDebugEvalValue struct{}

func (hostileDebugEvalValue) String() string      { panic("String called") }
func (hostileDebugEvalValue) Hash() uint64        { panic("Hash called") }
func (hostileDebugEvalValue) Copy() runtime.Value { panic("Copy called") }
func (hostileDebugEvalValue) Type() runtime.Type  { panic("Type called") }

func TestEvaluateDebugExpressionParsesFerretStringEscapes(t *testing.T) {
	value, err := evaluateExpression(context.Background(), `'a\n\t\\b'`, evalScope{
		params: runtime.NewParams(),
		values: vm.NewDebugValueAccess(),
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := value, runtime.NewString("a\n\t\\\\b"); got != want {
		t.Fatalf("unexpected string value: got %q, want %q", got, want)
	}
}

func TestEvaluateDebugExpressionRejectsOpaqueValuesWithoutCallingHostMethods(t *testing.T) {
	_, err := evaluateExpression(
		context.Background(),
		"opaque AND true",
		evalScope{
			locals: map[string]runtime.Value{"opaque": hostileDebugEvalValue{}},
			params: runtime.NewParams(),
			values: vm.NewDebugValueAccess(),
		},
	)
	if err == nil {
		t.Fatal("expected opaque value to be rejected")
	}
}

func TestUnsupportedDebugExpressionIncludesParsedText(t *testing.T) {
	p := parser.New("[1]")
	err := unsupportedDebugExpression(p.Expression())

	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected invalid operation, got %v", err)
	}
	if !strings.Contains(err.Error(), "[1]") {
		t.Fatalf("expected expression text in error, got %v", err)
	}
}

func TestUnsupportedDebugExpressionHandlesValuesWithoutText(t *testing.T) {
	const want = "invalid operation: expression is not supported by the safe debugger evaluator"

	for _, value := range []any{nil, struct{}{}} {
		err := unsupportedDebugExpression(value)

		if !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("expected invalid operation for %#v, got %v", value, err)
		}
		if err.Error() != want {
			t.Fatalf("unexpected error for %#v: got %q, want %q", value, err, want)
		}
	}
}
