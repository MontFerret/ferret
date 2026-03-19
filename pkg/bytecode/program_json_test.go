package bytecode

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestProgramJSONRoundTripPreservesAggregateSelectorSlots(t *testing.T) {
	prog := &Program{
		ISAVersion: Version,
		Registers:  2,
		Bytecode: []Instruction{
			NewInstruction(OpAggregateUpdate, NewRegister(1), NewRegister(2)),
			NewInstruction(OpReturn, NewRegister(1)),
		},
		Metadata: Metadata{
			AggregateSelectorSlots: []int{3, -1},
		},
		Constants: []runtime.Value{runtime.NewInt(1)},
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
}
