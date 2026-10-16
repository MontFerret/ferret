package encodingownership_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestValueAdopterContextPreservesOperationLifetime(t *testing.T) {
	type key struct{}
	parent, cancel := context.WithCancel(context.WithValue(t.Context(), key{}, "caller"))
	defer cancel()
	var adopted runtime.Value
	ctx := encodingownership.WithValueAdopter(parent, func(value runtime.Value, _ bool) { adopted = value })
	if ctx.Value(key{}) != "caller" || ctx.Done() != parent.Done() || encodingownership.ValueAdopterFromContext(parent) != nil {
		t.Fatal("adoption changed or escaped the operation context")
	}

	cancel()
	encodingownership.ValueAdopterFromContext(ctx)(runtime.Int(1), true)
	if adopted != runtime.Int(1) || ctx.Err() != context.Canceled {
		t.Fatal("cancellation prevented ownership registration")
	}

	if encodingownership.WithValueAdopter(nil, nil) != nil || encodingownership.ValueAdopterFromContext(nil) != nil {
		t.Fatal("nil context was replaced")
	}
}
