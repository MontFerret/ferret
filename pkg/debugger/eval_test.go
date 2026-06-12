package debugger

import (
	"context"
	"testing"

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
