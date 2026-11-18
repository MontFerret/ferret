package data

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDataSetBufferedLengthSkipsUnknownBacking(t *testing.T) {
	// Every list method panics if called on this nil embedded interface.
	type unreadableList struct{ runtime.List }
	for _, ds := range []*DataSet{nil, {}, {values: unreadableList{}}} {
		if got := ds.BufferedLength(); got != -1 {
			t.Fatalf("missing native buffer has length %d", got)
		}
	}
}
