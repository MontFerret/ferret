package ferret

import "github.com/MontFerret/ferret/v2/pkg/module"

type (
	// Module is an engine extension registered during Ferret initialization.
	Module = module.Module

	// EngineInitHook runs during engine initialization.
	EngineInitHook = module.EngineInitHook

	// EngineCloseHook runs during engine shutdown.
	EngineCloseHook = module.EngineCloseHook

	// BeforeCompileHook runs before compilation starts.
	BeforeCompileHook = module.BeforeCompileHook

	// AfterCompileHook runs after each compilation attempt.
	AfterCompileHook = module.AfterCompileHook

	// PlanCloseHook runs when a compiled plan is closed.
	PlanCloseHook = module.PlanCloseHook

	// BeforeRunHook runs before each session execution and may derive its context.
	BeforeRunHook = module.BeforeRunHook

	// AfterRunHook runs after each session execution attempt.
	AfterRunHook = module.AfterRunHook

	// SessionCloseHook runs when a session is closed.
	SessionCloseHook = module.SessionCloseHook
)
