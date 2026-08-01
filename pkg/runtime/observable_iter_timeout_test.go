package runtime

import (
	"testing"
	"time"
)

func TestStreamIteratorTimeoutConstructorKeepsMillisecondUnits(t *testing.T) {
	t.Parallel()

	iterator := NewStreamIteratorWithTimeout(nil, time.Duration(2)).(*StreamIterator)
	if iterator.timeout != 2*time.Millisecond {
		t.Fatalf("legacy timeout = %s, want 2ms", iterator.timeout)
	}
}
