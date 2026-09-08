package ferret

import "github.com/MontFerret/ferret/v2/pkg/engine"

type (
	// Output is the encoded result returned from session or engine execution.
	Output = engine.Output

	// Value is the shared Ferret value contract accepted at the embedding boundary.
	Value = engine.Value

	// Params contains pre-converted Ferret values keyed by query parameter name.
	Params = engine.Params
)
