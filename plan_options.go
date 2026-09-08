package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/engine"
)

// PlanOption configures one compilation without changing the engine defaults.
type PlanOption = engine.PlanOption

// WithPlanOptimizationLevel overrides optimization for one compilation without
// changing the engine defaults. Debug compilation accepts only OptimizationNone.
func WithPlanOptimizationLevel(level OptimizationLevel) PlanOption {
	return engine.WithPlanOptimizationLevel(level)
}
