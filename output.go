package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// Output is the caller-owned, one-shot handle returned by execution.
	Output = encoding.Output
	// Content is detached encoded data obtained through collection or debugging.
	Content = encoding.Content
	// Metadata is the immutable descriptor of a complete encoded representation.
	Metadata = encoding.Metadata
	// Consumer receives synchronous, borrowed, read-only encoded chunks.
	Consumer = encoding.Consumer

	// Value is the shared Ferret value contract accepted at the embedding boundary.
	Value = runtime.Value

	// Params contains pre-converted Ferret values keyed by query parameter name.
	Params = runtime.Params
)

var (
	// ErrOutputInUse identifies a competing consumption attempt.
	ErrOutputInUse = encoding.ErrOutputInUse
	// ErrOutputClosed identifies consumption after stopping, or interrupted by Close.
	ErrOutputClosed = encoding.ErrOutputClosed
)
