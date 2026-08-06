package debugger

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

var debugArithmeticHostType = runtime.NewType("debugger", "ArithmeticHost", func(runtime.Value) bool {
	return true
})

type debugArithmeticHost struct {
	stringCalls int
}

func (v *debugArithmeticHost) String() string {
	v.stringCalls++

	return "host"
}

func (v *debugArithmeticHost) Hash() uint64 {
	return 1
}

func (v *debugArithmeticHost) Copy() runtime.Value {
	panic("debug arithmetic inspected Copy")
}

func (v *debugArithmeticHost) Type() runtime.Type {
	return debugArithmeticHostType
}

func TestEvaluateDebugArithmeticUsesRuntimeContract(t *testing.T) {
	t.Parallel()

	host := &debugArithmeticHost{}
	scope := evalScope{
		locals: map[string]runtime.Value{
			"array": runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
			"host":  host,
		},
		params: runtime.NewParams(),
		values: vm.NewDebugValueAccess(),
	}

	for expression, expected := range map[string]runtime.Value{
		`5.5 % 2`:     runtime.NewFloat(1.5),
		`"a" + array`: runtime.NewString("a[1,2]"),
		`host + "a"`:  runtime.NewString("hosta"),
	} {
		actual, err := evaluateExpression(context.Background(), expression, scope)
		if err != nil || actual != expected {
			t.Fatalf("%s = %v, %v; want %v", expression, actual, err, expected)
		}
	}

	if host.stringCalls != 1 {
		t.Fatalf("host String calls = %d, want 1", host.stringCalls)
	}

	for expression, expected := range map[string]string{
		`"10" - 2`:  "invalid operation: operator '-' cannot be applied to String and Int",
		`true * 2`:  "invalid operation: operator '*' cannot be applied to Boolean and Int",
		`array / 2`: "invalid operation: operator '/' cannot be applied to Array and Int",
		`host % 2`:  "invalid operation: operator '%' cannot be applied to debugger.ArithmeticHost and Int",
	} {
		_, err := evaluateExpression(context.Background(), expression, scope)
		if !errors.Is(err, runtime.ErrInvalidOperation) || err.Error() != expected {
			t.Fatalf("%s error = %v, want %q", expression, err, expected)
		}
	}

	if host.stringCalls != 1 {
		t.Fatalf("rejected host arithmetic called String; calls = %d", host.stringCalls)
	}
}
