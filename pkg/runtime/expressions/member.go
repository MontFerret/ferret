package expressions

import (
	"context"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MemberExpression struct {
	src             core.SourceMap
	source          core.Expression
	path            []*MemberPathSegment
	preCompiledPath []core.Value
}

func NewMemberExpression(src core.SourceMap, source core.Expression, path []*MemberPathSegment, preCompiledPath []core.Value) (*MemberExpression, error) {
	if source == nil {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	if len(path) == 0 {
		return nil, core.Error(core.ErrMissedArgument, "path expressions")
	}

	return &MemberExpression{src, source, path, preCompiledPath}, nil
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

	var segments = e.preCompiledPath

	if e.preCompiledPath == nil {
		segments = make([]core.Value, len(e.path))

		// unfold the path
		for i, seg := range e.path {
			segment, err := seg.exp.Exec(ctx, scope)

			if err != nil {
				return values.None, err
			}

			segments[i] = segment
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
		segmentIdx := pathErr.Segment()
		// if invalid index is returned, we ignore the optionality check
		// and return the pathErr
		if segmentIdx >= len(e.path) {
			return values.None, errors.New(pathErr.Format(segments))
		}

		segment := e.path[segmentIdx]

		if !segment.optional {
			return values.None, errors.New(pathErr.Format(segments))
		}

		return values.None, nil
	}

	return out, nil
}
