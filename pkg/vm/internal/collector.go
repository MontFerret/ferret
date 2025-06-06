package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	CollectorType int

	Collector interface {
		runtime.Value
		runtime.Iterable

		Collect(ctx context.Context, key, value runtime.Value) error
	}

	BaseCollector struct{}
)

const (
	CollectorTypeCounter CollectorType = iota
	CollectorTypeKey
	CollectorTypeKeyCounter
	CollectorTypeKeyGroup
)

func NewCollector(typ CollectorType) Collector {
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

func (*BaseCollector) MarshalJSON() ([]byte, error) {
	panic("not supported")
}

func (*BaseCollector) String() string {
	return "[Collector]"
}

func (*BaseCollector) Unwrap() interface{} {
	panic("not supported")
}

func (*BaseCollector) Hash() uint64 {
	panic("not supported")
}

func (*BaseCollector) Copy() runtime.Value {
	panic("not supported")
}
