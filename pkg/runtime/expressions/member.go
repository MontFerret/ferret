package expressions

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MemberExpression struct {
	src    core.SourceMap
	source core.Expression
	path   []*MemberPathSegment
}

func NewMemberExpression(src core.SourceMap, source core.Expression, path []*MemberPathSegment) (*MemberExpression, error) {
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
		if e.path[0].optional {
			return values.None, nil
		}

		return values.None, core.SourceError(
			e.src,
			err,
		)
	}

	out := val
	path := make([]core.Value, 1)

	for _, seg := range e.path {
		segment, err := seg.exp.Exec(ctx, scope)

		if err != nil {
			return values.None, err
		}

		path[0] = segment
		c, err := values.GetIn(ctx, out, path)

		if err != nil {
			if !seg.optional {
				return values.None, core.SourceError(e.src, err)
			}

			return values.None, nil
		}

		out = c
	}

	return out, nil
}
