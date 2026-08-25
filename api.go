package ferret

import (
	"github.com/MontFerret/api"
	apisource "github.com/MontFerret/api/source"
)

type (
	// Source is a named FQL source file accepted by the embedding API.
	Source = apisource.File

	// Output is the encoded result returned from session or engine execution.
	Output = api.Output

	// PlanOption configures compilation of a reusable execution plan.
	PlanOption = api.PlanOption

	// PlanOptions is the target configured by a PlanOption.
	PlanOptions = api.PlanOptions

	// SessionOption configures a Session created from a Plan.
	SessionOption = api.SessionOption

	// SessionOptions is the target configured by a SessionOption.
	SessionOptions = api.SessionOptions

	// OptimizationLevel selects the compiler optimization pipeline.
	OptimizationLevel = api.OptimizationLevel
)

const (
	// OptimizationNone disables compiler optimizations.
	OptimizationNone = api.OptimizationNone
	// OptimizationBasic enables basic compiler optimizations.
	OptimizationBasic = api.OptimizationBasic
	// OptimizationFull enables the full compiler optimization pipeline.
	OptimizationFull = api.OptimizationFull
	// OptimizationAggressive enables aggressive compiler optimizations.
	OptimizationAggressive = api.OptimizationAggressive
)
