package runtime

import (
	"testing"
	"time"
)

func TestNewStreamIteratorWithTimeoutPreservesExactDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{name: "nanoseconds", timeout: 2 * time.Nanosecond},
		{name: "milliseconds", timeout: time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			iterator := NewStreamIteratorWithTimeout(nil, tt.timeout).(*StreamIterator)
			if iterator.timeout != tt.timeout {
				t.Fatalf("timeout = %s, want %s", iterator.timeout, tt.timeout)
			}
		})
	}
}

func TestNewStreamIteratorUsesDefaultTimeout(t *testing.T) {
	t.Parallel()

	iterator := NewStreamIterator(nil).(*StreamIterator)
	if iterator.timeout != 5*time.Second {
		t.Fatalf("timeout = %s, want 5s", iterator.timeout)
	}
}
