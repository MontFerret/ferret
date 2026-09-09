package ferret

import "github.com/MontFerret/ferret/v2/pkg/engine"

type (
	// Engine compiles and runs queries using its configured modules and host services.
	Engine = engine.Engine

	// Plan is a reusable compiled query. Callers close its sessions before closing the plan.
	Plan = engine.Plan

	// Session executes a plan with per-execution options and must be closed by its caller.
	Session = engine.Session
)

// New constructs an Engine with the supplied options. The caller owns the engine
// and closes it after closing any plans and sessions it created.
func New(opts ...Option) (*Engine, error) {
	return engine.New(opts...)
}
