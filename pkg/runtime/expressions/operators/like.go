package operators

import (
	"context"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type LikeOperator struct {
	*baseOperator
	negate bool
}

func NewLikeOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	negate bool,
) (*LikeOperator, error) {
	if left == nil {
		return nil, core.Error(core.ErrMissedArgument, "left expression")
	}

	if right == nil {
		return nil, core.Error(core.ErrMissedArgument, "right expression")
	}

	return &LikeOperator{&baseOperator{src, left, right}, negate}, nil
}

func (operator *LikeOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
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

func (operator *LikeOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
	err := core.ValidateType(right, types.String)

	if err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	err = core.ValidateType(left, types.String)

	if err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	r, err := glob.Compile(right.String())

	if err != nil {
		return nil, errors.Wrap(err, "invalid glob pattern")
	}

	result := r.Match(left.String())

	if operator.negate {
		return values.NewBoolean(!result), nil
	}

	return values.NewBoolean(result), nil
}
