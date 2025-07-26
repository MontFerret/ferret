package core

import "github.com/MontFerret/ferret/pkg/vm"

type (
	CollectorType int

	Collector struct {
		typ            CollectorType
		dst            vm.Operand
		projection     *CollectorProjection
		groupSelectors []*CollectSelector
		aggregation    *CollectorAggregation
	}
)

const (
	CollectorTypeCounter CollectorType = iota
	CollectorTypeKey
	CollectorTypeKeyCounter
	CollectorTypeKeyGroup
)

func NewCollector(type_ CollectorType, dst vm.Operand, projection *CollectorProjection, groupSelectors []*CollectSelector, aggregation *CollectorAggregation) *Collector {
	return &Collector{
		typ:            type_,
		dst:            dst,
		projection:     projection,
		groupSelectors: groupSelectors,
		aggregation:    aggregation,
	}
}

func DetermineCollectorType(withGrouping, withAggregation, withProjection, withCounter bool) CollectorType {
	if withGrouping {
		if withCounter {
			return CollectorTypeKeyCounter
		}

		return CollectorTypeKeyGroup
	}

	if withAggregation {
		return CollectorTypeKeyGroup
	}

	return CollectorTypeCounter
}

func (c *Collector) Type() CollectorType {
	return c.typ
}

func (c *Collector) Destination() vm.Operand {
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
