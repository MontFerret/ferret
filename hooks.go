package ferret

import (
	"context"
	"errors"
	"slices"
)

type (
	HookRegistrar interface {
		Engine() EngineHookRegistrar
		Session() SessionHookRegistrar
	}

	EngineHookRegistrar interface {
		OnInit(hook EngineInitHook)
		OnClose(hook EngineCloseHook)
	}

	SessionHookRegistrar interface {
		BeforeRun(hook BeforeRunHook)
		AfterRun(hook AfterRunHook)
		OnFailure(hook FailedRunHook)
		OnClose(hook SessionCloseHook)
	}
)

type (
	EngineInitHook func(ctx context.Context) error

	EngineCloseHook func() error

	BeforeRunHook func(ctx context.Context) (context.Context, error)

	FailedRunHook func(ctx context.Context, err error) error

	AfterRunHook func(ctx context.Context) error

	SessionCloseHook func() error
)

type (
	engineHooks interface {
		runInitHooks(ctx context.Context) error
		runCloseHooks() error
	}

	sessionHooks interface {
		runBeforeRunHooks(ctx context.Context) (context.Context, error)
		runAfterRunHooks(ctx context.Context) error
		runFailureHooks(ctx context.Context, err error) error
		runCloseHooks() error
	}

	hookRegistry struct {
		engine  *engineHookRegistry
		session *sessionHookRegistry
	}

	engineHookRegistry struct {
		onInit  []EngineInitHook
		onClose []EngineCloseHook
	}

	sessionHookRegistry struct {
		beforeRun []BeforeRunHook
		afterRun  []AfterRunHook
		onFailure []FailedRunHook
		onClose   []SessionCloseHook
	}
)

func newHookRegistry() *hookRegistry {
	return &hookRegistry{
		engine:  &engineHookRegistry{},
		session: &sessionHookRegistry{},
	}
}

func (hr *hookRegistry) Engine() EngineHookRegistrar {
	return hr.engine
}

func (hr *hookRegistry) Session() SessionHookRegistrar {
	return hr.session
}

func (hr *hookRegistry) clone() *hookRegistry {
	return &hookRegistry{
		engine:  hr.engine.clone(),
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

func (e *engineHookRegistry) runInitHooks(ctx context.Context) error {
	for _, hook := range e.onInit {
		if err := hook(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (e *engineHookRegistry) runCloseHooks() error {
	// Reverse the order of close hooks to ensure they run in LIFO order
	size := len(e.onClose)
	// We accumulate errors in a slice to return them all at once
	errs := make([]error, 0, size)

	for i := size - 1; i >= 0; i-- {
		if err := e.onClose[i](); err != nil {
			errs = append(errs, err)
			// We continue running the remaining hooks even if one fails to ensure all resources are attempted to be released
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

func (s *sessionHookRegistry) OnFailure(hook FailedRunHook) {
	if hook == nil {
		return
	}

	if s.onFailure == nil {
		s.onFailure = make([]FailedRunHook, 0, 1)
	}

	s.onFailure = append(s.onFailure, hook)
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
		onFailure: slices.Clone(s.onFailure),
		onClose:   slices.Clone(s.onClose),
	}
}

func (s *sessionHookRegistry) runBeforeRunHooks(ctx context.Context) (context.Context, error) {
	for _, hook := range s.beforeRun {
		var err error

		ctx, err = hook(ctx)

		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func (s *sessionHookRegistry) runAfterRunHooks(ctx context.Context) error {
	size := len(s.afterRun)

	// Reverse the order of after run hooks to ensure they run in LIFO order
	for i := size - 1; i >= 0; i-- {
		if err := s.afterRun[i](ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *sessionHookRegistry) runFailureHooks(ctx context.Context, err error) error {
	errs := make([]error, 0, len(s.onFailure))

	for _, hook := range s.onFailure {
		if err := hook(ctx, err); err != nil {
			errs = append(errs, err)
			// We continue running the remaining hooks even if one fails to ensure all failure handling is attempted
			continue
		}
	}

	return errors.Join(errs...)
}

func (s *sessionHookRegistry) runCloseHooks() error {
	size := len(s.onClose)
	// We accumulate errors in a slice to return them all at once
	errs := make([]error, 0, size)

	// Reverse the order of close hooks to ensure they run in LIFO order
	for i := size - 1; i >= 0; i-- {
		if err := s.onClose[i](); err != nil {
			errs = append(errs, err)
			// We continue running the remaining hooks even if one fails to ensure all resources are attempted to be released
			continue
		}
	}

	return errors.Join(errs...)
}
