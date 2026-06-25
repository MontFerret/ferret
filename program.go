package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ArtifactOption is a type alias for artifact.Option, used to define functional options for artifact marshaling.
type ArtifactOption = artifact.Option

// MarshalProgram serializes the given bytecode program into a byte slice using the provided artifact options.
// Returns an error if the program is nil or fails during the marshaling process.
func MarshalProgram(prog *bytecode.Program, opts ...ArtifactOption) ([]byte, error) {
	if prog == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "program cannot be nil")
	}

	return artifact.Marshal(prog, opts...)
}

// UnmarshalProgram decodes a byte slice into a *bytecode.Program object or returns an error if unmarshaling fails.
func UnmarshalProgram(data []byte) (*bytecode.Program, error) {
	return artifact.Unmarshal(data)
}
