package data

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func NewCollector(typ bytecode.CollectorType) Transformer {
	collector, err := NewCollectorSafe(typ)
	if err != nil {
		return NewNoopCollector()
	}

	return collector
}

func NewCollectorSafe(typ bytecode.CollectorType) (Transformer, error) {
	switch typ {
	case bytecode.CollectorTypeCounter:
		return NewCounterCollector(), nil
	case bytecode.CollectorTypeKey:
		return NewKeyCollector(), nil
	case bytecode.CollectorTypeKeyCounter:
		return NewKeyCounterCollector(), nil
	case bytecode.CollectorTypeKeyGroup:
		return NewKeyGroupCollector(), nil
	case bytecode.CollectorTypeAggregate:
		return nil, fmt.Errorf("collector type %d requires aggregate plan", typ)
	case bytecode.CollectorTypeAggregateGroup:
		return nil, fmt.Errorf("collector type %d requires aggregate plan", typ)
	default:
		return nil, fmt.Errorf("unknown collector type %d", typ)
	}
}
