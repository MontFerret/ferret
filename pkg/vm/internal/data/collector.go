package data

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

func NewCollector(typ bytecode.CollectorType) Transformer {
	switch typ {
	case bytecode.CollectorTypeCounter:
		return NewCounterCollector()
	case bytecode.CollectorTypeKey:
		return NewKeyCollector()
	case bytecode.CollectorTypeKeyCounter:
		return NewKeyCounterCollector()
	case bytecode.CollectorTypeKeyGroup:
		return NewKeyGroupCollector()
	case bytecode.CollectorTypeAggregate:
		panic("aggregate collector requires a plan")
	case bytecode.CollectorTypeAggregateGroup:
		panic("aggregate group collector requires a plan")
	default:
		panic("unknown collector type")
	}
}
