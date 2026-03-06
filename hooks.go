package ferret

import (
	"context"
	"errors"
	"slices"
)

type (
	HookRegistrar interface {
		Engine() EngineHookRegistrar
		Plan() PlanHookRegistrar
		Session() SessionHookRegistrar
	}

	EngineHookRegistrar interface {
		OnInit(hook EngineInitHook)
		OnClose(hook EngineCloseHook)
	}

	PlanHookRegistrar interface {
		BeforeCompile(hook BeforeCompileHook)
		AfterCompile(hook AfterCompileHook)
		OnClose(hook PlanCloseHook)
	}

	SessionHookRegistrar interface {
		BeforeRun(hook BeforeRunHook)
		AfterRun(hook AfterRunHook)
		OnClose(hook SessionCloseHook)
	}
)

type (
	EngineInitHook func() error

	EngineCloseHook func() error

	BeforeCompileHook func(ctx context.Context) error

	AfterCompileHook func(ctx context.Context, err error) error

	PlanCloseHook func() error

	BeforeRunHook func(ctx context.Context) (context.Context, error)

	AfterRunHook func(ctx context.Context, err error) error

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

	// Reverse the order of close sessionHooks to ensure they run in LIFO order
	size := len(e.onClose)
	// We accumulate errors in a slice to return them all at once
	errs := make([]error, 0, size)

	for i := size - 1; i >= 0; i-- {
		if err := e.onClose[i](); err != nil {
			errs = append(errs, err)
			// We continue running the remaining sessionHooks even if one fails to ensure all resources are attempted to be released
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
			// We continue running the remaining sessionHooks even if one fails to ensure all post-compilation handling is attempted
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
	// We accumulate errors in a slice to return them all at once
	errs := make([]error, 0, size)

	// Reverse the order of close sessionHooks to ensure they run in LIFO order
	for i := size - 1; i >= 0; i-- {
		if hookErr := p.onClose[i](); hookErr != nil {
			errs = append(errs, hookErr)
			// We continue running the remaining sessionHooks even if one fails to ensure all resources are attempted to be released
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

	// Reverse the order of after run sessionHooks to ensure they run in LIFO order
	for i := size - 1; i >= 0; i-- {
		if hookErr := s.afterRun[i](ctx, err); hookErr != nil {
			errs = append(errs, hookErr)
			// We continue running the remaining sessionHooks even if one fails to ensure all post-run handling is attempted
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
	// We accumulate errors in a slice to return them all at once
	errs := make([]error, 0, size)

	// Reverse the order of close sessionHooks to ensure they run in LIFO order
	for i := size - 1; i >= 0; i-- {
		if err := s.onClose[i](); err != nil {
			errs = append(errs, err)
			// We continue running the remaining sessionHooks even if one fails to ensure all resources are attempted to be released
			continue
		}
	}

	return errors.Join(errs...)
}
