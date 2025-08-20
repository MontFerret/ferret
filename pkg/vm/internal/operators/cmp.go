package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
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
		return Equals
	case NEQ:
		return NotEquals
	case GT:
		return GreaterThan
	case GTE:
		return GreaterThanOrEqual
	case LT:
		return LessThan
	case LTE:
		return LessThanOrEqual
	case IN:
		return In
	default:
		return Equals
	}
}

func Compare(_ context.Context, left, right runtime.Value) runtime.Int {
	return runtime.Int(runtime.CompareValues(right, left))
}

func Equals(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) == 0
}

func NotEquals(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) != 0
}

func GreaterThan(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) > 0
}

func GreaterThanOrEqual(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) >= 0
}

func LessThan(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) < 0
}

func LessThanOrEqual(_ context.Context, left, right runtime.Value) runtime.Boolean {
	return runtime.CompareValues(left, right) <= 0
}

func In(ctx context.Context, left, right runtime.Value) runtime.Boolean {
	// If "left in right" -> right.Contains(left)
	return Contains(ctx, right, left)
}
