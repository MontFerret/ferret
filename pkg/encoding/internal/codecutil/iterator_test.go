package codecutil_test

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding/internal/codecutil"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestIteratorBorrowsComparableValue(t *testing.T) {
	closes := 0
	source := acquiredIterator{closes: &closes}

	handle, err := codecutil.NewIterator(t.Context(), source)
	if err != nil {
		t.Fatal(err)
	}

	if err := handle.Close(); err != nil || closes != 0 {
		t.Fatalf("borrowed value identity: closes=%d err=%v", closes, err)
	}
}

func TestIteratorClosesDistinctNoncomparableValuesOnce(t *testing.T) {
	for _, mode := range []string{"value", "pointer"} {
		t.Run(mode, func(t *testing.T) {
			sourceCloses, iteratorCloses := 0, 0
			iterator := acquiredIterator{marker: []int{1}, closes: &iteratorCloses}
			source := acquiredIterator{iterator: iterator, marker: []int{1}, closes: &sourceCloses, ctx: t.Context()}
			input := runtime.Iterable(source)
			if mode == "pointer" {
				input = &source
			}

			handle, err := codecutil.NewIterator(t.Context(), input)
			if err != nil {
				t.Fatal(err)
			}

			for range 2 {
				if err := handle.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if sourceCloses != 0 || iteratorCloses != 1 {
				t.Fatalf("ownership: source=%d iterator=%d", sourceCloses, iteratorCloses)
			}
		})
	}
}

func TestIteratorBorrowsNoncomparableResourceAliases(t *testing.T) {
	sourceCloses, iteratorCloses := 0, 0
	iterator := identifiedIterator{acquiredIterator: acquiredIterator{marker: []int{1}, closes: &iteratorCloses}, id: 1}
	source := identifiedIterator{acquiredIterator: acquiredIterator{iterator: iterator, marker: []int{2}, closes: &sourceCloses}, id: 1}

	handle, err := codecutil.NewIterator(t.Context(), source)
	if err != nil {
		t.Fatal(err)
	}

	if err := handle.Close(); err != nil || sourceCloses != 0 || iteratorCloses != 0 {
		t.Fatalf("borrowed identity: source=%d iterator=%d err=%v", sourceCloses, iteratorCloses, err)
	}
}

func TestIteratorAcquisitionFailureDoesNotOwnReturnedIterator(t *testing.T) {
	closes := 0
	failure := errors.New("acquisition failed")
	source := acquiredIterator{iterator: acquiredIterator{closes: &closes}, err: failure}

	handle, err := codecutil.NewIterator(t.Context(), source)
	if !errors.Is(err, failure) {
		t.Fatalf("lost acquisition error: %v", err)
	}

	if err := handle.Close(); err != nil || closes != 0 {
		t.Fatalf("failed acquisition took ownership: closes=%d err=%v", closes, err)
	}
}
