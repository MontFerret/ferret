package ferret

import (
	"errors"
	"fmt"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

type compileOptions struct {
	optimizationLevel compiler.OptimizationLevel
	hasOptimization   bool
}

// WithOptimizationLevel sets the optimization level for the execution plan.
func WithOptimizationLevel(level OptimizationLevel) PlanOption {
	return api.WithOptimizationLevel(level)
}

func newCompileOptions(setters []PlanOption) (compileOptions, error) {
	var opts compileOptions
	var failures []error

	for _, setter := range setters {
		if setter == nil {
			continue
		}

		if optionErr := setter(&opts); optionErr != nil {
			failures = append(failures, optionErr)
		}
	}

	return opts, errors.Join(failures...)
}

func (o *compileOptions) SetOptimizationLevel(level OptimizationLevel) error {
	switch level {
	case OptimizationNone,
		OptimizationBasic,
		OptimizationFull,
		OptimizationAggressive:
		o.optimizationLevel = compiler.OptimizationLevel(level)
		o.hasOptimization = true

		return nil
	default:
		return fmt.Errorf("invalid optimization level: %d", level)
	}
}
