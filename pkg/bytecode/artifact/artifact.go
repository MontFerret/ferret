package artifact

import programformat "github.com/MontFerret/ferret/v2/pkg/bytecode/format"

// FormatID identifies a serialized program payload format in an artifact
// header.
type FormatID uint8

const (
	FormatJSON    FormatID = 1
	FormatMsgPack FormatID = 2
	DefaultFormat          = FormatMsgPack
)

type (
	// Options controls artifact marshaling.
	Options struct {
		Format FormatID
	}

	// RegisteredFormat associates an artifact format id with a payload format
	// implementation.
	RegisteredFormat struct {
		Format programformat.Format
		ID     FormatID
	}
)
