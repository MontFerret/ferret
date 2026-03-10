package vm

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Comparator int
type ComparatorPredicate func(ctx context.Context, a, b runtime.Value) runtime.Boolean

const (
	EQ Comparator = iota
	NEQ
	GT
	GTE
	LT
	LTE
	IN
)

func ComparatorFromByte(op int) Comparator {
	val := Comparator(op)

	if val < EQ || val > IN {
		return -1
	}

	return val
}

func (op Comparator) Predicate() ComparatorPredicate {
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

func cmp(_ context.Context, left, right runtime.Value) runtime.Int {
	return runtime.Int(runtime.CompareValues(right, left))
}

func eq(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) == 0
}

func ne(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) != 0
}

func gt(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) > 0
}

func gte(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) >= 0
}

func lt(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) < 0
}

func lte(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) <= 0
}

func checkInclusion(ctx context.Context, left, right runtime.Value) runtime.Boolean {
	// If "left in right" -> right.contains(left)
	return contains(ctx, right, left)
}

func contains(ctx context.Context, input runtime.Value, value runtime.Value) runtime.Boolean {
	switch val := input.(type) {
	case runtime.List:
		idx, err := val.IndexOf(ctx, value)

		if err != nil {
			return runtime.False
		}

		return idx > -1
	case runtime.Map:
		containsValue, err := val.Contains(ctx, value)

		if err != nil {
			return runtime.False
		}

		return containsValue
	case runtime.String:
		return runtime.Boolean(strings.Contains(val.String(), value.String()))
	default:
		return false
	}
}
