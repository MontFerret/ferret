package data

import "github.com/MontFerret/ferret/pkg/runtime"

type CollectorType int

const (
	CollectorTypeCounter CollectorType = iota
	CollectorTypeKey
	CollectorTypeKeyCounter
	CollectorTypeKeyGroup
)

func NewCollector(typ CollectorType, alloc runtime.Allocator) Transformer {
	switch typ {
	case CollectorTypeCounter:
		return NewCounterCollector(alloc)
	case CollectorTypeKey:
		return NewKeyCollector(alloc)
	case CollectorTypeKeyCounter:
		return NewKeyCounterCollector(alloc)
	case CollectorTypeKeyGroup:
		return NewKeyGroupCollector(alloc)
	default:
		panic("unknown collector type")
	}
}
