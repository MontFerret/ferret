package universal

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type planOptions struct {
	native []engine.PlanOption
	debug  bool
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

	if o.debug && native != engine.OptimizationNone {
		return fmt.Errorf("debug compilation requires optimization none")
	}

	o.native = append(o.native, engine.WithPlanOptimizationLevel(native))

	return nil
}

func newPlanOptions(ctx context.Context, debug bool, setters []api.PlanOption) (*planOptions, error) {
	if err := checkOptionContext(ctx); err != nil {
		return nil, err
	}

	opts := &planOptions{debug: debug}
	var failures []error
	for _, setter := range setters {
		if setter != nil {
			if err := setter(opts); err != nil {
				failures = append(failures, err)
			}
		}
	}

	if err := ctx.Err(); err != nil {
		failures = append(failures, err)
	}

	if err := errors.Join(failures...); err != nil {
		return nil, err
	}

	return opts, nil
}
