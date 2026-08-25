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

func newCompileOptions(setters []api.PlanOption) (compileOptions, error) {
	var opts compileOptions
	var err error

	for _, setter := range setters {
		if setter == nil {
			continue
		}

		if optionErr := setter(&opts); optionErr != nil {
			err = errors.Join(err, optionErr)
		}
	}

	return opts, err
}

func (o *compileOptions) SetOptimizationLevel(level api.OptimizationLevel) error {
	switch level {
	case api.OptimizationNone,
		api.OptimizationBasic,
		api.OptimizationFull,
		api.OptimizationAggressive:
		o.optimizationLevel = compiler.OptimizationLevel(level)
		o.hasOptimization = true

		return nil
	default:
		return fmt.Errorf("invalid optimization level: %d", level)
	}
}
