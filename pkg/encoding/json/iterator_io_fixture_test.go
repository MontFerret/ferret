package json_test

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	closingIterable struct {
		iter *closingIterator
		iterOnly
	}

	closingIterator struct {
		closeErr error
		iterOnlyIter
		closes int
	}

	iteratorWriter struct {
		iter    *closingIterator
		failure error
		writes  int
	}
)

func (v *closingIterable) Iterate(context.Context) (runtime.Iterator, error) {
	return v.iter, nil
}

func (v *closingIterator) Close() error {
	v.closes++

	return v.closeErr
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
