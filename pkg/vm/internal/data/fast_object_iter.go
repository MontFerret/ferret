package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	fastObjectEntry struct {
		key  string
		slot int
	}

	fastObjectIterator struct {
		entries []fastObjectEntry
		slots   []runtime.Value
		pos     int
	}

	fastObjectDictIterator struct {
		dict map[string]runtime.Value
		keys []string
		pos  int
	}
)

func (iter *fastObjectIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if iter.pos >= len(iter.entries) {
		return runtime.None, runtime.None, io.EOF
	}

	entry := iter.entries[iter.pos]
	value := iter.slots[entry.slot]
	if value == nil {
		value = runtime.None
	}
	iter.pos++

	return value, runtime.String(entry.key), nil
}

func (iter *fastObjectDictIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if iter.pos >= len(iter.keys) {
		return runtime.None, runtime.None, io.EOF
	}

	key := iter.keys[iter.pos]
	value := iter.dict[key]
	iter.pos++

	return value, runtime.String(key), nil
}
