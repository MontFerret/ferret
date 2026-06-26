package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
)

// ProgramArtifactOption is a type alias for artifact.Option, used to define functional options for artifact marshaling.
type ProgramArtifactOption = artifact.Option

var (
	// WithProgramFormat specifies an artifact format by applying the provided FormatID to the options configuration.
	WithProgramFormat = artifact.WithFormat

	// MarshalProgram serializes the given bytecode program into a byte slice using the provided artifact options.
	// Returns an error if the program is nil or fails during the marshaling process.
	MarshalProgram = artifact.Marshal

	// UnmarshalProgram decodes a byte slice into a *bytecode.Program object or returns an error if unmarshaling fails.
	UnmarshalProgram = artifact.Unmarshal
)
