package operators

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type (
	ArrayOperatorVariant int

	ArrayOperator struct {
		*baseOperator
		variant    ArrayOperatorVariant
		comparator core.Predicate
	}
)

const (
	ArrayOperatorVariantAll ArrayOperatorVariant = iota
	ArrayOperatorVariantAny
	ArrayOperatorVariantNone
)

func ToArrayOperatorVariant(name string) (ArrayOperatorVariant, error) {
	switch strings.ToUpper(name) {
	case "ALL":
		return ArrayOperatorVariantAll, nil
	case "ANY":
		return ArrayOperatorVariantAny, nil
	case "NONE":
		return ArrayOperatorVariantNone, nil
	default:
		return ArrayOperatorVariant(-1), core.Error(core.ErrInvalidArgument, name)
	}
}

func NewArrayOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	variantStr string,
	comparator core.Predicate,
) (*ArrayOperator, error) {
	if left == nil {
		return nil, core.Error(core.ErrMissedArgument, "left expression")
	}

	if right == nil {
		return nil, core.Error(core.ErrMissedArgument, "right expression")
	}

	variant, err := ToArrayOperatorVariant(variantStr)

	if err != nil {
		return nil, err
	}

	if comparator == nil {
		return nil, core.Error(core.ErrMissedArgument, "comparator expression")
	}

	base := &baseOperator{src, left, right}

	return &ArrayOperator{base, variant, comparator}, nil
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

	switch operator.variant {
	case ArrayOperatorVariantAll:
		return operator.all(ctx, arr, right)
	case ArrayOperatorVariantAny:
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
