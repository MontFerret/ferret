package artifact

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	programformat "github.com/MontFerret/ferret/v2/pkg/bytecode/format"
)

type Loader struct {
	formats map[FormatID]programformat.Format
}

// NewLoader creates an artifact loader with an explicit set of registered
// payload formats. Nil or duplicate registrations panic.
func NewLoader(formats ...RegisteredFormat) *Loader {
	loader := &Loader{
		formats: make(map[FormatID]programformat.Format, len(formats)),
	}

	for _, registered := range formats {
		if registered.ID == 0 {
			panic("artifact.NewLoader: format id must be non-zero")
		}

		if registered.Format == nil {
			panic(fmt.Sprintf("artifact.NewLoader: format %d is nil", registered.ID))
		}

		if _, exists := loader.formats[registered.ID]; exists {
			panic(fmt.Sprintf("artifact.NewLoader: duplicate format id %d", registered.ID))
		}

		loader.formats[registered.ID] = registered.Format
	}

	return loader
}

// Load decodes a self-describing bytecode artifact and validates header,
// schema, ISA, and payload compatibility before returning a program.
func (l *Loader) Load(data []byte) (*bytecode.Program, error) {
	if l == nil {
		return nil, fmt.Errorf("%w: loader is nil", ErrInvalidArtifact)
	}

	header, err := decodeHeader(data)
	if err != nil {
		return nil, err
	}

	if header.Magic != magic {
		return nil, fmt.Errorf("%w: expected %q, got %q", ErrInvalidMagic, string(magic[:]), string(header.Magic[:]))
	}

	if header.SchemaVersion != schemaVersion {
		return nil, fmt.Errorf("%w: expected %d, got %d", ErrUnsupportedSchema, schemaVersion, header.SchemaVersion)
	}

	if header.Flags != 0 {
		return nil, fmt.Errorf("%w: flags must be 0 in schema version %d, got 0x%x", ErrInvalidHeader, schemaVersion, header.Flags)
	}

	if header.ISAVersion != uint16(bytecode.Version) {
		return nil, fmt.Errorf("%w: expected %d, got %d", ErrIncompatibleISA, bytecode.Version, header.ISAVersion)
	}

	format, exists := l.formats[header.Format]
	if !exists || format == nil {
		return nil, fmt.Errorf("%w: format id %d", ErrUnknownFormat, header.Format)
	}

	if len(data) != headerSize+int(header.PayloadLength) {
		return nil, fmt.Errorf("%w: payload length mismatch: header=%d actual=%d", ErrInvalidHeader, header.PayloadLength, len(data)-headerSize)
	}

	payload := data[headerSize:]
	program, err := format.Unmarshal(payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}

	if program.ISAVersion != int(header.ISAVersion) {
		return nil, fmt.Errorf("%w: payload isaVersion %d does not match header %d", ErrIncompatibleISA, program.ISAVersion, header.ISAVersion)
	}

	return program, nil
}
