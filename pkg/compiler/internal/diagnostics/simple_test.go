package diagnostics

import (
	"testing"
)

func TestSimpleFunctions(t *testing.T) {
	// Test NewTrackingTokenStream constructor
	stream := NewTrackingTokenStream(nil, nil)
	if stream == nil {
		t.Error("NewTrackingTokenStream should not return nil")
	}

	// Test is() function with nil node
	result := is(nil, "TEST")
	if result {
		t.Error("is(nil, 'TEST') should return false")
	}

	// Test anyIs with nil nodes
	result2 := anyIs(nil, nil, "TEST")
	if result2 != nil {
		t.Error("anyIs(nil, nil, 'TEST') should return nil")
	}
}

func TestSpanFromToken_Simple(t *testing.T) {
	// Test SpanFromToken with nil
	result := SpanFromToken(nil)
	if result.Start != 0 || result.End != 0 {
		t.Errorf("SpanFromToken(nil) = %v, want {0 0}", result)
	}
}

func TestTokenNode_Simple(t *testing.T) {
	// Test TokenNode methods with nil token
	node := &TokenNode{token: nil}

	if node.GetText() != "" {
		t.Error("GetText() with nil token should return empty string")
	}

	if node.String() != "" {
		t.Error("String() with nil token should return empty string")
	}

	if node.Token() != nil {
		t.Error("Token() should return nil when token is nil")
	}

	if node.Prev() != nil {
		t.Error("Prev() should return nil when no previous node")
	}

	if node.Next() != nil {
		t.Error("Next() should return nil when no next node")
	}
}
