package bytecode

import (
	"strings"
	"testing"

	"github.com/goccy/go-json"
)

func TestProgramJSONRoundTripPreservesUdfDisplayName(t *testing.T) {
	original := &Program{
		ISAVersion: Version,
		Registers:  1,
		Bytecode: []Instruction{
			NewInstruction(OpReturn, NewRegister(0)),
		},
		Functions: Functions{
			UserDefined: []UDF{
				{
					Name:        "BOO",
					DisplayName: "boo",
					Entry:       0,
					Registers:   2,
					Params:      0,
				},
			},
		},
	}

	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Program
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if got, want := len(decoded.Functions.UserDefined), 1; got != want {
		t.Fatalf("unexpected udf count: got %d, want %d", got, want)
	}

	if got, want := decoded.Functions.UserDefined[0].Name, "BOO"; got != want {
		t.Fatalf("unexpected normalized name: got %q, want %q", got, want)
	}

	if got, want := decoded.Functions.UserDefined[0].DisplayName, "boo"; got != want {
		t.Fatalf("unexpected display name: got %q, want %q", got, want)
	}
}

func TestProgramJSONDecodeWithoutUdfDisplayNameKeepsFallbackEmpty(t *testing.T) {
	original := &Program{
		ISAVersion: Version,
		Registers:  1,
		Bytecode: []Instruction{
			NewInstruction(OpReturn, NewRegister(0)),
		},
		Functions: Functions{
			UserDefined: []UDF{
				{
					Name:      "BOO",
					Entry:     0,
					Registers: 2,
					Params:    0,
				},
			},
		},
	}

	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Simulate older bytecode payloads with no displayName field.
	if strings.Contains(string(encoded), "\"displayName\"") {
		t.Fatalf("unexpected displayName in encoded payload: %s", string(encoded))
	}

	var decoded Program
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if got, want := len(decoded.Functions.UserDefined), 1; got != want {
		t.Fatalf("unexpected udf count: got %d, want %d", got, want)
	}

	if got := decoded.Functions.UserDefined[0].DisplayName; got != "" {
		t.Fatalf("expected empty display name fallback for legacy payload, got %q", got)
	}
}
