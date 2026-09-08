package ferret

import (
	"fmt"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

type (
	// PlanOption configures one compilation without changing the engine defaults.
	PlanOption = gooptions.Option[planConfig]

	planConfig struct {
		level compiler.OptimizationLevel
		debug bool
	}
)

func newPlanConfig(level compiler.OptimizationLevel, debug bool, setters []PlanOption) (planConfig, error) {
	if len(setters) == 0 {
		return planConfig{level: level, debug: debug}, nil
	}

	return gooptions.ApplyTo[planConfig](planConfig{level: level, debug: debug}, setters...)
}

// WithPlanOptimizationLevel overrides optimization for one compilation without
// changing the engine defaults. Debug compilation accepts only OptimizationNone.
func WithPlanOptimizationLevel(level OptimizationLevel) PlanOption {
	return func(cfg *planConfig) error {
		if cfg.debug && level != OptimizationNone {
			return fmt.Errorf("debug compilation requires optimization none")
		}

		selected, err := level.compilerLevel()
		if err != nil {
			return err
		}

		cfg.level = selected

		return nil
	}
}
