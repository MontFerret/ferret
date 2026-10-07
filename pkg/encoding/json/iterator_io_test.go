package json_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type closingIterable struct {
	iter *closingIterator
	iterOnly
}
type closingIterator struct {
	closeErr error
	iterOnlyIter
	closes int
}

func (v *closingIterable) Iterate(context.Context) (runtime.Iterator, error) { return v.iter, nil }
func (v *closingIterator) Close() error                                      { v.closes++; return v.closeErr }

type iteratorWriter struct {
	iter    *closingIterator
	failure error
	writes  int
}

func (w *iteratorWriter) Write(data []byte) (int, error) {
	w.writes++
	if w.iter.closes != 0 {
		return 0, errors.New("wrote after iterator closure")
	}
	// '[' succeeds; the first item's write fails before the second Next call.
	if w.writes == 2 {
		return 0, w.failure
	}
	return len(data), nil
}

func TestJSONWritesIterableBeforeClosingAndSettlesFailures(t *testing.T) {
	writeErr, closeErr := errors.New("write failed"), errors.New("close failed")
	iter := &closingIterator{iterOnlyIter: iterOnlyIter{items: []runtime.Value{runtime.Int(1), runtime.Int(2)}}, closeErr: closeErr}
	value := &closingIterable{iter: iter}
	writer := &iteratorWriter{iter: iter, failure: writeErr}
	err := json.Default.Encode(t.Context(), writer, value)
	if !errors.Is(err, writeErr) || !errors.Is(err, closeErr) || iter.closes != 1 || iter.idx != 1 || writer.writes != 2 {
		t.Fatalf("iterator settlement: next=%d closes=%d writes=%d err=%v", iter.idx, iter.closes, writer.writes, err)
	}
	iter = &closingIterator{iterOnlyIter: iterOnlyIter{items: []runtime.Value{runtime.Int(1)}}}
	if err := json.Default.Encode(t.Context(), io.Discard, &closingIterable{iter: iter}); err != nil || iter.closes != 1 {
		t.Fatalf("successful iterator settlement: closes=%d err=%v", iter.closes, err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	items := make([]runtime.Value, 1000)
	for i := range items {
		items[i] = runtime.Int(1)
	}
	iter = &closingIterator{iterOnlyIter: iterOnlyIter{items: items}}
	enc := json.Default.EncodeWith().PreHook(func(context.Context, runtime.Value) error { cancel(); return nil }).Encoder()
	err = enc.Encode(ctx, io.Discard, &closingIterable{iter: iter})
	if !errors.Is(err, context.Canceled) || iter.closes != 1 || iter.idx > 256 {
		t.Fatalf("canceled iterator settlement: next=%d closes=%d err=%v", iter.idx, iter.closes, err)
	}
}
