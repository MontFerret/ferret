package ferret

import "github.com/MontFerret/ferret/v2/pkg/engine"

// OptimizationLevel controls the optimizer pipeline used during normal query compilation.
type OptimizationLevel = engine.OptimizationLevel

const (
	// OptimizationNone disables optimizer passes, including tail-call elimination.
	OptimizationNone OptimizationLevel = engine.OptimizationNone

	// OptimizationBasic enables the reduced pipeline, including safe tail-call
	// elimination, without register coalescing.
	OptimizationBasic OptimizationLevel = engine.OptimizationBasic

	// OptimizationFull is the default and enables the complete supported pipeline.
	OptimizationFull OptimizationLevel = engine.OptimizationFull
)
