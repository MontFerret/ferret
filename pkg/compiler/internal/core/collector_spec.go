package core

import "github.com/MontFerret/ferret/pkg/vm"

type (
	CollectorType int

	CollectorSpec struct {
		typ                  CollectorType
		dst                  vm.Operand
		projection           *CollectorProjection
		groupSelectors       []*CollectSelector
		aggregationSelectors []*AggregateSelector
	}
)

const (
	CollectorTypeCounter CollectorType = iota
	CollectorTypeKey
	CollectorTypeKeyCounter
	CollectorTypeKeyGroup
)

func NewCollectorSpec(type_ CollectorType, dst vm.Operand, projection *CollectorProjection, groupSelectors []*CollectSelector, aggregationSelectors []*AggregateSelector) *CollectorSpec {
	return &CollectorSpec{
		typ:                  type_,
		dst:                  dst,
		projection:           projection,
		groupSelectors:       groupSelectors,
		aggregationSelectors: aggregationSelectors,
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

func (c *CollectorSpec) Type() CollectorType {
	return c.typ
}

func (c *CollectorSpec) Destination() vm.Operand {
	return c.dst
}

func (c *CollectorSpec) Projection() *CollectorProjection {
	return c.projection
}

func (c *CollectorSpec) GroupSelectors() []*CollectSelector {
	return c.groupSelectors
}

func (c *CollectorSpec) AggregationSelectors() []*AggregateSelector {
	return c.aggregationSelectors
}

func (c *CollectorSpec) HasProjection() bool {
	return c.projection != nil
}

func (c *CollectorSpec) HasGrouping() bool {
	return len(c.groupSelectors) > 0
}

func (c *CollectorSpec) HasAggregation() bool {
	return len(c.aggregationSelectors) > 0
}
