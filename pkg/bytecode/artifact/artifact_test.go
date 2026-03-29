package artifact

import (
	"bytes"
	"errors"
	"testing"

	gojson "github.com/goccy/go-json"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	formatjson "github.com/MontFerret/ferret/v2/pkg/bytecode/format/json"
	formatmsgpack "github.com/MontFerret/ferret/v2/pkg/bytecode/format/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/internal/persist"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestMarshalAndUnmarshal_DefaultMessagePack(t *testing.T) {
	program := newArtifactTestProgram()

	data, err := Marshal(program, Options{})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	decoded, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got, want := decoded.ISAVersion, program.ISAVersion; got != want {
		t.Fatalf("unexpected isaVersion: got %d, want %d", got, want)
	}
}

func TestMarshalAndUnmarshal_JSON(t *testing.T) {
	program := newArtifactTestProgram()

	data, err := Marshal(program, Options{Format: FormatJSON})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	header, err := decodeHeader(data)
	if err != nil {
		t.Fatalf("decodeHeader() error = %v", err)
	}

	if got, want := header.Format, FormatJSON; got != want {
		t.Fatalf("unexpected format id: got %d, want %d", got, want)
	}

	if _, err := Unmarshal(data); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
}

func TestMarshalRejectsUnknownFormat(t *testing.T) {
	_, err := Marshal(newArtifactTestProgram(), Options{Format: FormatID(99)})
	if !errors.Is(err, ErrUnknownFormat) {
		t.Fatalf("expected ErrUnknownFormat, got %v", err)
	}
}

func TestMarshalAllowsConcatImmediateCountAtRegisterLimit(t *testing.T) {
	program := newArtifactTestProgram()
	program.Registers = 3
	program.Bytecode = []bytecode.Instruction{
		bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(0), bytecode.NewRegister(0), bytecode.Operand(program.Registers)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	}

	if _, err := Marshal(program, Options{}); err != nil {
		t.Fatalf("expected Marshal() to accept valid concat immediate count, got %v", err)
	}
}

func TestPayloadLengthForHeader(t *testing.T) {
	t.Run("max_uint32", func(t *testing.T) {
		length, err := payloadLengthForHeader(^uint64(0) >> 32)
		if err != nil {
			t.Fatalf("expected max uint32 payload length to be valid, got %v", err)
		}

		if got, want := length, ^uint32(0); got != want {
			t.Fatalf("unexpected payload length: got %d, want %d", got, want)
		}
	})

	t.Run("overflow", func(t *testing.T) {
		_, err := payloadLengthForHeader(uint64(^uint32(0)) + 1)
		if !errors.Is(err, ErrInvalidHeader) {
			t.Fatalf("expected ErrInvalidHeader, got %v", err)
		}
	})
}

func TestLoaderRejectsInvalidHeaders(t *testing.T) {
	program := newArtifactTestProgram()
	data, err := Marshal(program, Options{})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	t.Run("short_input", func(t *testing.T) {
		_, err := Unmarshal([]byte("x"))
		if !errors.Is(err, ErrInvalidHeader) {
			t.Fatalf("expected ErrInvalidHeader, got %v", err)
		}
	})

	t.Run("invalid_magic", func(t *testing.T) {
		mutated := append([]byte(nil), data...)
		copy(mutated[:4], []byte("XXXX"))

		_, err := Unmarshal(mutated)
		if !errors.Is(err, ErrInvalidMagic) {
			t.Fatalf("expected ErrInvalidMagic, got %v", err)
		}
	})

	t.Run("unsupported_schema", func(t *testing.T) {
		mutated := append([]byte(nil), data...)
		mutated[5]++

		_, err := Unmarshal(mutated)
		if !errors.Is(err, ErrUnsupportedSchema) {
			t.Fatalf("expected ErrUnsupportedSchema, got %v", err)
		}
	})

	t.Run("unknown_format", func(t *testing.T) {
		mutated := append([]byte(nil), data...)
		mutated[4] = byte(99)

		_, err := Unmarshal(mutated)
		if !errors.Is(err, ErrUnknownFormat) {
			t.Fatalf("expected ErrUnknownFormat, got %v", err)
		}
	})

	t.Run("unregistered_format", func(t *testing.T) {
		jsonData, err := Marshal(program, Options{Format: FormatJSON})
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}

		loader := NewLoader(RegisteredFormat{ID: FormatMsgPack, Format: formatmsgpack.Default})
		_, err = loader.Load(jsonData)
		if !errors.Is(err, ErrUnknownFormat) {
			t.Fatalf("expected ErrUnknownFormat, got %v", err)
		}
	})

	t.Run("non_zero_flags", func(t *testing.T) {
		mutated := append([]byte(nil), data...)
		header, err := decodeHeader(mutated)
		if err != nil {
			t.Fatalf("decodeHeader() error = %v", err)
		}

		header.Flags = 1
		encodeHeader(mutated[:headerSize], header)

		_, err = Unmarshal(mutated)
		if !errors.Is(err, ErrInvalidHeader) {
			t.Fatalf("expected ErrInvalidHeader, got %v", err)
		}
	})

	t.Run("payload_length_mismatch", func(t *testing.T) {
		mutated := append([]byte(nil), data...)
		header, err := decodeHeader(mutated)
		if err != nil {
			t.Fatalf("decodeHeader() error = %v", err)
		}

		header.PayloadLength++
		encodeHeader(mutated[:headerSize], header)

		_, err = Unmarshal(mutated)
		if !errors.Is(err, ErrInvalidHeader) {
			t.Fatalf("expected ErrInvalidHeader, got %v", err)
		}
	})

	t.Run("trailing_bytes", func(t *testing.T) {
		mutated := append(append([]byte(nil), data...), 0x00)

		_, err := Unmarshal(mutated)
		if !errors.Is(err, ErrInvalidHeader) {
			t.Fatalf("expected ErrInvalidHeader, got %v", err)
		}
	})

	t.Run("incompatible_isa", func(t *testing.T) {
		mutated := append([]byte(nil), data...)
		header, err := decodeHeader(mutated)
		if err != nil {
			t.Fatalf("decodeHeader() error = %v", err)
		}

		header.ISAVersion++
		encodeHeader(mutated[:headerSize], header)

		_, err = Unmarshal(mutated)
		if !errors.Is(err, ErrIncompatibleISA) {
			t.Fatalf("expected ErrIncompatibleISA, got %v", err)
		}
	})

	t.Run("payload_header_isa_mismatch", func(t *testing.T) {
		jsonData, err := Marshal(program, Options{Format: FormatJSON})
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}

		header, err := decodeHeader(jsonData)
		if err != nil {
			t.Fatalf("decodeHeader() error = %v", err)
		}

		payload := append([]byte(nil), jsonData[headerSize:]...)
		payload = bytes.Replace(payload, []byte(`"isaVersion":1`), []byte(`"isaVersion":2`), 1)

		mutated := make([]byte, headerSize+len(payload))
		copy(mutated[:headerSize], jsonData[:headerSize])
		copy(mutated[headerSize:], payload)
		header.PayloadLength = uint32(len(payload))
		encodeHeader(mutated[:headerSize], header)

		_, err = Unmarshal(mutated)
		if !errors.Is(err, ErrInvalidPayload) && !errors.Is(err, ErrIncompatibleISA) {
			t.Fatalf("expected payload/header ISA mismatch error, got %v", err)
		}
	})
}

func TestNewLoaderPanicsOnInvalidRegistrations(t *testing.T) {
	t.Run("nil_format", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()

		NewLoader(RegisteredFormat{ID: FormatJSON})
	})

	t.Run("duplicate_id", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()

		NewLoader(
			RegisteredFormat{ID: FormatJSON, Format: formatjson.Default},
			RegisteredFormat{ID: FormatJSON, Format: formatjson.Default},
		)
	})
}

func newArtifactTestProgram() *bytecode.Program {
	return &bytecode.Program{
		Source: source.New("artifact.fql", "RETURN 1"),
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Constants: []runtime.Value{
			runtime.NewString("ok"),
		},
		Metadata: bytecode.Metadata{
			Labels:                 map[int]string{1: "done"},
			CompilerVersion:        "test",
			AggregatePlans:         []bytecode.AggregatePlan{bytecode.NewAggregatePlan([]runtime.String{runtime.NewString("group")}, []bytecode.AggregateKind{bytecode.AggregateCount}, false)},
			AggregateSelectorSlots: []int{-1, -1},
			MatchFailTargets:       []int{-1, -1},
			DebugSpans:             []source.Span{{Start: 0, End: 8}, {Start: 0, End: 8}},
			OptimizationLevel:      1,
		},
		ISAVersion: bytecode.Version,
		Registers:  1,
	}
}

func TestLoaderRejectsMalformedPayload(t *testing.T) {
	isaVersion := bytecode.Version
	registers := 1
	bytecodeFrame := []persist.InstructionFrame{
		{
			Opcode:   uint8(bytecode.OpReturn),
			Operands: [3]int64{0, 0, 0},
		},
	}

	frame := persist.ProgramFrame{
		ISAVersion: &isaVersion,
		Registers:  &registers,
		Bytecode:   &bytecodeFrame,
		Constants: []persist.ConstantFrame{
			{Type: "array"},
		},
	}

	payload, err := gojson.Marshal(frame)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	header := header{
		Magic:         magic,
		Format:        FormatJSON,
		SchemaVersion: schemaVersion,
		ISAVersion:    uint16(bytecode.Version),
		PayloadLength: uint32(len(payload)),
	}

	data := make([]byte, headerSize+len(payload))
	encodeHeader(data[:headerSize], header)
	copy(data[headerSize:], payload)

	_, err = Unmarshal(data)
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestLoaderRejectsOversizedCollectionPreallocation(t *testing.T) {
	isaVersion := bytecode.Version
	registers := 1
	bytecodeFrame := []persist.InstructionFrame{
		{
			Opcode:   uint8(bytecode.OpLoadArray),
			Operands: [3]int64{0, int64(bytecode.MaxCollectionPreallocation + 1), 0},
		},
		{
			Opcode:   uint8(bytecode.OpReturn),
			Operands: [3]int64{0, 0, 0},
		},
	}

	frame := persist.ProgramFrame{
		ISAVersion: &isaVersion,
		Registers:  &registers,
		Bytecode:   &bytecodeFrame,
	}

	payload, err := gojson.Marshal(frame)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	header := header{
		Magic:         magic,
		Format:        FormatJSON,
		SchemaVersion: schemaVersion,
		ISAVersion:    uint16(bytecode.Version),
		PayloadLength: uint32(len(payload)),
	}

	data := make([]byte, headerSize+len(payload))
	encodeHeader(data[:headerSize], header)
	copy(data[headerSize:], payload)

	_, err = Unmarshal(data)
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}

	if !errors.Is(err, bytecode.ErrInvalidInstruction) {
		t.Fatalf("expected bytecode.ErrInvalidInstruction, got %v", err)
	}
}

func TestLoaderRejectsOversizedMultiSorterDirectionCount(t *testing.T) {
	isaVersion := bytecode.Version
	registers := 1
	bytecodeFrame := []persist.InstructionFrame{
		{
			Opcode:   uint8(bytecode.OpDataSetMultiSorter),
			Operands: [3]int64{0, 0, int64(bytecode.MaxEncodedSortDirections + 1)},
		},
		{
			Opcode:   uint8(bytecode.OpReturn),
			Operands: [3]int64{0, 0, 0},
		},
	}

	frame := persist.ProgramFrame{
		ISAVersion: &isaVersion,
		Registers:  &registers,
		Bytecode:   &bytecodeFrame,
	}

	payload, err := gojson.Marshal(frame)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	header := header{
		Magic:         magic,
		Format:        FormatJSON,
		SchemaVersion: schemaVersion,
		ISAVersion:    uint16(bytecode.Version),
		PayloadLength: uint32(len(payload)),
	}

	data := make([]byte, headerSize+len(payload))
	encodeHeader(data[:headerSize], header)
	copy(data[headerSize:], payload)

	_, err = Unmarshal(data)
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}

	if !errors.Is(err, bytecode.ErrInvalidInstruction) {
		t.Fatalf("expected bytecode.ErrInvalidInstruction, got %v", err)
	}
}

func TestLoaderRejectsConcatRangeOutOfBounds(t *testing.T) {
	isaVersion := bytecode.Version
	registers := 3
	bytecodeFrame := []persist.InstructionFrame{
		{
			Opcode:   uint8(bytecode.OpConcat),
			Operands: [3]int64{0, 2, 2},
		},
		{
			Opcode:   uint8(bytecode.OpReturn),
			Operands: [3]int64{0, 0, 0},
		},
	}

	frame := persist.ProgramFrame{
		ISAVersion: &isaVersion,
		Registers:  &registers,
		Bytecode:   &bytecodeFrame,
	}

	payload, err := gojson.Marshal(frame)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	header := header{
		Magic:         magic,
		Format:        FormatJSON,
		SchemaVersion: schemaVersion,
		ISAVersion:    uint16(bytecode.Version),
		PayloadLength: uint32(len(payload)),
	}

	data := make([]byte, headerSize+len(payload))
	encodeHeader(data[:headerSize], header)
	copy(data[headerSize:], payload)

	_, err = Unmarshal(data)
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}

	if !errors.Is(err, bytecode.ErrInvalidInstruction) {
		t.Fatalf("expected bytecode.ErrInvalidInstruction, got %v", err)
	}
}

func TestLoaderAcceptsMaxEncodedSortDirections(t *testing.T) {
	isaVersion := bytecode.Version
	registers := 1
	bytecodeFrame := []persist.InstructionFrame{
		{
			Opcode:   uint8(bytecode.OpDataSetMultiSorter),
			Operands: [3]int64{0, int64(1 << (bytecode.MaxEncodedSortDirections - 1)), int64(bytecode.MaxEncodedSortDirections)},
		},
		{
			Opcode:   uint8(bytecode.OpReturn),
			Operands: [3]int64{0, 0, 0},
		},
	}

	frame := persist.ProgramFrame{
		ISAVersion: &isaVersion,
		Registers:  &registers,
		Bytecode:   &bytecodeFrame,
	}

	payload, err := gojson.Marshal(frame)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	header := header{
		Magic:         magic,
		Format:        FormatJSON,
		SchemaVersion: schemaVersion,
		ISAVersion:    uint16(bytecode.Version),
		PayloadLength: uint32(len(payload)),
	}

	data := make([]byte, headerSize+len(payload))
	encodeHeader(data[:headerSize], header)
	copy(data[headerSize:], payload)

	if _, err := Unmarshal(data); err != nil {
		t.Fatalf("expected max encoded sort direction count to be accepted, got %v", err)
	}
}
