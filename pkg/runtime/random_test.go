package runtime_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDeprecatedStandaloneRandomHelpers(t *testing.T) {
	if got := runtime.RandomDefault(); got < 0 || got >= 1 {
		t.Fatalf("standalone unit draw = %v", got)
	}

	if got := runtime.Random(5, 5); got != 5 {
		t.Fatalf("standalone equal bounds = %v", got)
	}

	if got := runtime.Random2(0); got != 0 {
		t.Fatalf("standalone zero midpoint = %v", got)
	}
}
