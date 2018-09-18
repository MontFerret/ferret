package core

import (
	"github.com/pkg/errors"
	"io"
)

type (
	CloseFunc func()

	Scope struct {
		closed   bool
		vars     map[string]Value
		parent   *Scope
		children []*Scope
	}
)

func NewRootScope() (*Scope, CloseFunc) {
	scope := NewScope(nil)

	return scope, func() {
		scope.close()
	}
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		closed:   false,
		vars:     make(map[string]Value),
		parent:   parent,
		children: make([]*Scope, 0, 5),
	}
}

func (s *Scope) SetVariable(name string, val Value) error {
	_, exists := s.vars[name]

	// it already has been declared in the current scope
	if exists {
		return errors.Wrapf(ErrNotUnique, "variable is already declared '%s'", name)
	}

	s.vars[name] = val

	return nil
}

func (s *Scope) HasVariable(name string) bool {
	_, exists := s.vars[name]

	// does not exist in the current scope
	// try to find in a parent scope
	if !exists {
		if s.parent != nil {
			return s.parent.HasVariable(name)
		}
	}

	return exists
}

func (s *Scope) GetVariable(name string) (Value, error) {
	out, exists := s.vars[name]

	// does not exist in the current scope
	// try to find in the parent scope
	if !exists {
		if s.parent != nil {
			return s.parent.GetVariable(name)
		}

		return nil, errors.Wrapf(ErrNotFound, "variable '%s'", name)
	}

	return out, nil
}

func (s *Scope) Fork() *Scope {
	child := NewScope(s)

	s.children = append(s.children, child)

	return child
}

func (s *Scope) close() error {
	if s.closed {
		return errors.Wrap(ErrInvalidOperation, "scope is already closed")
	}

	s.closed = true

	// close all active child scopes
	for _, c := range s.children {
		c.close()
	}

	// do clean up
	// if some of the variables implements io.Closer interface
	// we need to close them
	for _, v := range s.vars {
		closer, ok := v.(io.Closer)

		if ok {
			closer.Close()
		}
	}

	s.children = nil
	s.vars = nil

	return nil
}
