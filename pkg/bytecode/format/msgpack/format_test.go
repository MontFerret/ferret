package msgpack

import (
	"errors"
	"testing"
	"time"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/internal/persist"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestFormatRoundTrip(t *testing.T) {
	program := newTestProgram()

	data, err := Default.Marshal(program)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	decoded, err := Default.Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got, want := decoded.Source.Name(), program.Source.Name(); got != want {
		t.Fatalf("unexpected source name: got %q, want %q", got, want)
	}

	if got := decoded.Metadata.AggregatePlans[0].Index["group"]; got != 0 {
		t.Fatalf("expected aggregate plan index to be rebuilt, got %d", got)
	}

	if line, col := decoded.Source.LocationAt(source.Span{Start: 7, End: 7}); line == 0 || col == 0 {
		t.Fatalf("expected source lines to be rebuilt, got line=%d col=%d", line, col)
	}
}

func TestFormatRejectsMalformedPayload(t *testing.T) {
	if _, err := Default.Unmarshal([]byte{0xc1}); err == nil {
		t.Fatal("expected malformed payload error")
	}
}

func TestFormatRejectsMissingRequiredFields(t *testing.T) {
	data := mustMarshalFrame(t, persist.ProgramFrame{})

	_, err := Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidProgram) {
		t.Fatalf("expected ErrInvalidProgram, got %v", err)
	}
}

func TestFormatRejectsInvalidConstants(t *testing.T) {
	frame := validFrame()
	frame.Constants = []persist.ConstantFrame{{Type: "array"}}

	data := mustMarshalFrame(t, frame)
	_, err := Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidConstant) {
		t.Fatalf("expected ErrInvalidConstant, got %v", err)
	}

	frame = validFrame()
	frame.Constants = []persist.ConstantFrame{{Type: "datetime"}}

	data = mustMarshalFrame(t, frame)
	_, err = Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidConstant) {
		t.Fatalf("expected ErrInvalidConstant for missing datetime value, got %v", err)
	}
}

func TestFormatRejectsDuplicateHostsAndLabels(t *testing.T) {
	frame := validFrame()
	frame.Functions.Host = []persist.HostFunctionFrame{
		{Name: "dup", Arity: 1},
		{Name: "dup", Arity: 2},
	}

	data := mustMarshalFrame(t, frame)
	_, err := Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidProgram) {
		t.Fatalf("expected ErrInvalidProgram for duplicate host names, got %v", err)
	}

	frame = validFrame()
	frame.Metadata.Labels = []persist.LabelFrame{
		{PC: 1, Name: "first"},
		{PC: 1, Name: "second"},
	}

	data = mustMarshalFrame(t, frame)
	_, err = Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidProgram) {
		t.Fatalf("expected ErrInvalidProgram for duplicate label pcs, got %v", err)
	}
}

func mustMarshalFrame(t *testing.T, frame persist.ProgramFrame) []byte {
	t.Helper()

	data, err := vmmsgpack.Marshal(frame)
	if err != nil {
		t.Fatalf("marshal frame: %v", err)
	}

	return data
}

func validFrame() persist.ProgramFrame {
	isaVersion := bytecode.Version
	registers := 3
	bytecodeFrame := []persist.InstructionFrame{
		{
			Opcode:   uint8(bytecode.OpLoadConst),
			Operands: [3]int64{0, -1, 0},
		},
		{
			Opcode:   uint8(bytecode.OpReturn),
			Operands: [3]int64{0, 0, 0},
		},
	}
	value := int64(7)

	return persist.ProgramFrame{
		ISAVersion: &isaVersion,
		Registers:  &registers,
		Bytecode:   &bytecodeFrame,
		Constants: []persist.ConstantFrame{
			{Type: "int", Int: &value},
		},
	}
}

func newTestProgram() *bytecode.Program {
	return &bytecode.Program{
		Source: source.New("roundtrip.fql", "RETURN 1\nRETURN 2"),
		Functions: bytecode.Functions{
			Host: map[string]int{
				"now": 0,
				"sum": 2,
			},
			UserDefined: []bytecode.UDF{
				{
					Name:        "main",
					DisplayName: "main",
					Entry:       1,
					Registers:   1,
					Params:      0,
				},
			},
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewFloat(1.5),
			runtime.NewString("hello"),
			runtime.NewBinary([]byte("abc")),
			runtime.NewDateTime(time.Unix(1700000000, 0).UTC()),
			runtime.True,
			runtime.None,
		},
		Params: []string{"input"},
		Metadata: bytecode.Metadata{
			Labels:                 map[int]string{1: "done"},
			CompilerVersion:        "test",
			AggregatePlans:         []bytecode.AggregatePlan{bytecode.NewAggregatePlan([]runtime.String{runtime.NewString("group")}, []bytecode.AggregateKind{bytecode.AggregateCount}, true)},
			AggregateSelectorSlots: []int{-1, -1},
			MatchFailTargets:       []int{-1, -1},
			DebugSpans:             []source.Span{{Start: 0, End: 8}, {Start: 9, End: 17}},
			OptimizationLevel:      1,
		},
		ISAVersion: bytecode.Version,
		Registers:  3,
	}
}
