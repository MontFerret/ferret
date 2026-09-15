package runtime

import "github.com/MontFerret/ferret/v2/pkg/rnd"

// RandomDefault creates a standalone source and draws a Float in [0, 1).
// It does not participate in a Ferret Session's random sequence.
// Deprecated: Use rnd.New or rnd.NewSeed and retain the source across draws.
func RandomDefault() float64 {
	return rnd.New().Float64()
}

// Random creates a standalone source and applies the historical rounded
// rand(max, min) calculation. It does not participate in a Session's sequence.
// Deprecated: Use a retained rnd.Source and its LegacyFloat64 method for migration.
func Random(max float64, min float64) float64 {
	return rnd.New().LegacyFloat64(max, min)
}

// Random2 creates a standalone source and draws using historical mid/2 and mid*2
// bounds. It does not participate in a Ferret Session's random sequence.
// Deprecated: Use a retained rnd.Source and its LegacyFloat64 method for migration.
func Random2(mid float64) float64 {
	return rnd.New().LegacyFloat64(mid*2, mid/2)
}
