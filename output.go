package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// Output is the encoded result returned from session or engine execution.
	Output = encoding.Output

	// Value is the shared Ferret value contract accepted at the embedding boundary.
	Value = runtime.Value

	// Params contains pre-converted Ferret values keyed by query parameter name.
	Params = runtime.Params
)
