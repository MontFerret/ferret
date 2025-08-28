package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestHas(t *testing.T) {
	tests := []struct {
		name   string
		msg    string
		substr string
		want   bool
	}{
		{"contains substr", "This is a test message", "test", true},
		{"case insensitive", "This is a TEST message", "test", true},
		{"not contains", "This is a message", "example", false},
		{"empty substr", "This is a message", "", true},
		{"empty msg", "", "test", false},
		{"both empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := has(tt.msg, tt.substr)
			if result != tt.want {
				t.Errorf("has(%q, %q) = %v, want %v", tt.msg, tt.substr, result, tt.want)
			}
		})
	}
}

func TestIsMismatched(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"contains mismatched input", "mismatched input 'FOR' expecting", true},
		{"case insensitive", "MISMATCHED INPUT 'FOR' expecting", true},
		{"does not contain", "syntax error at position 5", false},
		{"partial match", "matched input error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMismatched(tt.msg)
			if result != tt.want {
				t.Errorf("isMismatched(%q) = %v, want %v", tt.msg, result, tt.want)
			}
		})
	}
}

func TestIsNoAlternative(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"contains no viable alternative", "no viable alternative at input 'FOR'", true},
		{"case insensitive", "NO VIABLE ALTERNATIVE at input 'FOR'", true},
		{"does not contain", "syntax error at position 5", false},
		{"partial match", "no alternative found", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNoAlternative(tt.msg)
			if result != tt.want {
				t.Errorf("isNoAlternative(%q) = %v, want %v", tt.msg, result, tt.want)
			}
		})
	}
}

func TestIsMissing(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"contains missing", "missing ';' at 'FOR'", true},
		{"case insensitive", "MISSING token", true},
		{"does not contain", "syntax error at position 5", false},
		{"partial match", "miss token", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMissing(tt.msg)
			if result != tt.want {
				t.Errorf("isMissing(%q) = %v, want %v", tt.msg, result, tt.want)
			}
		})
	}
}

func TestIsMissingToken(t *testing.T) {
	tests := []struct {
		name  string
		msg   string
		token string
		want  bool
	}{
		{"contains both missing and token", "missing ';' at 'FOR'", ";", true},
		{"case insensitive", "MISSING ';' AT 'FOR'", ";", true},
		{"missing but wrong token", "missing ';' at 'FOR'", ":", false},
		{"no missing", "syntax error at ';'", ";", false},
		{"no token", "missing something", ";", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMissingToken(tt.msg, tt.token)
			if result != tt.want {
				t.Errorf("isMissingToken(%q, %q) = %v, want %v", tt.msg, tt.token, result, tt.want)
			}
		})
	}
}

func TestExtractNoAlternativeInput(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "valid pattern",
			msg:  "no viable alternative at input 'FOR x'",
			want: "FOR x",
		},
		{
			name: "with quotes",
			msg:  "no viable alternative at input 'LET'",
			want: "LET",
		},
		{
			name: "no match",
			msg:  "some other error message",
			want: "",
		},
		{
			name: "empty message",
			msg:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractNoAlternativeInput(tt.msg)
			if result != tt.want {
				t.Errorf("extractNoAlternativeInput(%q) = %q, want %q", tt.msg, result, tt.want)
			}
		})
	}
}

func TestExtractNoAlternativeInputs(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want []string
	}{
		{
			name: "multiple words",
			msg:  "no viable alternative at input 'FOR x IN'",
			want: []string{"FOR", "x", "IN"},
		},
		{
			name: "single word",
			msg:  "no viable alternative at input 'LET'",
			want: []string{"LET"},
		},
		{
			name: "no match",
			msg:  "some other error message",
			want: []string{},
		},
		{
			name: "empty message",
			msg:  "",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractNoAlternativeInputs(tt.msg)
			if len(result) != len(tt.want) {
				t.Errorf("extractNoAlternativeInputs(%q) length = %d, want %d", tt.msg, len(result), len(tt.want))
				return
			}
			for i, expected := range tt.want {
				if result[i] != expected {
					t.Errorf("extractNoAlternativeInputs(%q)[%d] = %q, want %q", tt.msg, i, result[i], expected)
				}
			}
		})
	}
}

func TestIsExtraneous(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"contains extraneous input", "extraneous input 'FOR' expecting", true},
		{"case insensitive", "EXTRANEOUS INPUT 'FOR' expecting", true},
		{"does not contain", "syntax error at position 5", false},
		{"partial match", "extra input error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExtraneous(tt.msg)
			if result != tt.want {
				t.Errorf("isExtraneous(%q) = %v, want %v", tt.msg, result, tt.want)
			}
		})
	}
}

func TestExtractExtraneousInput(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "valid pattern",
			msg:  "extraneous input 'FOR' expecting ';'",
			want: "FOR",
		},
		{
			name: "with quotes",
			msg:  "extraneous input 'LET x' expecting '='",
			want: "LET x",
		},
		{
			name: "no match",
			msg:  "some other error message",
			want: "",
		},
		{
			name: "empty message",
			msg:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractExtraneousInput(tt.msg)
			if result != tt.want {
				t.Errorf("extractExtraneousInput(%q) = %q, want %q", tt.msg, result, tt.want)
			}
		})
	}
}

func TestExtractMismatchedInput(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "valid pattern",
			msg:  "mismatched input 'FOR' expecting ';'",
			want: "FOR",
		},
		{
			name: "with quotes",
			msg:  "mismatched input 'LET x' expecting '='",
			want: "LET x",
		},
		{
			name: "no match",
			msg:  "some other error message",
			want: "",
		},
		{
			name: "empty message",
			msg:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractMismatchedInput(tt.msg)
			if result != tt.want {
				t.Errorf("extractMismatchedInput(%q) = %q, want %q", tt.msg, result, tt.want)
			}
		})
	}
}

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
