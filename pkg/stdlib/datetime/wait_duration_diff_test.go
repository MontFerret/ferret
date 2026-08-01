package datetime

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestWaitDurationDiff(t *testing.T) {
	t.Parallel()

	start := runtime.NewDateTime(time.Unix(0, 1))
	end := runtime.NewDateTime(time.Unix(0, 4))

	actual, err := waitDurationDiff(t.Context(), start, end)
	if err != nil {
		t.Fatal(err)
	}
	if actual != runtime.Duration(3) {
		t.Fatalf("duration difference = %v, want 3ns", actual)
	}
}
