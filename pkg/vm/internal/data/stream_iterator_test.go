package data

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestStreamIteratorDurationTimeoutPrecision(t *testing.T) {
	t.Parallel()

	iterator := newStreamIterator(nil, runtime.Duration(2)).(*streamIterator)
	if iterator.timeout != 2*time.Nanosecond {
		t.Fatalf("native timeout = %s, want 2ns", iterator.timeout)
	}
}
