package ferret

import "github.com/MontFerret/ferret/v2/pkg/engine"

type (
	// Module is an engine extension registered during Ferret initialization.
	Module = engine.Module

	// EngineInitHook runs during engine initialization.
	EngineInitHook = engine.EngineInitHook

	// EngineCloseHook runs during engine shutdown.
	EngineCloseHook = engine.EngineCloseHook

	// BeforeCompileHook runs before compilation starts.
	BeforeCompileHook = engine.BeforeCompileHook

	// AfterCompileHook runs after each compilation attempt.
	AfterCompileHook = engine.AfterCompileHook

	// PlanCloseHook runs when a compiled plan is closed.
	PlanCloseHook = engine.PlanCloseHook

	// BeforeRunHook runs before each session execution and may derive its context.
	BeforeRunHook = engine.BeforeRunHook

	// AfterRunHook runs after each session execution attempt.
	AfterRunHook = engine.AfterRunHook

	// SessionCloseHook runs when a session is closed.
	SessionCloseHook = engine.SessionCloseHook
)
