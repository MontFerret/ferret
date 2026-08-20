package runtime

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

func dispatchAdd(ctx context.Context, left, right Value) (Value, error) {
	if receiver, ok := left.(Addable); ok {
		result, err := receiver.Add(ctx, right)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	if receiver, ok := right.(Addable); ok {
		result, err := receiver.RightAdd(ctx, left)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	return None, binaryOperatorTypeError(operator.Add, left, right)
}

func dispatchSubtract(ctx context.Context, left, right Value) (Value, error) {
	if receiver, ok := left.(Subtractable); ok {
		result, err := receiver.Subtract(ctx, right)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	if receiver, ok := right.(Subtractable); ok {
		result, err := receiver.RightSubtract(ctx, left)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	return None, binaryOperatorTypeError(operator.Subtract, left, right)
}

func dispatchMultiply(ctx context.Context, left, right Value) (Value, error) {
	if receiver, ok := left.(Multipliable); ok {
		result, err := receiver.Multiply(ctx, right)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	if receiver, ok := right.(Multipliable); ok {
		result, err := receiver.RightMultiply(ctx, left)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	return None, binaryOperatorTypeError(operator.Multiply, left, right)
}

func dispatchDivide(ctx context.Context, left, right Value) (Value, error) {
	if receiver, ok := left.(Dividable); ok {
		result, err := receiver.Divide(ctx, right)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	if receiver, ok := right.(Dividable); ok {
		result, err := receiver.RightDivide(ctx, left)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	return None, binaryOperatorTypeError(operator.Divide, left, right)
}

func dispatchMod(ctx context.Context, left, right Value) (Value, error) {
	if receiver, ok := left.(Modulable); ok {
		result, err := receiver.Mod(ctx, right)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	if receiver, ok := right.(Modulable); ok {
		result, err := receiver.RightMod(ctx, left)
		if err == nil {
			return result, nil
		}

		if !declinesOperands(err) {
			return None, err
		}
	}

	return None, binaryOperatorTypeError(operator.Modulus, left, right)
}

func declinesOperands(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	return errors.Is(err, ErrUnsupportedOperands)
}
