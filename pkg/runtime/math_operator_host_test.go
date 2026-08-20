package runtime_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type arithmeticHostValue struct {
	stringCalls int
}

func (v *arithmeticHostValue) String() string {
	v.stringCalls++

	return "host"
}

func (v *arithmeticHostValue) Hash() uint64 {
	return 1
}

func (v *arithmeticHostValue) Copy() runtime.Value {
	panic("arithmetic dispatch inspected Copy")
}

func (v *arithmeticHostValue) Type() runtime.Type {
	return comparisonContractType
}

func TestStringTriggeredHostConcatenation(t *testing.T) {
	t.Parallel()

	host := &arithmeticHostValue{}
	left, err := runtime.Add(t.Context(), runtime.NewString("left:"), host)
	if err != nil || left != runtime.NewString("left:host") {
		t.Fatalf("String + host = %v, %v", left, err)
	}

	right, err := runtime.Add(t.Context(), host, runtime.NewString(":right"))
	if err != nil || right != runtime.NewString("host:right") {
		t.Fatalf("host + String = %v, %v", right, err)
	}

	if host.stringCalls != 2 {
		t.Fatalf("String calls = %d, want 2", host.stringCalls)
	}
}

func TestArithmeticDoesNotInspectUnsupportedHostValues(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	opaque := &opaqueHostValue{}
	typeName := runtime.TypeName(runtime.TypeOf(opaque))
	tests := []struct {
		operation func() (runtime.Value, error)
		expected  string
		name      string
	}{
		{name: "addition", operation: func() (runtime.Value, error) { return runtime.Add(ctx, runtime.NewInt(1), opaque) }, expected: fmt.Sprintf("invalid operation: operator '+' cannot be applied to Int and %s", typeName)},
		{name: "subtraction", operation: func() (runtime.Value, error) { return runtime.Subtract(ctx, opaque, runtime.NewInt(1)) }, expected: fmt.Sprintf("invalid operation: operator '-' cannot be applied to %s and Int", typeName)},
		{name: "multiplication", operation: func() (runtime.Value, error) { return runtime.Multiply(ctx, opaque, runtime.NewInt(1)) }, expected: fmt.Sprintf("invalid operation: operator '*' cannot be applied to %s and Int", typeName)},
		{name: "division", operation: func() (runtime.Value, error) { return runtime.Divide(ctx, opaque, runtime.NewInt(1)) }, expected: fmt.Sprintf("invalid operation: operator '/' cannot be applied to %s and Int", typeName)},
		{name: "modulo", operation: func() (runtime.Value, error) { return runtime.Mod(ctx, opaque, runtime.NewInt(1)) }, expected: fmt.Sprintf("invalid operation: operator '%%' cannot be applied to %s and Int", typeName)},
		{name: "increment", operation: func() (runtime.Value, error) { return runtime.Increment(ctx, opaque) }, expected: fmt.Sprintf("invalid operation: operator '++' cannot be applied to %s", typeName)},
		{name: "decrement", operation: func() (runtime.Value, error) { return runtime.Decrement(ctx, opaque) }, expected: fmt.Sprintf("invalid operation: operator '--' cannot be applied to %s", typeName)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.operation()
			if !errors.Is(err, runtime.ErrInvalidOperation) || err.Error() != test.expected {
				t.Fatalf("error = %v, want %q", err, test.expected)
			}
		})
	}
}
