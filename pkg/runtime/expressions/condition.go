package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ConditionExpression struct {
	src        core.SourceMap
	test       core.Expression
	consequent core.Expression
	alternate  core.Expression
}

func NewConditionExpression(
	src core.SourceMap,
	test core.Expression,
	consequent core.Expression,
	alternate core.Expression,
) (*ConditionExpression, error) {
	if test == nil {
		return nil, core.Error(core.ErrMissedArgument, "test expression")
	}

	if alternate == nil {
		return nil, core.Error(core.ErrMissedArgument, "alternate expression")
	}

	return &ConditionExpression{
		src,
		test,
		consequent,
		alternate,
	}, nil
}

func (e *ConditionExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	out, err := e.test.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	cond := values.ToBoolean(out)

	var next core.Expression

	if cond == values.True {
		next = e.consequent

		// shortcut version
		if next == nil {
			return out, nil
		}
	} else {
		next = e.alternate
	}

	res, err := next.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	return res, nil
}
