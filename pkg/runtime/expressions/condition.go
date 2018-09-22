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
	if core.IsNil(test) {
		return nil, core.Error(core.ErrMissedArgument, "test expression")
	}

	if core.IsNil(alternate) {
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

	cond := e.evalTestValue(out)

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

func (e *ConditionExpression) evalTestValue(value core.Value) values.Boolean {
	switch value.Type() {
	case core.BooleanType:
		return value.(values.Boolean)
	case core.NoneType:
		return values.False
	case core.StringType:
		return values.NewBoolean(value.String() != "")
	case core.IntType:
		return values.NewBoolean(value.(values.Int) != 0)
	case core.FloatType:
		return values.NewBoolean(value.(values.Float) != 0)
	default:
		return values.True
	}
}
