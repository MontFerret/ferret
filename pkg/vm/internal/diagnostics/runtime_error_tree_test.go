package diagnostics

import (
	"errors"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestRuntimeAggregateUnwrapPreservesEveryCause(t *testing.T) {
	first, second := errors.New("first"), errors.New("second")
	set := NewRuntimeErrorSet(2)
	set.Add(&RuntimeError{Diagnostic: diagnostics.NewUnexpectedErrorWith(source.NewAnonymous("RETURN 1"), "first", first)})
	set.Add(&RuntimeError{Diagnostic: diagnostics.NewUnexpectedErrorWith(source.NewAnonymous("RETURN 2"), "second", second)})
	wrapped := fmt.Errorf("aggregate: %w", set)

	result := WrapRuntimeError(nil, 0, nil, wrapped)
	if result != wrapped || !errors.Is(result, first) || !errors.Is(result, second) {
		t.Fatalf("aggregate lost causes: %v", result)
	}

	var diagnostic *diagnostics.Diagnostic
	if !errors.As(result, &diagnostic) || diagnostic != set.First().Diagnostic {
		t.Fatal("native diagnostic is not traversable")
	}

	children := set.Unwrap()
	children[0] = nil

	if !errors.Is(set, first) {
		t.Fatal("unwrap granted mutation of the aggregate")
	}
}

func TestWrappedRuntimeErrorPreservesJoinedSiblings(t *testing.T) {
	first, second := errors.New("first cause"), errors.New("second cause")
	diagnostic := &RuntimeError{Diagnostic: diagnostics.NewUnexpectedErrorWith(source.NewAnonymous("RETURN 1"), "first", first)}
	joined := errors.Join(fmt.Errorf("host: %w", diagnostic), second)

	result := WrapRuntimeError(nil, 0, nil, joined)
	if result != joined || !errors.Is(result, first) || !errors.Is(result, second) {
		t.Fatalf("runtime error lost joined siblings: %v", result)
	}
}
