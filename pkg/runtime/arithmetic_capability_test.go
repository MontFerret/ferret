package runtime_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type arithmeticOperation func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error)

func TestArithmeticCapabilityDispatchDirections(t *testing.T) {
	t.Parallel()

	operations := []struct {
		operation   arithmeticOperation
		leftMethod  string
		rightMethod string
		name        string
	}{
		{name: "addition", operation: runtime.Add, leftMethod: "Add", rightMethod: "RightAdd"},
		{name: "subtraction", operation: runtime.Subtract, leftMethod: "Subtract", rightMethod: "RightSubtract"},
		{name: "multiplication", operation: runtime.Multiply, leftMethod: "Multiply", rightMethod: "RightMultiply"},
		{name: "division", operation: runtime.Divide, leftMethod: "Divide", rightMethod: "RightDivide"},
		{name: "modulus", operation: runtime.Modulo, leftMethod: "Mod", rightMethod: "RightMod"},
	}

	for _, operation := range operations {
		t.Run(operation.name+"/left", func(t *testing.T) {
			host := newArithmeticCapabilityValue("host")
			other := runtime.NewInt(10)
			ctx := context.WithValue(t.Context(), struct{}{}, operation.leftMethod)

			result, err := operation.operation(ctx, host, other)
			if err != nil || result != runtime.NewString(operation.leftMethod) {
				t.Fatalf("operation result = %v, %v", result, err)
			}

			assertArithmeticCapabilityCall(t, host, operation.leftMethod, ctx, other)
		})

		t.Run(operation.name+"/right", func(t *testing.T) {
			host := newArithmeticCapabilityValue("host")
			other := runtime.NewInt(10)
			ctx := context.WithValue(t.Context(), struct{}{}, operation.rightMethod)

			result, err := operation.operation(ctx, other, host)
			if err != nil || result != runtime.NewString(operation.rightMethod) {
				t.Fatalf("operation result = %v, %v", result, err)
			}

			assertArithmeticCapabilityCall(t, host, operation.rightMethod, ctx, other)
		})
	}
}

func TestArithmeticCapabilityNegotiation(t *testing.T) {
	t.Parallel()

	left := newArithmeticCapabilityValue("left")
	left.responses["Add"] = arithmeticCapabilityResponse{
		value: runtime.None,
		err:   fmt.Errorf("left declined: %w", runtime.ErrUnsupportedOperands),
	}
	right := newArithmeticCapabilityValue("right")
	right.responses["RightAdd"] = arithmeticCapabilityResponse{value: runtime.NewString("fallback")}

	result, err := runtime.Add(t.Context(), left, right)
	if err != nil || result != runtime.NewString("fallback") {
		t.Fatalf("Add() = %v, %v", result, err)
	}

	assertArithmeticCapabilityMethod(t, left, "Add")
	assertArithmeticCapabilityMethod(t, right, "RightAdd")
}

func TestArithmeticCapabilityPropagatesRealErrors(t *testing.T) {
	t.Parallel()

	hostErr := errors.New("currency mismatch")
	left := newArithmeticCapabilityValue("left")
	left.responses["Subtract"] = arithmeticCapabilityResponse{value: runtime.None, err: hostErr}
	right := newArithmeticCapabilityValue("right")

	result, err := runtime.Subtract(t.Context(), left, right)
	if result != runtime.None || !errors.Is(err, hostErr) {
		t.Fatalf("Subtract() = %v, %v", result, err)
	}

	if len(right.calls) != 0 {
		t.Fatalf("right capability calls = %v, want none", right.calls)
	}
}

func TestArithmeticCapabilityPropagatesCancellationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target error
		name   string
	}{
		{name: "canceled", target: context.Canceled},
		{name: "deadline", target: context.DeadlineExceeded},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			left := newArithmeticCapabilityValue("left")
			left.responses["Multiply"] = arithmeticCapabilityResponse{
				value: runtime.None,
				err: errors.Join(
					fmt.Errorf("host stopped: %w", test.target),
					runtime.ErrUnsupportedOperands,
				),
			}
			right := newArithmeticCapabilityValue("right")

			result, err := runtime.Multiply(t.Context(), left, right)
			if result != runtime.None || !errors.Is(err, test.target) {
				t.Fatalf("Multiply() = %v, %v", result, err)
			}

			if len(right.calls) != 0 {
				t.Fatalf("right capability calls = %v, want none", right.calls)
			}
		})
	}
}

func TestArithmeticCapabilityDeclinesToNormalOperatorError(t *testing.T) {
	t.Parallel()

	left := newArithmeticCapabilityValue("left")
	left.responses["Divide"] = arithmeticCapabilityResponse{value: runtime.None, err: runtime.ErrUnsupportedOperands}
	right := newArithmeticCapabilityValue("right")
	right.responses["RightDivide"] = arithmeticCapabilityResponse{value: runtime.None, err: runtime.ErrUnsupportedOperands}

	result, err := runtime.Divide(t.Context(), left, right)
	if result != runtime.None || !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("Divide() = %v, %v", result, err)
	}

	if errors.Is(err, runtime.ErrUnsupportedOperands) {
		t.Fatalf("Divide() leaked ErrUnsupportedOperands: %v", err)
	}

	expected := "invalid operation: operator '/' cannot be applied to runtime_test.ArithmeticCapabilityValue and runtime_test.ArithmeticCapabilityValue"
	if err.Error() != expected {
		t.Fatalf("Divide() error = %q, want %q", err, expected)
	}
}

func TestArithmeticPrimitiveCapabilitiesRemainIndependent(t *testing.T) {
	t.Parallel()

	type typedRuntimeValue interface {
		runtime.Value
		runtime.Typed
	}

	addImplementation := newArithmeticCapabilityValue("add")
	addOnly := struct {
		typedRuntimeValue
		runtime.Addable
	}{addImplementation, addImplementation}

	if _, ok := any(addOnly).(runtime.Subtractable); ok {
		t.Fatal("Addable value unexpectedly implements Subtractable")
	}

	result, err := runtime.Add(t.Context(), addOnly, runtime.NewInt(1))
	if err != nil || result != runtime.NewString("Add") {
		t.Fatalf("Add() = %v, %v", result, err)
	}

	if _, err := runtime.Subtract(t.Context(), addOnly, runtime.NewInt(1)); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("Subtract() error = %v, want ErrInvalidOperation", err)
	}

	multiplyImplementation := newArithmeticCapabilityValue("multiply")
	multiplyOnly := struct {
		typedRuntimeValue
		runtime.Multipliable
	}{multiplyImplementation, multiplyImplementation}

	if _, ok := any(multiplyOnly).(runtime.Dividable); ok {
		t.Fatal("Multipliable value unexpectedly implements Dividable")
	}

	if _, ok := any(multiplyOnly).(runtime.Modulable); ok {
		t.Fatal("Multipliable value unexpectedly implements Modulable")
	}

	result, err = runtime.Multiply(t.Context(), multiplyOnly, runtime.NewInt(2))
	if err != nil || result != runtime.NewString("Multiply") {
		t.Fatalf("Multiply() = %v, %v", result, err)
	}

	if _, err := runtime.Divide(t.Context(), multiplyOnly, runtime.NewInt(2)); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("Divide() error = %v, want ErrInvalidOperation", err)
	}

	if _, err := runtime.Modulo(t.Context(), multiplyOnly, runtime.NewInt(2)); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("Modulo() error = %v, want ErrInvalidOperation", err)
	}
}

func TestStringConcatenationPrecedesAddableDispatch(t *testing.T) {
	t.Parallel()

	host := newArithmeticCapabilityValue("host")

	left, err := runtime.Add(t.Context(), runtime.NewString("left:"), host)
	if err != nil || left != runtime.NewString("left:host") {
		t.Fatalf("String + host = %v, %v", left, err)
	}

	right, err := runtime.Add(t.Context(), host, runtime.NewString(":right"))
	if err != nil || right != runtime.NewString("host:right") {
		t.Fatalf("host + String = %v, %v", right, err)
	}

	if len(host.calls) != 0 {
		t.Fatalf("arithmetic capability calls = %v, want none", host.calls)
	}
}

func assertArithmeticCapabilityCall(
	t *testing.T,
	host *arithmeticCapabilityValue,
	method string,
	ctx context.Context,
	operand runtime.Value,
) {
	t.Helper()

	if len(host.calls) != 1 {
		t.Fatalf("capability calls = %v, want one", host.calls)
	}

	call := host.calls[0]
	if call.method != method || call.ctx != ctx || call.operand != operand {
		t.Fatalf("capability call = %#v, want method=%s context=%p operand=%v", call, method, ctx, operand)
	}
}

func assertArithmeticCapabilityMethod(t *testing.T, host *arithmeticCapabilityValue, method string) {
	t.Helper()

	if len(host.calls) != 1 || host.calls[0].method != method {
		t.Fatalf("capability calls = %v, want %s", host.calls, method)
	}
}
