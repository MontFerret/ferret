package core

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	baseAggregateSelector struct {
		name          runtime.String
		funcName      runtime.String
		protectedCall bool
	}

	// AggregateSelector stores the finalized aggregate metadata that is needed
	// during the post-collection evaluation stage (arguments are already materialized in the collector).
	AggregateSelector struct {
		ctx antlr.ParserRuleContext
		baseAggregateSelector
		args int
	}

	// CompiledAggregateSelector stores selector metadata together with compiled argument registers.
	// It is a transient representation used while building collector input.
	CompiledAggregateSelector struct {
		ctx fql.ICollectAggregateSelectorContext
		baseAggregateSelector
		args RegisterSequence
	}
)

func NewAggregateSelector(name runtime.String, args int, funcName runtime.String, protectedCall bool, ctx antlr.ParserRuleContext) *AggregateSelector {
	return &AggregateSelector{
		baseAggregateSelector: baseAggregateSelector{
			name:          name,
			funcName:      funcName,
			protectedCall: protectedCall,
		},
		args: args,
		ctx:  ctx,
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

func NewCompiledAggregateSelector(name runtime.String, args RegisterSequence, funcName runtime.String, protectedCall bool, ctx fql.ICollectAggregateSelectorContext) *CompiledAggregateSelector {
	return &CompiledAggregateSelector{
		baseAggregateSelector: baseAggregateSelector{
			name:          name,
			funcName:      funcName,
			protectedCall: protectedCall,
		},
		args: args,
		ctx:  ctx,
	}
}

func (s *CompiledAggregateSelector) Name() runtime.String {
	return s.name
}

func (s *CompiledAggregateSelector) Args() RegisterSequence {
	return s.args
}

func (s *CompiledAggregateSelector) FuncName() runtime.String {
	return s.funcName
}

func (s *CompiledAggregateSelector) ProtectedCall() bool {
	return s.protectedCall
}

func (s *CompiledAggregateSelector) Context() fql.ICollectAggregateSelectorContext {
	return s.ctx
}
