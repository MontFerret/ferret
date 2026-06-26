package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
)

var (
	// WithProgramFormat specifies an artifact format by applying the provided FormatID to the options configuration.
	WithProgramFormat = artifact.WithFormat

	// MarshalProgram serializes the given bytecode program into a byte slice using the provided artifact options.
	// Returns an error if the program is nil or fails during the marshaling process.
	MarshalProgram = artifact.Marshal

	// UnmarshalProgram decodes a byte slice into a *bytecode.Program object or returns an error if unmarshaling fails.
	UnmarshalProgram = artifact.Unmarshal
)

const (
	// ProgramFormatJSON specifies the JSON format identifier for program serialization or artifact representation.
	ProgramFormatJSON = artifact.FormatJSON

	// ProgramFormatMsgPack specifies the MsgPack format identifier for program serialization or artifact representation.
	ProgramFormatMsgPack = artifact.FormatMsgPack
)
