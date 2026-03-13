package core

import "github.com/antlr4-go/antlr/v4"

type CollectorProjection struct {
	ctx            antlr.ParserRuleContext
	groupsVariable string
	countVariable  string
}

func NewCollectorGroupProjection(groupsVariable string, ctx antlr.ParserRuleContext) *CollectorProjection {
	return &CollectorProjection{
		groupsVariable: groupsVariable,
		countVariable:  "",
		ctx:            ctx,
	}
}

func NewCollectorCountProjection(countVariable string, ctx antlr.ParserRuleContext) *CollectorProjection {
	return &CollectorProjection{
		groupsVariable: "",
		countVariable:  countVariable,
		ctx:            ctx,
	}
}

func (p *CollectorProjection) VariableName() string {
	if p.groupsVariable != "" {
		return p.groupsVariable
	}

	if p.countVariable != "" {
		return p.countVariable
	}

	return ""
}

func (p *CollectorProjection) IsGrouped() bool {
	return p.groupsVariable != ""
}

func (p *CollectorProjection) IsCounted() bool {
	return p.countVariable != ""
}

func (p *CollectorProjection) Context() antlr.ParserRuleContext {
	return p.ctx
}
