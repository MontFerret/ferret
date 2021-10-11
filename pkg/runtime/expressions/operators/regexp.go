package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"regexp"
)

type (
	RegexpOperatorVariant int

	RegexpOperator struct {
		*baseOperator
		variant RegexpOperatorVariant
	}
)

const (
	RegexpOperatorVariantNegative RegexpOperatorVariant = 0
	RegexpOperatorVariantPositive RegexpOperatorVariant = 1
)

var regexpVariants = map[string]RegexpOperatorVariant{
	"!~": RegexpOperatorVariantNegative,
	"=~": RegexpOperatorVariantPositive,
}

func NewRegexpOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	operatorStr string,
) (*RegexpOperator, error) {
	variant, exists := regexpVariants[operatorStr]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator")
	}

	return &RegexpOperator{
		&baseOperator{
			src,
			left,
			right,
		},
		variant,
	}, nil
}

func (operator *RegexpOperator) Type() RegexpOperatorVariant {
	return operator.variant
}

func (operator *RegexpOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return operator.Eval(ctx, left, right)
}

func (operator *RegexpOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
	leftStr := left.String()
	rightStr := right.String()

	r, err := regexp.Compile(rightStr)

	if err != nil {
		return values.None, err
	}

	if operator.variant == RegexpOperatorVariantPositive {
		return values.NewBoolean(r.MatchString(leftStr)), nil
	}

	return values.NewBoolean(!r.MatchString(leftStr)), nil
}
