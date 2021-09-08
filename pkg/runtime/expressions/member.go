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
	member, err := e.source.Exec(ctx, scope)

	if err != nil {
		if e.path[0].optional {
			return values.None, nil
		}

		return values.None, core.SourceError(
			e.src,
			err,
		)
	}

	// keep information about all optional path segments
	optionals := make(map[int]bool)
	segments := make([]core.Value, len(e.path))

	// unfold the path
	for i, seg := range e.path {
		segment, err := seg.exp.Exec(ctx, scope)

		if err != nil {
			return values.None, err
		}

		segments[i] = segment

		if seg.optional {
			optionals[i] = true
		}
	}

	var pathErr core.PathError
	var out core.Value = values.None

	getter, ok := member.(core.Getter)

	if ok {
		out, pathErr = getter.GetIn(ctx, segments)
	} else {
		out, pathErr = values.GetIn(ctx, member, segments)
	}

	if pathErr != nil {
		_, isOptional := optionals[int(pathErr.Segment())]

		// we either cannot determine what segment caused the issue
		// or it is not optional, thus we just return the error
		if !isOptional {
			return values.None, pathErr
		}

		return values.None, nil
	}

	return out, nil
}
