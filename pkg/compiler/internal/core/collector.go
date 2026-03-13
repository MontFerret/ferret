package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type Collector struct {
	projection     *CollectorProjection
	aggregation    *CollectorAggregation
	groupSelectors []*CollectSelector
	typ            bytecode.CollectorType
	dst            bytecode.Operand
}

func NewCollector(type_ bytecode.CollectorType, dst bytecode.Operand, projection *CollectorProjection, groupSelectors []*CollectSelector, aggregation *CollectorAggregation) *Collector {
	return &Collector{
		typ:            type_,
		dst:            dst,
		projection:     projection,
		groupSelectors: groupSelectors,
		aggregation:    aggregation,
	}
}

func DetermineCollectorType(withGrouping, withAggregation, withProjection, withCounter bool) bytecode.CollectorType {
	if withGrouping {
		if withCounter {
			return bytecode.CollectorTypeKeyCounter
		}

		return bytecode.CollectorTypeKeyGroup
	}

	if withAggregation {
		return bytecode.CollectorTypeKeyGroup
	}

	return bytecode.CollectorTypeCounter
}

func (c *Collector) Type() bytecode.CollectorType {
	return c.typ
}

func (c *Collector) Destination() bytecode.Operand {
	return c.dst
}

func (c *Collector) Projection() *CollectorProjection {
	return c.projection
}

func (c *Collector) GroupSelectors() []*CollectSelector {
	return c.groupSelectors
}

func (c *Collector) Aggregation() *CollectorAggregation {
	return c.aggregation
}

func (c *Collector) HasProjection() bool {
	return c.projection != nil
}

func (c *Collector) HasGrouping() bool {
	return len(c.groupSelectors) > 0
}

func (c *Collector) HasAggregation() bool {
	return c.aggregation != nil
}
