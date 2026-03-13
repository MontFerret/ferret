package core

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type CollectSelector struct {
	ctx  antlr.ParserRuleContext
	name runtime.String
}

func NewCollectSelector(name runtime.String, ctx antlr.ParserRuleContext) *CollectSelector {
	return &CollectSelector{
		name: name,
		ctx:  ctx,
	}
}

func (s *CollectSelector) Name() runtime.String {
	return s.name
}

func (s *CollectSelector) Context() antlr.ParserRuleContext {
	return s.ctx
}
