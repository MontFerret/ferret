package host

import (
	"context"
	"errors"
	"slices"

	"github.com/MontFerret/ferret/v2/pkg/module"
)

type (
	// Hooks groups the mutable registrars and snapshots their callbacks for execution.
	Hooks struct {
		EngineHooks  *EngineHooks
		PlanHooks    *PlanHooks
		SessionHooks *SessionHooks
	}

	// EngineHooks stores initialization and shutdown callbacks.
	EngineHooks struct {
		onInit  []module.EngineInitHook
		onClose []module.EngineCloseHook
	}

	// PlanHooks stores compilation and plan shutdown callbacks.
	PlanHooks struct {
		beforeCompile []module.BeforeCompileHook
		afterCompile  []module.AfterCompileHook
		onClose       []module.PlanCloseHook
	}

	// SessionHooks stores execution and session shutdown callbacks.
	SessionHooks struct {
		beforeRun []module.BeforeRunHook
		afterRun  []module.AfterRunHook
		onClose   []module.SessionCloseHook
	}
)

// NewHooks creates empty registrars for each native lifecycle stage.
func NewHooks() *Hooks {
	return &Hooks{
		EngineHooks:  &EngineHooks{},
		PlanHooks:    &PlanHooks{},
		SessionHooks: &SessionHooks{},
	}
}

// Engine exposes engine hook registration.
func (hr *Hooks) Engine() module.EngineHookRegistrar {
	return hr.EngineHooks
}

// Plan exposes plan hook registration.
func (hr *Hooks) Plan() module.PlanHookRegistrar {
	return hr.PlanHooks
}

// Session exposes session hook registration.
func (hr *Hooks) Session() module.SessionHookRegistrar {
	return hr.SessionHooks
}

// Clone isolates the execution callbacks from subsequent module registration.
func (hr *Hooks) Clone() *Hooks {
	return &Hooks{
		EngineHooks:  hr.EngineHooks.clone(),
		PlanHooks:    hr.PlanHooks.clone(),
		SessionHooks: hr.SessionHooks.clone(),
	}
}

// OnInit appends a non-nil initialization callback.
func (e *EngineHooks) OnInit(hook module.EngineInitHook) {
	if hook == nil {
		return
	}

	if e.onInit == nil {
		e.onInit = make([]module.EngineInitHook, 0, 1)
	}

	e.onInit = append(e.onInit, hook)
}

// OnClose appends a non-nil shutdown callback.
func (e *EngineHooks) OnClose(hook module.EngineCloseHook) {
	if hook == nil {
		return
	}

	if e.onClose == nil {
		e.onClose = make([]module.EngineCloseHook, 0, 1)
	}

	e.onClose = append(e.onClose, hook)
}

// RunInit runs callbacks in registration order, stopping at the first error.
func (e *EngineHooks) RunInit() error {
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

// RunClose unwinds all shutdown callbacks and joins their failures.
func (e *EngineHooks) RunClose() error {
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

// BeforeCompile appends a non-nil before-compilation callback.
func (p *PlanHooks) BeforeCompile(hook module.BeforeCompileHook) {
	if hook == nil {
		return
	}

	if p.beforeCompile == nil {
		p.beforeCompile = make([]module.BeforeCompileHook, 0, 1)
	}

	p.beforeCompile = append(p.beforeCompile, hook)
}

// AfterCompile appends a non-nil after-compilation callback.
func (p *PlanHooks) AfterCompile(hook module.AfterCompileHook) {
	if hook == nil {
		return
	}

	if p.afterCompile == nil {
		p.afterCompile = make([]module.AfterCompileHook, 0, 1)
	}

	p.afterCompile = append(p.afterCompile, hook)
}

// OnClose appends a non-nil shutdown callback.
func (p *PlanHooks) OnClose(hook module.PlanCloseHook) {
	if hook == nil {
		return
	}

	if p.onClose == nil {
		p.onClose = make([]module.PlanCloseHook, 0, 1)
	}

	p.onClose = append(p.onClose, hook)
}

// RunBeforeCompile runs callbacks in registration order, stopping at the first error.
func (p *PlanHooks) RunBeforeCompile(ctx context.Context) error {
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

// RunAfterCompile unwinds callbacks with the same primary error and joins their failures.
func (p *PlanHooks) RunAfterCompile(ctx context.Context, err error) error {
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

// RunClose unwinds all shutdown callbacks and joins their failures.
func (p *PlanHooks) RunClose() error {
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

// BeforeRun appends a non-nil before-run callback.
func (s *SessionHooks) BeforeRun(hook module.BeforeRunHook) {
	if hook == nil {
		return
	}

	if s.beforeRun == nil {
		s.beforeRun = make([]module.BeforeRunHook, 0, 1)
	}

	s.beforeRun = append(s.beforeRun, hook)
}

// AfterRun appends a non-nil after-run callback.
func (s *SessionHooks) AfterRun(hook module.AfterRunHook) {
	if hook == nil {
		return
	}

	if s.afterRun == nil {
		s.afterRun = make([]module.AfterRunHook, 0, 1)
	}

	s.afterRun = append(s.afterRun, hook)
}

// OnClose appends a non-nil shutdown callback.
func (s *SessionHooks) OnClose(hook module.SessionCloseHook) {
	if hook == nil {
		return
	}

	if s.onClose == nil {
		s.onClose = make([]module.SessionCloseHook, 0, 1)
	}

	s.onClose = append(s.onClose, hook)
}

// RunBeforeRun chains callback contexts in registration order, stopping at the first error.
func (s *SessionHooks) RunBeforeRun(ctx context.Context) (context.Context, error) {
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

// RunAfterRun unwinds callbacks with the same primary error and joins their failures.
func (s *SessionHooks) RunAfterRun(ctx context.Context, err error) error {
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

// RunClose unwinds all shutdown callbacks and joins their failures.
func (s *SessionHooks) RunClose() error {
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

func (e *EngineHooks) clone() *EngineHooks {
	return &EngineHooks{
		onInit:  slices.Clone(e.onInit),
		onClose: slices.Clone(e.onClose),
	}
}

func (p *PlanHooks) clone() *PlanHooks {
	return &PlanHooks{
		beforeCompile: slices.Clone(p.beforeCompile),
		afterCompile:  slices.Clone(p.afterCompile),
		onClose:       slices.Clone(p.onClose),
	}
}

func (s *SessionHooks) clone() *SessionHooks {
	return &SessionHooks{
		beforeRun: slices.Clone(s.beforeRun),
		afterRun:  slices.Clone(s.afterRun),
		onClose:   slices.Clone(s.onClose),
	}
}
