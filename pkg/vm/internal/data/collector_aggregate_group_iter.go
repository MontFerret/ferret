package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type groupedAggregateIterator struct {
	entries          []*groupedAggregateEntry
	idx              int
	trackGroupValues bool
}

func newGroupedAggregateIterator(entries []*groupedAggregateEntry, trackGroupValues bool) runtime.Iterator {
	return &groupedAggregateIterator{
		entries:          entries,
		trackGroupValues: trackGroupValues,
	}
}

func (it *groupedAggregateIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it == nil || it.idx >= len(it.entries) {
		return runtime.None, runtime.None, io.EOF
	}

	entry := it.entries[it.idx]
	it.idx++

	if !it.trackGroupValues || entry.group == nil {
		return runtime.None, entry.key, nil
	}

	return entry.group, entry.key, nil
}
