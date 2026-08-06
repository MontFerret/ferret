package runtime

import (
	"context"
	"errors"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

// EvaluateComparison applies a module-internal comparison operator to two values.
// External callers should use EqualValues or CompareValues for stable comparison contracts.
func EvaluateComparison(ctx context.Context, op operator.Binary, left, right Value) (Boolean, error) {
	switch op {
	case operator.Equal:
		return EqualValues(ctx, left, right)
	case operator.NotEqual:
		return evaluateNotEqual(ctx, left, right)
	case operator.Less:
		return evaluateLess(ctx, left, right)
	case operator.LessOrEqual:
		return evaluateLessOrEqual(ctx, left, right)
	case operator.Greater:
		return evaluateGreater(ctx, left, right)
	case operator.GreaterOrEqual:
		return evaluateGreaterOrEqual(ctx, left, right)
	case operator.In:
		return containsValue(ctx, right, left)
	default:
		return False, unsupportedComparisonOperatorError(op)
	}
}

// ResolveComparison returns a module-internal predicate for repeated evaluation.
// External callers should use EqualValues or CompareValues for stable comparison contracts.
func ResolveComparison(
	op operator.Binary,
) (func(context.Context, Value, Value) (Boolean, error), error) {
	switch op {
	case operator.Equal:
		return EqualValues, nil
	case operator.NotEqual:
		return evaluateNotEqual, nil
	case operator.Less:
		return evaluateLess, nil
	case operator.LessOrEqual:
		return evaluateLessOrEqual, nil
	case operator.Greater:
		return evaluateGreater, nil
	case operator.GreaterOrEqual:
		return evaluateGreaterOrEqual, nil
	case operator.In:
		return evaluateIn, nil
	default:
		return nil, unsupportedComparisonOperatorError(op)
	}
}

func evaluateNotEqual(ctx context.Context, left, right Value) (Boolean, error) {
	equal, err := EqualValues(ctx, left, right)
	if err != nil {
		return False, err
	}

	return !equal, nil
}

func evaluateLess(ctx context.Context, left, right Value) (Boolean, error) {
	comparison, err := CompareValues(ctx, left, right)
	if err != nil {
		return False, comparisonOperatorError(operator.Less, left, right, err)
	}

	return comparison < Equal, nil
}

func evaluateLessOrEqual(ctx context.Context, left, right Value) (Boolean, error) {
	comparison, err := CompareValues(ctx, left, right)
	if err != nil {
		return False, comparisonOperatorError(operator.LessOrEqual, left, right, err)
	}

	return comparison <= Equal, nil
}

func evaluateGreater(ctx context.Context, left, right Value) (Boolean, error) {
	comparison, err := CompareValues(ctx, left, right)
	if err != nil {
		return False, comparisonOperatorError(operator.Greater, left, right, err)
	}

	return comparison > Equal, nil
}

func evaluateGreaterOrEqual(ctx context.Context, left, right Value) (Boolean, error) {
	comparison, err := CompareValues(ctx, left, right)
	if err != nil {
		return False, comparisonOperatorError(operator.GreaterOrEqual, left, right, err)
	}

	return comparison >= Equal, nil
}

func evaluateIn(ctx context.Context, left, right Value) (Boolean, error) {
	return containsValue(ctx, right, left)
}

func comparisonOperatorError(op operator.Binary, left, right Value, err error) error {
	if errors.Is(err, ErrInvalidOperation) {
		return binaryOperatorTypeError(op, left, right)
	}

	return err
}

func containsValue(ctx context.Context, input, value Value) (Boolean, error) {
	switch collection := input.(type) {
	case List:
		idx, err := collection.IndexOf(ctx, value)
		if err != nil {
			return False, err
		}

		return idx > -1, nil
	case Map:
		return collection.Contains(ctx, value)
	case String:
		return Boolean(strings.Contains(collection.String(), value.String())), nil
	default:
		return False, nil
	}
}

func unsupportedComparisonOperatorError(op operator.Binary) error {
	return Errorf(ErrInvalidOperation, "operator %q is not a comparison operator", op.String())
}
