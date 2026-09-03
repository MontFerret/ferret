package ferret

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

// OptimizationLevel controls the optimizer pipeline used during normal query compilation.
type OptimizationLevel int

const (
	// OptimizationNone disables optimizer passes.
	OptimizationNone OptimizationLevel = iota
	// OptimizationBasic enables the reduced pipeline without register coalescing.
	OptimizationBasic
	// OptimizationFull is the default and enables the complete supported pipeline.
	OptimizationFull
)

// String returns the semantic name of the optimization level.
func (l OptimizationLevel) String() string {
	switch l {
	case OptimizationNone:
		return "none"
	case OptimizationBasic:
		return "basic"
	case OptimizationFull:
		return "full"
	default:
		return "unknown"
	}
}

func (l OptimizationLevel) compilerLevel() (compiler.OptimizationLevel, error) {
	switch l {
	case OptimizationNone:
		return compiler.None, nil
	case OptimizationBasic:
		return compiler.Basic, nil
	case OptimizationFull:
		return compiler.Full, nil
	default:
		return 0, fmt.Errorf("unsupported optimization level %d", l)
	}
}
