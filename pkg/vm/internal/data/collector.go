package data

type CollectorType int

const (
	CollectorTypeCounter CollectorType = iota
	CollectorTypeKey
	CollectorTypeKeyCounter
	CollectorTypeKeyGroup
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
	default:
		panic("unknown collector type")
	}
}
