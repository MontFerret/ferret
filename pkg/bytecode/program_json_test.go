package bytecode

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestProgramJSONPreservesSourcePresence(t *testing.T) {
	tests := []struct {
		name        string
		src         source.Source
		wantPresent bool
	}{
		{name: "zero source"},
		{name: "populated source", src: source.New("query.fql", "RETURN 1"), wantPresent: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			program := &Program{Source: tc.src}

			encoded, err := json.Marshal(program)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			var payload map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &payload); err != nil {
				t.Fatalf("decode payload: %v", err)
			}

			if _, present := payload["source"]; present != tc.wantPresent {
				t.Fatalf("source presence = %t, want %t; payload: %s", present, tc.wantPresent, encoded)
			}

			decoded := Program{Source: source.New("stale.fql", "RETURN 0")}
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			if got, want := decoded.Source.ID(), tc.src.ID(); got != want {
				t.Fatalf("source ID = %v, want %v", got, want)
			}

			if got, want := decoded.Source.Name(), tc.src.Name(); got != want {
				t.Fatalf("source name = %q, want %q", got, want)
			}

			if got, want := decoded.Source.Content(), tc.src.Content(); got != want {
				t.Fatalf("source content = %q, want %q", got, want)
			}
		})
	}

	t.Run("null source", func(t *testing.T) {
		decoded := Program{Source: source.New("stale.fql", "RETURN 0")}
		if err := json.Unmarshal([]byte(`{"source":null}`), &decoded); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		if got, want := decoded.Source.ID(), (source.ID{}); got != want {
			t.Fatalf("source ID = %v, want %v", got, want)
		}
	})
}

func TestProgramJSONRoundTripPreservesInstructionMetadata(t *testing.T) {
	prog := &Program{
		ISAVersion: Version,
		Registers:  2,
		Bytecode: []Instruction{
			NewInstruction(OpAggregateUpdate, NewRegister(1), NewRegister(2)),
			NewInstruction(OpReturn, NewRegister(1)),
		},
		Metadata: Metadata{
			AggregateSelectorSlots: []int{3, -1},
			MatchFailTargets:       []int{-1, 7},
		},
		Constants: []runtime.Value{runtime.NewInt(1), runtime.NewDuration(1500_000_000)},
	}

	encoded, err := json.Marshal(prog)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if strings.Contains(string(encoded), "inlineSlot") {
		t.Fatalf("expected instruction JSON to omit inlineSlot, got %s", string(encoded))
	}

	if !strings.Contains(string(encoded), "aggregateSelectorSlots") {
		t.Fatalf("expected metadata JSON to include aggregateSelectorSlots, got %s", string(encoded))
	}

	if !strings.Contains(string(encoded), "matchFailTargets") {
		t.Fatalf("expected metadata JSON to include matchFailTargets, got %s", string(encoded))
	}
	if !strings.Contains(string(encoded), `"type":"duration","value":"1.5s"`) {
		t.Fatalf("expected normalized duration constant, got %s", string(encoded))
	}

	var decoded Program
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(decoded.Metadata.AggregateSelectorSlots) != 2 {
		t.Fatalf("expected 2 aggregate selector slots, got %d", len(decoded.Metadata.AggregateSelectorSlots))
	}

	if got, want := decoded.Metadata.AggregateSelectorSlots[0], 3; got != want {
		t.Fatalf("unexpected first aggregate selector slot: got %d, want %d", got, want)
	}

	if got, want := decoded.Metadata.AggregateSelectorSlots[1], -1; got != want {
		t.Fatalf("unexpected second aggregate selector slot: got %d, want %d", got, want)
	}

	if len(decoded.Metadata.MatchFailTargets) != 2 {
		t.Fatalf("expected 2 match fail targets, got %d", len(decoded.Metadata.MatchFailTargets))
	}

	if got, want := decoded.Metadata.MatchFailTargets[0], -1; got != want {
		t.Fatalf("unexpected first match fail target: got %d, want %d", got, want)
	}

	if got, want := decoded.Metadata.MatchFailTargets[1], 7; got != want {
		t.Fatalf("unexpected second match fail target: got %d, want %d", got, want)
	}
	if got, ok := decoded.Constants[1].(runtime.Duration); !ok || got != runtime.NewDuration(1500_000_000) {
		t.Fatalf("duration constant = %v (%T)", decoded.Constants[1], decoded.Constants[1])
	}
}

func TestProgramJSONRejectsMalformedDurationConstant(t *testing.T) {
	_, err := decodeConstant(constantJSON{
		Type:  "duration",
		Value: json.RawMessage(`"not-a-duration"`),
	})
	if err == nil {
		t.Fatal("expected malformed duration constant error")
	}
}

func TestProgramJSONPreservesOrderedHostSignatures(t *testing.T) {
	program := &Program{
		ISAVersion: Version,
		Registers:  1,
		Bytecode: []Instruction{
			NewInstruction(OpReturn, NewRegister(0)),
		},
		Functions: Functions{Host: []HostFunction{
			{Name: "DB::POSTGRES::PICK", ArgCount: 2},
			{Name: "Db::Postgres::Pick", ArgCount: 1},
		}},
	}

	encoded, err := json.Marshal(program)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Program
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !slices.Equal(decoded.Functions.Host, program.Functions.Host) {
		t.Fatalf("unexpected host signature order or spelling: got %v, want %v", decoded.Functions.Host, program.Functions.Host)
	}
}

func TestProgramJSONRejectsLegacyHostMetadata(t *testing.T) {
	for _, payload := range []string{
		`{"isaVersion":1,"registers":1,"bytecode":[{"opcode":39,"operands":[0,0,0]}],"functions":{"host":[{"name":"PICK"}]}}`,
		`{"isaVersion":1,"registers":1,"bytecode":[{"opcode":39,"operands":[0,0,0]}],"functions":{"host":[{"name":"PICK","arity":1}]}}`,
	} {
		var decoded Program
		if err := json.Unmarshal([]byte(payload), &decoded); err == nil || !strings.Contains(err.Error(), "missing argCount") {
			t.Fatalf("expected legacy host metadata rejection, got %v", err)
		}
	}
}

func TestProgramJSONRejectsDuplicateHostSignatures(t *testing.T) {
	payload := `{"isaVersion":1,"registers":1,"bytecode":[{"opcode":39,"operands":[0,0,0]}],"functions":{"host":[{"name":"DB::POSTGRES::PICK","argCount":1},{"name":"db::postgres::pick","argCount":1}]}}`

	var decoded Program
	if err := json.Unmarshal([]byte(payload), &decoded); err == nil || !strings.Contains(err.Error(), "duplicate host function signature") {
		t.Fatalf("expected duplicate host signature rejection, got %v", err)
	}
}
