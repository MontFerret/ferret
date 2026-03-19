package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type kvEntriesIterator struct {
	entries []*KV
	idx     int
}

func newKVEntriesIterator(entries []*KV) runtime.Iterator {
	return &kvEntriesIterator{entries: entries}
}

func (it *kvEntriesIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it == nil || it.idx >= len(it.entries) {
		return runtime.None, runtime.None, io.EOF
	}

	entry := it.entries[it.idx]
	it.idx++

	return entry.Value, entry.Key, nil
}
