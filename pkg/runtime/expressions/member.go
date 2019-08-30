package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MemberExpression struct {
	src    core.SourceMap
	source core.Expression
	path   []core.Expression
}

func NewMemberExpression(src core.SourceMap, source core.Expression, path []core.Expression) (*MemberExpression, error) {
	if source == nil {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	if len(path) == 0 {
		return nil, core.Error(core.ErrMissedArgument, "path expressions")
	}

	return &MemberExpression{src, source, path}, nil
}

func (e *MemberExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	val, err := e.source.Exec(ctx, scope)

	if err != nil {
		return values.None, core.SourceError(
			e.src,
			err,
		)
	}

	strPath := make([]core.Value, len(e.path))

	for idx, exp := range e.path {
		segment, err := exp.Exec(ctx, scope)

		if err != nil {
			return values.None, err
		}

		strPath[idx] = segment
	}

	out, err := values.GetIn(ctx, val, strPath)

	if err != nil {
		return values.None, core.SourceError(e.src, err)
	}

	return out, nil
}
