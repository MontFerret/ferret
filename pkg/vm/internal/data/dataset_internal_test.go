package data

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDataSetBufferedLengthSkipsUnknownBacking(t *testing.T) {
	// Every list method panics if called on this nil embedded interface.
	type unreadableList struct{ runtime.List }
	for _, ds := range []*DataSet{nil, {}, {values: unreadableList{}}} {
		if _, _, err := ds.BufferedAt(t.Context(), 0); !errors.Is(err, errors.ErrUnsupported) {
			t.Fatalf("unknown buffer read: %v", err)
		}
		if got := ds.BufferedLength(); got != -1 {
			t.Fatalf("missing native buffer has length %d", got)
		}
	}
}
