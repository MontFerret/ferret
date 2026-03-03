package data

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type aggregateState struct {
	sum       float64
	count     runtime.Int
	min       float64
	max       float64
	hasNumber bool
}

func validateAggregateSelectorIndex(idx int, count int) error {
	if idx < 0 || idx >= count {
		return runtime.Errorf(runtime.ErrInvalidArgument, "aggregate selector index out of range")
	}

	return nil
}

func updateAggregateState(state *aggregateState, kind bytecode.AggregateKind, value runtime.Value) {
	switch kind {
	case bytecode.AggregateCount:
		state.count++
	case bytecode.AggregateSum:
		if runtime.IsNumber(value) {
			state.sum += aggregateToFloat(value)
		}
	case bytecode.AggregateAverage:
		if runtime.IsNumber(value) {
			state.sum += aggregateToFloat(value)
			state.count++
		}
	case bytecode.AggregateMin:
		if runtime.IsNumber(value) {
			v := aggregateToFloat(value)
			if !state.hasNumber || v < state.min {
				state.min = v
			}
			state.hasNumber = true
		}
	case bytecode.AggregateMax:
		if runtime.IsNumber(value) {
			v := aggregateToFloat(value)
			if !state.hasNumber || v > state.max {
				state.max = v
			}
			state.hasNumber = true
		}
	}
}

func aggregateValueFor(state aggregateState, kind bytecode.AggregateKind) runtime.Value {
	switch kind {
	case bytecode.AggregateCount:
		return state.count
	case bytecode.AggregateSum:
		return runtime.NewFloat(state.sum)
	case bytecode.AggregateAverage:
		if state.count == 0 {
			return runtime.ZeroFloat
		}

		return runtime.NewFloat(state.sum / float64(state.count))
	case bytecode.AggregateMin:
		if !state.hasNumber {
			return runtime.None
		}

		return runtime.NewFloat(state.min)
	case bytecode.AggregateMax:
		if !state.hasNumber {
			return runtime.None
		}

		return runtime.NewFloat(state.max)
	default:
		return runtime.None
	}
}

func aggregateToFloat(arg runtime.Value) float64 {
	switch v := arg.(type) {
	case runtime.Float:
		return float64(v)
	case runtime.Int:
		return float64(v)
	default:
		return 0
	}
}
