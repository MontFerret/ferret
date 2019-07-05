package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"regexp"
)

type (
	RegexpOperatorType int
	RegexpOperator     struct {
		*baseOperator
		opType RegexpOperatorType
	}
)

const (
	RegexpOperatorTypeNegative RegexpOperatorType = 0
	RegexpOperatorTypePositive RegexpOperatorType = 1
)

var regexpOperators = map[string]RegexpOperatorType{
	"!~": RegexpOperatorTypeNegative,
	"=~": RegexpOperatorTypePositive,
}

func NewRegexpOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	operator string,
) (*RegexpOperator, error) {
	op, exists := regexpOperators[operator]

	if !exists {
		return nil, core.Error(core.ErrInvalidArgument, "operator")
	}

	return &RegexpOperator{
		&baseOperator{
			src,
			left,
			right,
		},
		op,
	}, nil
}

func (operator *RegexpOperator) Type() RegexpOperatorType {
	return operator.opType
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

	if operator.opType == RegexpOperatorTypePositive {
		return values.NewBoolean(r.MatchString(leftStr)), nil
	}

	return values.NewBoolean(!r.MatchString(leftStr)), nil
}
