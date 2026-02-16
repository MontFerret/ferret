package core

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	CollectorAggregation struct {
		state    vm.Operand
		selector []*AggregateSelector
		fused    bool
	}

	AggregateSelector struct {
		name          runtime.String
		args          int
		funcName      runtime.String
		protectedCall bool
		ctx           antlr.ParserRuleContext
	}
)

func NewCollectorAggregation(state vm.Operand, selector []*AggregateSelector) *CollectorAggregation {
	return &CollectorAggregation{
		state:    state,
		selector: selector,
	}
}

func NewCollectorAggregationFused(state vm.Operand, selector []*AggregateSelector) *CollectorAggregation {
	return &CollectorAggregation{
		state:    state,
		selector: selector,
		fused:    true,
	}
}

func (c *CollectorAggregation) State() vm.Operand {
	return c.state
}

func (c *CollectorAggregation) Selectors() []*AggregateSelector {
	return c.selector
}

func (c *CollectorAggregation) IsFused() bool {
	return c.fused
}

func NewAggregateSelector(name runtime.String, args int, funcName runtime.String, protectedCall bool, ctx antlr.ParserRuleContext) *AggregateSelector {
	return &AggregateSelector{
		name:          name,
		args:          args,
		funcName:      funcName,
		protectedCall: protectedCall,
		ctx:           ctx,
	}
}

func (s *AggregateSelector) Name() runtime.String {
	return s.name
}

func (s *AggregateSelector) Args() int {
	return s.args
}

func (s *AggregateSelector) FuncName() runtime.String {
	return s.funcName
}

func (s *AggregateSelector) ProtectedCall() bool {
	return s.protectedCall
}

func (s *AggregateSelector) Context() antlr.ParserRuleContext {
	return s.ctx
}
