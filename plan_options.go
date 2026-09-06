package ferret

import (
	"errors"
	"fmt"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

type (
	// PlanOption configures one compilation without changing the engine defaults.
	PlanOption  = api.PlanOption
	planOptions struct {
		level compiler.OptimizationLevel
		debug bool
	}
)

func newPlanOptions(level compiler.OptimizationLevel, debug bool, setters []PlanOption) (planOptions, error) {
	if len(setters) == 0 {
		return planOptions{level: level, debug: debug}, nil
	}

	opts := planOptions{level: level, debug: debug}
	var failures []error
	for _, setter := range setters {
		if setter != nil {
			if err := setter(&opts); err != nil {
				failures = append(failures, err)
			}
		}
	}

	return opts, errors.Join(failures...)
}

func (o *planOptions) SetOptimizationLevel(level api.OptimizationLevel) error {
	if o.debug && level != api.OptimizationNone {
		return fmt.Errorf("debug compilation requires optimization none")
	}

	switch level {
	case api.OptimizationNone:
		o.level = compiler.None
	case api.OptimizationBasic:
		o.level = compiler.Basic
	case api.OptimizationFull:
		o.level = compiler.Full
	default:
		return fmt.Errorf("unsupported optimization level %d", level)
	}

	return nil
}
