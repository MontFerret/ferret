package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ParameterExpression struct {
	src  core.SourceMap
	name string
}

func NewParameterExpression(src core.SourceMap, name string) (*ParameterExpression, error) {
	if name == "" {
		return nil, core.Error(core.ErrMissedArgument, "name")
	}

	return &ParameterExpression{src, name}, nil
}

func (e *ParameterExpression) Exec(ctx context.Context, _ *core.Scope) (core.Value, error) {
	param, err := core.ParamFrom(ctx, e.name)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	return param, nil
}
