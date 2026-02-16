package core

import (
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type CollectorAggregation struct {
	state    vm.Operand
	selector []*AggregateSelector
	fused    bool
}

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
