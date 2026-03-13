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
					Name:        "boo",
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

	if got, want := decoded.Functions.UserDefined[0].Name, "boo"; got != want {
		t.Fatalf("unexpected exact-case name: got %q, want %q", got, want)
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
					Name:      "boo",
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

func TestProgramJSONRoundTripPreservesCaseDistinctUdfNames(t *testing.T) {
	original := &Program{
		ISAVersion: Version,
		Registers:  1,
		Bytecode: []Instruction{
			NewInstruction(OpReturn, NewRegister(0)),
		},
		Functions: Functions{
			UserDefined: []UDF{
				{Name: "a", DisplayName: "a", Entry: 0, Registers: 1, Params: 0},
				{Name: "A", DisplayName: "A", Entry: 0, Registers: 1, Params: 0},
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

	if got, want := len(decoded.Functions.UserDefined), 2; got != want {
		t.Fatalf("unexpected udf count: got %d, want %d", got, want)
	}

	if got, want := decoded.Functions.UserDefined[0].Name, "a"; got != want {
		t.Fatalf("unexpected lowercase name: got %q, want %q", got, want)
	}

	if got, want := decoded.Functions.UserDefined[1].Name, "A"; got != want {
		t.Fatalf("unexpected uppercase name: got %q, want %q", got, want)
	}
}
