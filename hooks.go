package ferret

import (
	"context"
	"errors"
	"slices"
)

type (
	// HookRegistrar provides access to lifecycle hook registrars for engine, plan, and session stages.
	HookRegistrar interface {
		// Engine returns the registrar for engine lifecycle hooks.
		Engine() EngineHookRegistrar
		// Plan returns the registrar for plan lifecycle hooks.
		Plan() PlanHookRegistrar
		// Session returns the registrar for session lifecycle hooks.
		Session() SessionHookRegistrar
	}

	// EngineHookRegistrar registers hooks for engine initialization and shutdown.
	EngineHookRegistrar interface {
		// OnInit registers a hook executed in FIFO order during engine initialization.
		// A nil hook is ignored.
		OnInit(hook EngineInitHook)
		// OnClose registers a hook executed in LIFO order when the engine is closed.
		// A nil hook is ignored.
		OnClose(hook EngineCloseHook)
	}

	// PlanHookRegistrar registers hooks for compilation and plan shutdown.
	PlanHookRegistrar interface {
		// BeforeCompile registers a hook executed in FIFO order before compilation starts.
		// A nil hook is ignored.
		BeforeCompile(hook BeforeCompileHook)
		// AfterCompile registers a hook executed in LIFO order after each compilation attempt.
		// It receives the compilation error (if any), and a nil hook is ignored.
		AfterCompile(hook AfterCompileHook)
		// OnClose registers a hook executed in LIFO order when a plan is closed.
		// A nil hook is ignored.
		OnClose(hook PlanCloseHook)
	}

	// SessionHookRegistrar registers hooks for session execution and shutdown.
	SessionHookRegistrar interface {
		// BeforeRun registers a hook executed in FIFO order before each session run.
		// Hooks can replace the context passed to subsequent hooks and VM execution.
		// A nil hook is ignored.
		BeforeRun(hook BeforeRunHook)
		// AfterRun registers a hook executed in LIFO order after each run attempt.
		// It receives the run error (if any), and a nil hook is ignored.
		AfterRun(hook AfterRunHook)
		// OnClose registers a hook executed in LIFO order when a session is closed.
		// A nil hook is ignored.
		OnClose(hook SessionCloseHook)
	}
)

type (
	// EngineInitHook runs during engine initialization.
	// Returning an error stops initialization immediately.
	EngineInitHook func() error

	// EngineCloseHook runs during engine shutdown.
	// Close hooks are executed in LIFO order and their errors are aggregated.
	EngineCloseHook func() error

	// BeforeCompileHook runs before compilation starts.
	// Hooks run in FIFO order and stop on the first error.
	BeforeCompileHook func(ctx context.Context) error

	// AfterCompileHook runs after each compilation attempt.
	// Hooks run in LIFO order and receive the compilation error (if any).
	AfterCompileHook func(ctx context.Context, err error) error

	// PlanCloseHook runs when a plan is closed.
	// Close hooks are executed in LIFO order and their errors are aggregated.
	PlanCloseHook func() error

	// BeforeRunHook runs before each session run.
	// It can return a derived context for subsequent hooks and VM execution.
	BeforeRunHook func(ctx context.Context) (context.Context, error)

	// AfterRunHook runs after each session run attempt.
	// Hooks run in LIFO order, receive the run error (if any), and aggregate hook errors.
	AfterRunHook func(ctx context.Context, err error) error

	// SessionCloseHook runs when a session is closed.
	// Close hooks are executed in LIFO order and their errors are aggregated.
	SessionCloseHook func() error
)

type (
	engineHooks interface {
		runInitHooks() error
		runCloseHooks() error
	}

	planHooks interface {
		runBeforeCompileHooks(ctx context.Context) error
		runAfterCompileHooks(ctx context.Context, err error) error
		runCloseHooks() error
	}

	sessionHooks interface {
		runBeforeRunHooks(ctx context.Context) (context.Context, error)
		runAfterRunHooks(ctx context.Context, err error) error
		runCloseHooks() error
	}

	hookRegistry struct {
		engine  *engineHookRegistry
		plan    *planHookRegistry
		session *sessionHookRegistry
	}

	engineHookRegistry struct {
		onInit  []EngineInitHook
		onClose []EngineCloseHook
	}

	planHookRegistry struct {
		beforeCompile []BeforeCompileHook
		afterCompile  []AfterCompileHook
		onClose       []PlanCloseHook
	}

	sessionHookRegistry struct {
		beforeRun []BeforeRunHook
		afterRun  []AfterRunHook
		onClose   []SessionCloseHook
	}
)

func newHookRegistry() *hookRegistry {
	return &hookRegistry{
		engine:  &engineHookRegistry{},
		plan:    &planHookRegistry{},
		session: &sessionHookRegistry{},
	}
}

func (hr *hookRegistry) Engine() EngineHookRegistrar {
	return hr.engine
}

func (hr *hookRegistry) Plan() PlanHookRegistrar {
	return hr.plan
}

func (hr *hookRegistry) Session() SessionHookRegistrar {
	return hr.session
}

func (hr *hookRegistry) clone() *hookRegistry {
	return &hookRegistry{
		engine:  hr.engine.clone(),
		plan:    hr.plan.clone(),
		session: hr.session.clone(),
	}
}

func (e *engineHookRegistry) OnInit(hook EngineInitHook) {
	if hook == nil {
		return
	}

	if e.onInit == nil {
		e.onInit = make([]EngineInitHook, 0, 1)
	}

	e.onInit = append(e.onInit, hook)
}

func (e *engineHookRegistry) OnClose(hook EngineCloseHook) {
	if hook == nil {
		return
	}

	if e.onClose == nil {
		e.onClose = make([]EngineCloseHook, 0, 1)
	}

	e.onClose = append(e.onClose, hook)
}

func (e *engineHookRegistry) clone() *engineHookRegistry {
	return &engineHookRegistry{
		onInit:  slices.Clone(e.onInit),
		onClose: slices.Clone(e.onClose),
	}
}

func (e *engineHookRegistry) runInitHooks() error {
	if len(e.onInit) == 0 {
		return nil
	}

	for _, hook := range e.onInit {
		if err := hook(); err != nil {
			return err
		}
	}

	return nil
}

func (e *engineHookRegistry) runCloseHooks() error {
	if len(e.onClose) == 0 {
		return nil
	}

	size := len(e.onClose)
	errs := make([]error, 0, size)

	// Close hooks run in reverse registration order (LIFO).
	for i := size - 1; i >= 0; i-- {
		if err := e.onClose[i](); err != nil {
			errs = append(errs, err)
			// Continue so remaining hooks can attempt resource cleanup.
			continue
		}
	}

	return errors.Join(errs...)
}

func (p *planHookRegistry) BeforeCompile(hook BeforeCompileHook) {
	if hook == nil {
		return
	}

	if p.beforeCompile == nil {
		p.beforeCompile = make([]BeforeCompileHook, 0, 1)
	}

	p.beforeCompile = append(p.beforeCompile, hook)
}

func (p *planHookRegistry) AfterCompile(hook AfterCompileHook) {
	if hook == nil {
		return
	}

	if p.afterCompile == nil {
		p.afterCompile = make([]AfterCompileHook, 0, 1)
	}

	p.afterCompile = append(p.afterCompile, hook)
}

func (p *planHookRegistry) OnClose(hook PlanCloseHook) {
	if hook == nil {
		return
	}

	if p.onClose == nil {
		p.onClose = make([]PlanCloseHook, 0, 1)
	}

	p.onClose = append(p.onClose, hook)
}

func (p *planHookRegistry) clone() *planHookRegistry {
	return &planHookRegistry{
		beforeCompile: slices.Clone(p.beforeCompile),
		afterCompile:  slices.Clone(p.afterCompile),
		onClose:       slices.Clone(p.onClose),
	}
}

func (p *planHookRegistry) runBeforeCompileHooks(ctx context.Context) error {
	if len(p.beforeCompile) == 0 {
		return nil
	}

	for _, hook := range p.beforeCompile {
		if err := hook(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (p *planHookRegistry) runAfterCompileHooks(ctx context.Context, err error) error {
	if len(p.afterCompile) == 0 {
		return nil
	}

	size := len(p.afterCompile)
	errs := make([]error, 0, size)

	for i := size - 1; i >= 0; i-- {
		if hookErr := p.afterCompile[i](ctx, err); hookErr != nil {
			errs = append(errs, hookErr)
			// Continue so remaining hooks can run their post-compilation handling.
			continue
		}
	}

	return errors.Join(errs...)
}

func (p *planHookRegistry) runCloseHooks() error {
	if len(p.onClose) == 0 {
		return nil
	}

	size := len(p.onClose)
	errs := make([]error, 0, size)

	// Close hooks run in reverse registration order (LIFO).
	for i := size - 1; i >= 0; i-- {
		if hookErr := p.onClose[i](); hookErr != nil {
			errs = append(errs, hookErr)
			// Continue so remaining hooks can attempt resource cleanup.
			continue
		}
	}

	return errors.Join(errs...)
}

func (s *sessionHookRegistry) BeforeRun(hook BeforeRunHook) {
	if hook == nil {
		return
	}

	if s.beforeRun == nil {
		s.beforeRun = make([]BeforeRunHook, 0, 1)
	}

	s.beforeRun = append(s.beforeRun, hook)
}

func (s *sessionHookRegistry) AfterRun(hook AfterRunHook) {
	if hook == nil {
		return
	}

	if s.afterRun == nil {
		s.afterRun = make([]AfterRunHook, 0, 1)
	}

	s.afterRun = append(s.afterRun, hook)
}

func (s *sessionHookRegistry) OnClose(hook SessionCloseHook) {
	if hook == nil {
		return
	}

	if s.onClose == nil {
		s.onClose = make([]SessionCloseHook, 0, 1)
	}

	s.onClose = append(s.onClose, hook)
}

func (s *sessionHookRegistry) clone() *sessionHookRegistry {
	return &sessionHookRegistry{
		beforeRun: slices.Clone(s.beforeRun),
		afterRun:  slices.Clone(s.afterRun),
		onClose:   slices.Clone(s.onClose),
	}
}

func (s *sessionHookRegistry) runBeforeRunHooks(ctx context.Context) (context.Context, error) {
	if len(s.beforeRun) == 0 {
		return ctx, nil
	}

	for _, hook := range s.beforeRun {
		var err error

		// Each hook receives the context returned by the previous hook.
		ctx, err = hook(ctx)

		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func (s *sessionHookRegistry) runAfterRunHooks(ctx context.Context, err error) error {
	if len(s.afterRun) == 0 {
		return nil
	}

	size := len(s.afterRun)
	errs := make([]error, 0, size)

	// After-run hooks execute in reverse registration order (LIFO).
	for i := size - 1; i >= 0; i-- {
		if hookErr := s.afterRun[i](ctx, err); hookErr != nil {
			errs = append(errs, hookErr)
			// Continue so remaining hooks can run their post-run handling.
			continue
		}
	}

	return errors.Join(errs...)
}

func (s *sessionHookRegistry) runCloseHooks() error {
	if len(s.onClose) == 0 {
		return nil
	}

	size := len(s.onClose)
	errs := make([]error, 0, size)

	// Close hooks run in reverse registration order (LIFO).
	for i := size - 1; i >= 0; i-- {
		if err := s.onClose[i](); err != nil {
			errs = append(errs, err)
			// Continue so remaining hooks can attempt resource cleanup.
			continue
		}
	}

	return errors.Join(errs...)
}
