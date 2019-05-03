package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type (
	ArrayOperatorType int
	ArrayOperator     struct {
		*baseOperator
		aotype     ArrayOperatorType
		comparator core.OperatorExpression
	}
)

const (
	ArrayOperatorTypeAll  ArrayOperatorType = 0
	ArrayOperatorTypeAny  ArrayOperatorType = 1
	ArrayOperatorTypeNone ArrayOperatorType = 2
)

func IsValidArrayOperatorType(aotype ArrayOperatorType) bool {
	switch aotype {
	case ArrayOperatorTypeAll, ArrayOperatorTypeAny, ArrayOperatorTypeNone:
		return true
	default:
		return false
	}
}

func ToIsValidArrayOperatorType(stype string) (ArrayOperatorType, error) {
	switch stype {
	case "ALL":
		return ArrayOperatorTypeAll, nil
	case "ANY":
		return ArrayOperatorTypeAny, nil
	case "NONE":
		return ArrayOperatorTypeNone, nil
	default:
		return ArrayOperatorType(-1), core.Error(core.ErrInvalidArgument, stype)
	}
}

func NewArrayOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	aotype ArrayOperatorType,
	comparator core.OperatorExpression,
) (*ArrayOperator, error) {
	if left == nil {
		return nil, core.Error(core.ErrMissedArgument, "left expression")
	}

	if right == nil {
		return nil, core.Error(core.ErrMissedArgument, "right expression")
	}

	if !IsValidArrayOperatorType(aotype) {
		return nil, core.Error(core.ErrInvalidArgument, "operator")
	}

	if comparator == nil {
		return nil, core.Error(core.ErrMissedArgument, "comparator expression")
	}

	base := &baseOperator{src, left, right}

	return &ArrayOperator{base, aotype, comparator}, nil
}

func (operator *ArrayOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return values.False, core.SourceError(operator.src, err)
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return values.False, core.SourceError(operator.src, err)
	}

	return operator.Eval(ctx, left, right)
}

func (operator *ArrayOperator) Eval(ctx context.Context, left, right core.Value) (core.Value, error) {
	err := core.ValidateType(left, types.Array)

	if err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	arr := left.(*values.Array)

	switch operator.aotype {
	case ArrayOperatorTypeAll:
		return operator.all(ctx, arr, right)
	case ArrayOperatorTypeAny:
		return operator.any(ctx, arr, right)
	default:
		return operator.none(ctx, arr, right)
	}
}

func (operator *ArrayOperator) all(ctx context.Context, arr *values.Array, value core.Value) (core.Value, error) {
	result := values.False
	var err error

	arr.ForEach(func(el core.Value, _ int) bool {
		out, e := operator.comparator.Eval(ctx, el, value)

		if e != nil {
			err = e
			return false
		}

		if out == values.True {
			result = values.True
		} else {
			result = values.False
			return false
		}

		return true
	})

	if err != nil {
		return values.False, err
	}

	return result, nil
}

func (operator *ArrayOperator) any(ctx context.Context, arr *values.Array, value core.Value) (core.Value, error) {
	result := values.False
	var err error

	arr.ForEach(func(el core.Value, _ int) bool {
		out, e := operator.comparator.Eval(ctx, el, value)

		if e != nil {
			err = e
			return false
		}

		if out == values.True {
			result = values.True
			return false
		}

		return true
	})

	if err != nil {
		return values.False, err
	}

	return result, nil
}

func (operator *ArrayOperator) none(ctx context.Context, arr *values.Array, value core.Value) (core.Value, error) {
	result := values.False
	var err error

	arr.ForEach(func(el core.Value, _ int) bool {
		out, e := operator.comparator.Eval(ctx, el, value)

		if e != nil {
			err = e
			return false
		}

		if out == values.False {
			result = values.True
		} else {
			result = values.False
			return false
		}

		return true
	})

	if err != nil {
		return values.False, err
	}

	return result, nil
}
