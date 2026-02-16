package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type CollectorAggregation struct {
	state    bytecode.Operand
	selector []*AggregateSelector
	fused    bool
}

func NewCollectorAggregation(state bytecode.Operand, selector []*AggregateSelector) *CollectorAggregation {
	return &CollectorAggregation{
		state:    state,
		selector: selector,
	}
}

func NewCollectorAggregationFused(state bytecode.Operand, selector []*AggregateSelector) *CollectorAggregation {
	return &CollectorAggregation{
		state:    state,
		selector: selector,
		fused:    true,
	}
}

func (c *CollectorAggregation) State() bytecode.Operand {
	return c.state
}

func (c *CollectorAggregation) Selectors() []*AggregateSelector {
	return c.selector
}

func (c *CollectorAggregation) IsFused() bool {
	return c.fused
}
