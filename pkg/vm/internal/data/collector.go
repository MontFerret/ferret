package data

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

type noopCollector struct{}

func NewNoopCollector() Transformer {
	return &noopCollector{}
}

func (c *noopCollector) String() string {
	return "[NoopCollector]"
}

func (c *noopCollector) Hash() uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte("vm.noop_collector"))

	return hasher.Sum64()
}

func (c *noopCollector) Copy() runtime.Value {
	return &noopCollector{}
}

func (c *noopCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return runtime.NewArray(0).Iterate(ctx)
}

func (c *noopCollector) Length(_ context.Context) (runtime.Int, error) {
	return 0, nil
}

func (c *noopCollector) Get(_ context.Context, key runtime.Value) (runtime.Value, error) {
	return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %v", key)
}

func (c *noopCollector) Set(_ context.Context, _, _ runtime.Value) error {
	return nil
}

func (c *noopCollector) Close() error {
	return nil
}
