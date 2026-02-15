package data

type CollectorType int

const (
	CollectorTypeCounter CollectorType = iota
	CollectorTypeKey
	CollectorTypeKeyCounter
	CollectorTypeKeyGroup
	CollectorTypeAggregate
	CollectorTypeAggregateGroup
)

func NewCollector(typ CollectorType) Transformer {
	switch typ {
	case CollectorTypeCounter:
		return NewCounterCollector()
	case CollectorTypeKey:
		return NewKeyCollector()
	case CollectorTypeKeyCounter:
		return NewKeyCounterCollector()
	case CollectorTypeKeyGroup:
		return NewKeyGroupCollector()
	case CollectorTypeAggregate:
		panic("aggregate collector requires a plan")
	case CollectorTypeAggregateGroup:
		panic("aggregate group collector requires a plan")
	default:
		panic("unknown collector type")
	}
}
