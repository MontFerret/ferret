package compiler

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	globalScope struct {
		params map[string]struct{}
	}

	scope struct {
		global *globalScope
		parent *scope
		vars   map[string]struct{}
	}
)

func newGlobalScope() *globalScope {
	return &globalScope{
		params: map[string]struct{}{},
	}
}

func newRootScope(global *globalScope) *scope {
	return &scope{
		global: global,
		vars:   make(map[string]struct{}),
	}
}

func newScope(parent *scope) *scope {
	s := newRootScope(parent.global)
	s.parent = parent

	return s
}

func (s *scope) AddParam(name string) {
	s.global.params[name] = struct{}{}
}

func (s *scope) HasVariable(name string) bool {
	_, exists := s.vars[name]

	if exists {
		return true
	}

	if s.parent != nil {
		return s.parent.HasVariable(name)
	}

	return false
}

func (s *scope) SetVariable(name string) error {
	_, exists := s.vars[name]

	if exists {
		return core.Error(ErrVariableNotUnique, name)
	}

	// TODO: add type detection
	s.vars[name] = struct{}{}

	return nil
}

func (s *scope) RemoveVariable(name string) error {
	_, exists := s.vars[name]

	if !exists {
		return core.Error(ErrVariableNotFound, name)
	}

	delete(s.vars, name)

	return nil
}

func (s *scope) ClearVariables() {
	s.vars = make(map[string]struct{})
}

func (s *scope) Fork() *scope {
	return newScope(s)
}
