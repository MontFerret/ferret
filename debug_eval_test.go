package ferret

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type hostileDebugEvalValue struct{}

func (hostileDebugEvalValue) String() string      { panic("String called") }
func (hostileDebugEvalValue) Hash() uint64        { panic("Hash called") }
func (hostileDebugEvalValue) Copy() runtime.Value { panic("Copy called") }
func (hostileDebugEvalValue) Type() runtime.Type  { panic("Type called") }

func TestEvaluateDebugExpressionParsesFerretStringEscapes(t *testing.T) {
	value, err := evaluateDebugExpression(context.Background(), `'a\n\t\\b'`, nil, runtime.NewParams())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := value, runtime.NewString("a\n\t\\\\b"); got != want {
		t.Fatalf("unexpected string value: got %q, want %q", got, want)
	}
}

func TestEvaluateDebugExpressionRejectsOpaqueValuesWithoutCallingHostMethods(t *testing.T) {
	_, err := evaluateDebugExpression(
		context.Background(),
		"opaque AND true",
		map[string]runtime.Value{"opaque": hostileDebugEvalValue{}},
		runtime.NewParams(),
	)
	if err == nil {
		t.Fatal("expected opaque value to be rejected")
	}
}
