package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	assignmentTarget struct {
		Root     string
		RootCtx  antlr.ParserRuleContext
		RootTok  antlr.Token
		Segments []assignmentTargetSegment
	}

	assignmentTargetSegment struct {
		Context  *fql.AssignmentTargetPathContext
		Property fql.IPropertyNameContext
		Computed fql.IComputedPropertyNameContext
		Safe     bool
	}
)

func newAssignmentTarget(ctx fql.IAssignmentTargetContext) (assignmentTarget, bool) {
	targetCtx, ok := ctx.(*fql.AssignmentTargetContext)
	if !ok || targetCtx == nil {
		return assignmentTarget{}, false
	}

	root := targetCtx.BindingIdentifier()
	if root == nil {
		return assignmentTarget{}, false
	}

	out := assignmentTarget{
		Root:    textOfBindingIdentifier(root),
		RootCtx: root.(antlr.ParserRuleContext),
		RootTok: root.GetStart(),
	}

	for _, path := range targetCtx.AllAssignmentTargetPath() {
		segmentCtx, ok := path.(*fql.AssignmentTargetPathContext)
		if !ok || segmentCtx == nil {
			return assignmentTarget{}, false
		}

		segment := assignmentTargetSegment{
			Context:  segmentCtx,
			Property: segmentCtx.PropertyName(),
			Computed: segmentCtx.ComputedPropertyName(),
			Safe:     segmentCtx.ErrorOperator() != nil,
		}

		if segment.Property == nil && segment.Computed == nil {
			return assignmentTarget{}, false
		}

		out.Segments = append(out.Segments, segment)
	}

	return out, true
}
