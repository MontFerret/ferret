package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
)

type (
	// ProgramFormat identifies the encoding used for a serialized Ferret program.
	ProgramFormat = artifact.FormatID

	// ProgramOption configures Ferret program serialization.
	ProgramOption = artifact.Option
)

const (
	// ProgramFormatJSON serializes the program payload as JSON.
	ProgramFormatJSON ProgramFormat = artifact.FormatJSON

	// ProgramFormatMsgPack serializes the program payload as MessagePack.
	ProgramFormatMsgPack ProgramFormat = artifact.FormatMsgPack
)

// WithProgramFormat selects the payload format used to serialize a Ferret program.
func WithProgramFormat(format ProgramFormat) ProgramOption {
	return artifact.WithFormat(format)
}

// MarshalProgram is a low-level compatibility helper that serializes a bytecode program.
// Ordinary embedding callers should use Plan.Marshal instead of handling bytecode directly.
func MarshalProgram(program *bytecode.Program, opts ...ProgramOption) ([]byte, error) {
	return artifact.Marshal(program, opts...)
}

// UnmarshalProgram is a low-level compatibility helper that decodes a bytecode program.
// Ordinary embedding callers should use Engine.Load to obtain an executable Plan.
func UnmarshalProgram(data []byte) (*bytecode.Program, error) {
	return artifact.Unmarshal(data)
}
