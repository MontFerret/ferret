package core

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	CollectorAggregation struct {
		state    vm.Operand
		selector []*AggregateSelector
	}

	AggregateSelector struct {
		name          runtime.String
		args          int
		funcName      runtime.String
		protectedCall bool
	}
)

func NewCollectorAggregation(state vm.Operand, selector []*AggregateSelector) *CollectorAggregation {
	return &CollectorAggregation{
		state:    state,
		selector: selector,
	}
}

func (c *CollectorAggregation) State() vm.Operand {
	return c.state
}

func (c *CollectorAggregation) Selectors() []*AggregateSelector {
	return c.selector
}

func NewAggregateSelector(name runtime.String, args int, funcName runtime.String, protectedCall bool) *AggregateSelector {
	return &AggregateSelector{
		name:          name,
		args:          args,
		funcName:      funcName,
		protectedCall: protectedCall,
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
