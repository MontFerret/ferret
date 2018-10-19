package compiler

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
)

type (
	scope struct {
		parent *scope
		vars   map[string]core.Type
	}
)

func newRootScope() *scope {
	return &scope{
		vars: make(map[string]core.Type),
	}
}

func newScope(parent *scope) *scope {
	s := newRootScope()
	s.parent = parent

	return s
}

func (s *scope) GetVariable(name string) (core.Type, error) {
	local, exists := s.vars[name]

	if exists {
		return local, nil
	}

	if s.parent != nil {
		parents, err := s.parent.GetVariable(name)

		if err != nil {
			return core.NoneType, err
		}

		return parents, nil
	}

	return core.NoneType, core.Error(ErrVariableNotFound, name)
}

func (s *scope) SetVariable(name string) error {
	_, exists := s.vars[name]

	if exists {
		return errors.Wrap(ErrVariableNotUnique, name)
	}

	// TODO: add type detection
	s.vars[name] = core.NoneType

	return nil
}

func (s *scope) RemoveVariable(name string) error {
	_, exists := s.vars[name]

	if !exists {
		return errors.Wrap(ErrVariableNotFound, name)
	}

	delete(s.vars, name)

	return nil
}

func (s *scope) Fork() *scope {
	return newScope(s)
}
