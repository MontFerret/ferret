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
	ctx := encodingownership.WithValueAdopter(parent, func(value runtime.Value) { adopted = value })
	if ctx.Value(key{}) != "caller" || ctx.Done() != parent.Done() || encodingownership.ValueAdopterFromContext(parent) != nil {
		t.Fatal("adoption changed or escaped the operation context")
	}

	cancel()
	encodingownership.ValueAdopterFromContext(ctx)(runtime.Int(1))
	if adopted != runtime.Int(1) || ctx.Err() != context.Canceled {
		t.Fatal("cancellation prevented ownership registration")
	}

	if encodingownership.WithValueAdopter(nil, nil) != nil || encodingownership.ValueAdopterFromContext(nil) != nil {
		t.Fatal("nil context was replaced")
	}

	operation := encodingownership.NewOperationContext(ctx)
	encodingownership.OperationValueAdopter(operation)(runtime.Int(2))
	if adopted != runtime.Int(2) || operation.Done() != ctx.Done() || operation.Err() != context.Canceled || operation.Value(key{}) != "caller" {
		t.Fatal("cached adoption changed the context or stopped on cancellation")
	}

	if encodingownership.NewOperationContext(parent) != parent || encodingownership.NewOperationContext(nil) != nil || encodingownership.OperationValueAdopter(parent) != nil || encodingownership.OperationValueAdopter(nil) != nil {
		t.Fatal("operation without adoption acquired a cache")
	}
}
