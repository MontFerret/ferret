package universal

import (
	"errors"
	"fmt"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type planOptions struct {
	native []engine.PlanOption
}

var _ api.PlanOptions = (*planOptions)(nil)

func (o *planOptions) SetOptimizationLevel(level api.OptimizationLevel) error {
	var native engine.OptimizationLevel
	switch level {
	case api.OptimizationNone:
		native = engine.OptimizationNone
	case api.OptimizationBasic:
		native = engine.OptimizationBasic
	case api.OptimizationFull:
		native = engine.OptimizationFull
	default:
		return fmt.Errorf("unsupported optimization level: %d", level)
	}

	o.native = append(o.native, engine.WithPlanOptimizationLevel(native))

	return nil
}

func newPlanOptions(setters []api.PlanOption) (*planOptions, error) {
	opts := &planOptions{}
	var failures []error
	for _, setter := range setters {
		if setter == nil {
			continue
		}

		if err := setter(opts); err != nil {
			failures = append(failures, err)
		}
	}

	return opts, errors.Join(failures...)
}
