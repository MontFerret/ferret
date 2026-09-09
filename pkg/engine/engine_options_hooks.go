package engine

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/module"
)

// WithEngineInitHook returns an Option that registers a hook to execute during engine initialization.
// It returns an error if hook is nil.
func WithEngineInitHook(hook module.EngineInitHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("engine init hook is nil")
		}

		opts.hooks.EngineHooks.OnInit(hook)

		return nil
	}
}

// WithEngineCloseHook returns an Option that registers a hook to execute when the engine is closed.
// It returns an error if hook is nil.
func WithEngineCloseHook(hook module.EngineCloseHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("engine close hook is nil")
		}

		opts.hooks.EngineHooks.OnClose(hook)

		return nil
	}
}

// WithBeforeCompileHook returns an Option that registers a hook to execute before each compilation attempt.
// It returns an error if hook is nil.
func WithBeforeCompileHook(hook module.BeforeCompileHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("before compile hook is nil")
		}

		opts.hooks.PlanHooks.BeforeCompile(hook)

		return nil
	}
}

// WithAfterCompileHook returns an Option that registers a hook to execute after each compilation attempt.
// The hook receives the compilation error (if any). It returns an error if hook is nil.
func WithAfterCompileHook(hook module.AfterCompileHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("after compile hook is nil")
		}

		opts.hooks.PlanHooks.AfterCompile(hook)

		return nil
	}
}

// WithPlanCloseHook returns an Option that registers a hook to execute when a plan is closed.
// It returns an error if hook is nil.
func WithPlanCloseHook(hook module.PlanCloseHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("plan close hook is nil")
		}

		opts.hooks.PlanHooks.OnClose(hook)

		return nil
	}
}

// WithBeforeRunHook returns an Option that registers a hook to execute before each session run.
// The hook can replace the context used by subsequent hooks and VM execution.
// It returns an error if hook is nil.
func WithBeforeRunHook(hook module.BeforeRunHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("before run hook is nil")
		}

		opts.hooks.SessionHooks.BeforeRun(hook)

		return nil
	}
}

// WithAfterRunHook returns an Option that registers a hook to execute after each session run attempt.
// The hook receives the run error (if any). It returns an error if hook is nil.
func WithAfterRunHook(hook module.AfterRunHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("after run hook is nil")
		}

		opts.hooks.SessionHooks.AfterRun(hook)

		return nil
	}
}

// WithSessionCloseHook returns an Option that registers a hook to execute when a session is closed.
// It returns an error if hook is nil.
func WithSessionCloseHook(hook module.SessionCloseHook) Option {
	return func(opts *config) error {
		if hook == nil {
			return fmt.Errorf("session close hook is nil")
		}

		opts.hooks.SessionHooks.OnClose(hook)

		return nil
	}
}
