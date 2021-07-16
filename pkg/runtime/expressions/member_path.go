package expressions

import "github.com/MontFerret/ferret/pkg/runtime/core"

type MemberPathSegment struct {
	exp      core.Expression
	optional bool
}

func NewMemberPathSegment(source core.Expression, optional bool) (*MemberPathSegment, error) {
	if source == nil {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	return &MemberPathSegment{source, optional}, nil
}
