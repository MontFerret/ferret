package json

import (
	"errors"
	"slices"
	"testing"
	"time"

	gojson "github.com/goccy/go-json"

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

	if got, want := decoded.Source.Content(), program.Source.Content(); got != want {
		t.Fatalf("unexpected source content: got %q, want %q", got, want)
	}

	if position := decoded.Source.PositionAt(source.Span{Start: 7, End: 7}); position.Line == 0 || position.Column == 0 {
		t.Fatalf("expected source lines to be rebuilt, got position=%+v", position)
	}

	if got := decoded.Metadata.AggregatePlans[0].Index["group"]; got != 0 {
		t.Fatalf("expected aggregate plan index to be rebuilt, got %d", got)
	}

	if got := decoded.Constants[1]; got != runtime.NewDuration(1500*time.Millisecond) {
		t.Fatalf("duration constant = %v (%T)", got, got)
	}
	if !slices.Equal(decoded.Functions.Host, program.Functions.Host) {
		t.Fatalf("host signature order mismatch: got %v, want %v", decoded.Functions.Host, program.Functions.Host)
	}
}

func TestFormatRejectsMalformedPayload(t *testing.T) {
	if _, err := Default.Unmarshal([]byte("{")); err == nil {
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
	frame.Constants = []persist.ConstantFrame{{Type: "bool"}}

	data = mustMarshalFrame(t, frame)
	_, err = Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidConstant) {
		t.Fatalf("expected ErrInvalidConstant for missing bool value, got %v", err)
	}

	invalidDuration := "not-a-duration"
	frame = validFrame()
	frame.Constants = []persist.ConstantFrame{{Type: "duration", Duration: &invalidDuration}}

	data = mustMarshalFrame(t, frame)
	_, err = Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidConstant) {
		t.Fatalf("expected ErrInvalidConstant for malformed duration, got %v", err)
	}
}

func TestFormatAllowsOverloadedHostsAndRejectsDuplicateSignaturesAndLabels(t *testing.T) {
	frame := validFrame()
	one := 1
	two := 2
	frame.Functions.Host = []persist.HostFunctionFrame{
		{Name: "DB::POSTGRES::DUP", ArgCount: &one},
		{Name: "db::postgres::dup", ArgCount: &two},
	}

	data := mustMarshalFrame(t, frame)
	decoded, err := Default.Unmarshal(data)
	if err != nil {
		t.Fatalf("expected overloaded host signatures to decode, got %v", err)
	}
	if got := decoded.Functions.Host[0].Name; got != "DB::POSTGRES::DUP" {
		t.Fatalf("expected stored host metadata spelling to survive, got %q", got)
	}

	frame.Functions.Host[1].ArgCount = &one
	data = mustMarshalFrame(t, frame)
	_, err = Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidProgram) {
		t.Fatalf("expected ErrInvalidProgram for duplicate host signatures, got %v", err)
	}

	frame = validFrame()
	frame.Functions.Host = []persist.HostFunctionFrame{{Name: "legacy"}}
	data = mustMarshalFrame(t, frame)
	_, err = Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidProgram) {
		t.Fatalf("expected ErrInvalidProgram for missing argCount, got %v", err)
	}

	var legacy map[string]any
	if err := gojson.Unmarshal(mustMarshalFrame(t, validFrame()), &legacy); err != nil {
		t.Fatalf("decode legacy test payload: %v", err)
	}
	functions := legacy["functions"].(map[string]any)
	functions["host"] = []any{map[string]any{"name": "legacy", "arity": 2}}
	data, err = gojson.Marshal(legacy)
	if err != nil {
		t.Fatalf("encode legacy test payload: %v", err)
	}
	_, err = Default.Unmarshal(data)
	if !errors.Is(err, bytecode.ErrInvalidProgram) {
		t.Fatalf("expected ErrInvalidProgram for legacy arity metadata, got %v", err)
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

	data, err := gojson.Marshal(frame)
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
			Host: []bytecode.HostFunction{
				{Name: "sum", ArgCount: 2},
				{Name: "sum", ArgCount: 1},
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
			runtime.NewDuration(1500 * time.Millisecond),
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
