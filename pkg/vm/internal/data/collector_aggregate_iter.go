package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type aggregateIterator struct {
	collector        *AggregateCollector
	groupKeys        []string
	aggregateIdx     int
	groupIdx         int
	emittedSingleKey bool
}

func (it *aggregateIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it.collector == nil || !it.collector.hasData {
		return runtime.None, runtime.None, io.EOF
	}

	if it.aggregateIdx < len(it.collector.plan.Keys) {
		idx := it.aggregateIdx
		it.aggregateIdx++

		return it.collector.valueFor(idx), it.collector.plan.Keys[idx], nil
	}

	if it.collector.hasSingleGroup && !it.emittedSingleKey {
		it.emittedSingleKey = true

		return it.collector.singleGroupValue, runtime.NewString(it.collector.singleGroupKey), nil
	}

	if it.groupIdx >= len(it.groupKeys) {
		return runtime.None, runtime.None, io.EOF
	}

	key := it.groupKeys[it.groupIdx]
	it.groupIdx++

	return it.collector.groups[key], runtime.NewString(key), nil
}
