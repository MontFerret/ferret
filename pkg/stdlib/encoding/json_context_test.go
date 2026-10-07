package encoding

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestJSONFunctionsPreserveCancellationAndStrictDocuments(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, call := range []func(context.Context, runtime.Value) (runtime.Value, error){JSONParse, JSONStringify} {
		if _, err := call(ctx, runtime.NewString("1")); !errors.Is(err, context.Canceled) {
			t.Fatalf("lost cancellation: %v", err)
		}
	}
	for _, text := range []string{"[1,]", "[1 2]", "[1", "1 {}"} {
		if _, err := JSONParse(t.Context(), runtime.NewString(text)); err == nil {
			t.Errorf("accepted %q", text)
		}
	}
}
