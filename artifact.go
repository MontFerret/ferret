package ferret

import "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"

type (
	ArtifactOption = artifact.Option
)

func WithArtifactFormat(format artifact.FormatID) ArtifactOption {
	return artifact.WithFormat(format)
}
