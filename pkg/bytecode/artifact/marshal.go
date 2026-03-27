package artifact

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	programformat "github.com/MontFerret/ferret/v2/pkg/bytecode/format"
	formatjson "github.com/MontFerret/ferret/v2/pkg/bytecode/format/json"
	formatmsgpack "github.com/MontFerret/ferret/v2/pkg/bytecode/format/msgpack"
)

var (
	builtinFormats = map[FormatID]programformat.Format{
		FormatJSON:    formatjson.Default,
		FormatMsgPack: formatmsgpack.Default,
	}
	defaultLoader = NewDefaultLoader()
)

func builtinRegisteredFormats() []RegisteredFormat {
	return []RegisteredFormat{
		{ID: FormatJSON, Format: formatjson.Default},
		{ID: FormatMsgPack, Format: formatmsgpack.Default},
	}
}

// Marshal serializes a program into a self-describing artifact using one of the
// built-in payload formats.
func Marshal(program *bytecode.Program, opts Options) ([]byte, error) {
	if err := bytecode.ValidateProgram(program); err != nil {
		return nil, err
	}

	if program.ISAVersion < 0 || program.ISAVersion > int(^uint16(0)) {
		return nil, fmt.Errorf("%w: program isaVersion %d overflows header field", ErrInvalidHeader, program.ISAVersion)
	}

	formatID := opts.Format
	if formatID == 0 {
		formatID = DefaultFormat
	}

	format, exists := builtinFormats[formatID]
	if !exists || format == nil {
		return nil, fmt.Errorf("%w: format id %d", ErrUnknownFormat, formatID)
	}

	payload, err := format.Marshal(program)
	if err != nil {
		return nil, err
	}

	if len(payload) > int(^uint32(0)) {
		return nil, fmt.Errorf("%w: payload length %d overflows header field", ErrInvalidHeader, len(payload))
	}

	header := header{
		Magic:         magic,
		Format:        formatID,
		SchemaVersion: schemaVersion,
		ISAVersion:    uint16(program.ISAVersion),
		Flags:         0,
		PayloadLength: uint32(len(payload)),
	}

	data := make([]byte, headerSize+len(payload))
	encodeHeader(data[:headerSize], header)
	copy(data[headerSize:], payload)

	return data, nil
}

// Unmarshal decodes a self-describing artifact using the built-in loader.
func Unmarshal(data []byte) (*bytecode.Program, error) {
	return defaultLoader.Load(data)
}
