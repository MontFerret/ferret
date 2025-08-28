package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestIsQuote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"double quote", "\"", true},
		{"single quote", "'", true},
		{"backtick", "`", true},
		{"empty string", "", false},
		{"letter", "a", false},
		{"multiple chars", "abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isQuote(tt.input)
			if result != tt.expected {
				t.Errorf("isQuote(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"single quote", "\"", true},
		{"single single quote", "'", true},
		{"single backtick", "`", true},
		{"double quoted string", "\"hello\"", true},
		{"single quoted string", "'hello'", true},
		{"backtick quoted string", "`hello`", true},
		{"mismatched quotes", "\"hello'", false},
		{"no quotes", "hello", false},
		{"single character", "h", false},
		{"mixed quotes start with double", "\"hello'", false},
		{"mixed quotes start with single", "'hello\"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidString(tt.input)
			if result != tt.expected {
				t.Errorf("isValidString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSpanFromTokenSafe_EdgeCases(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1") // Length: 9

	// Test nil token
	result := spanFromTokenSafe(nil, src)
	expected := file.Span{Start: 0, End: 1}
	if result != expected {
		t.Errorf("spanFromTokenSafe(nil, src) = %v, want %v", result, expected)
	}
}

func TestIsIdentifier_NilCases(t *testing.T) {
	if isIdentifier(nil) {
		t.Error("isIdentifier(nil) should return false")
	}

	node := &TokenNode{token: nil}
	if isIdentifier(node) {
		t.Error("isIdentifier with nil token should return false")
	}
}

func TestIsKeyword_NilCases(t *testing.T) {
	if isKeyword(nil) {
		t.Error("isKeyword(nil) should return false")
	}

	node := &TokenNode{token: nil}
	if isKeyword(node) {
		t.Error("isKeyword with nil token should return false")
	}
}
