package vm

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	arrayComparator          int
	arrayComparatorPredicate func(ctx context.Context, a, b runtime.Value) (runtime.Boolean, error)
)

const (
	EQ arrayComparator = iota
	NEQ
	GT
	GTE
	LT
	LTE
	IN
)

func comparatorFromByte(op int) arrayComparator {
	val := arrayComparator(op)

	if val < EQ || val > IN {
		return -1
	}

	return val
}

func (op arrayComparator) predicate() arrayComparatorPredicate {
	switch op {
	case EQ:
		return eq
	case NEQ:
		return ne
	case GT:
		return gt
	case GTE:
		return gte
	case LT:
		return lt
	case LTE:
		return lte
	case IN:
		return checkInclusion
	default:
		return eq
	}
}

func cmp(ctx context.Context, left, right runtime.Value) (runtime.Int, error) {
	result, err := runtime.CompareValues(ctx, right, left)
	return runtime.Int(result), err
}

func eq(ctx context.Context, left, right runtime.Value) (runtime.Boolean, error) {
	return runtime.EqualValues(ctx, left, right)
}

func ne(ctx context.Context, left, right runtime.Value) (runtime.Boolean, error) {
	result, err := runtime.EqualValues(ctx, left, right)
	if err != nil {
		return runtime.False, err
	}

	return !result, nil
}

func gt(ctx context.Context, left, right runtime.Value) (runtime.Boolean, error) {
	result, err := runtime.CompareValues(ctx, left, right)
	if err != nil {
		return runtime.False, relationalComparisonError(">", left, right, err)
	}

	return result > 0, nil
}

func gte(ctx context.Context, left, right runtime.Value) (runtime.Boolean, error) {
	result, err := runtime.CompareValues(ctx, left, right)
	if err != nil {
		return runtime.False, relationalComparisonError(">=", left, right, err)
	}

	return result >= 0, nil
}

func lt(ctx context.Context, left, right runtime.Value) (runtime.Boolean, error) {
	result, err := runtime.CompareValues(ctx, left, right)
	if err != nil {
		return runtime.False, relationalComparisonError("<", left, right, err)
	}

	return result < 0, nil
}

func lte(ctx context.Context, left, right runtime.Value) (runtime.Boolean, error) {
	result, err := runtime.CompareValues(ctx, left, right)
	if err != nil {
		return runtime.False, relationalComparisonError("<=", left, right, err)
	}

	return result <= 0, nil
}

func checkInclusion(ctx context.Context, left, right runtime.Value) (runtime.Boolean, error) {
	// If "left in right" -> right.contains(left)
	return contains(ctx, right, left)
}

func contains(ctx context.Context, input runtime.Value, value runtime.Value) (runtime.Boolean, error) {
	switch val := input.(type) {
	case runtime.List:
		idx, err := val.IndexOf(ctx, value)

		if err != nil {
			return runtime.False, err
		}

		return idx > -1, nil
	case runtime.Map:
		containsValue, err := val.Contains(ctx, value)

		if err != nil {
			return runtime.False, err
		}

		return containsValue, nil
	case runtime.String:
		return runtime.Boolean(strings.Contains(val.String(), value.String())), nil
	default:
		return false, nil
	}
}
