package asm

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDurationConstantFormatting(t *testing.T) {
	t.Parallel()

	if got := constantAsText(runtime.NewDuration(1500 * time.Millisecond)); got != "1.5s" {
		t.Fatalf("duration constant = %q, want 1.5s", got)
	}
}
