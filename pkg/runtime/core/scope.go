package core

import (
	"io"
)

type (
	CloseFunc func() error

	RootScope struct {
		closed      bool
		disposables []io.Closer
	}

	Scope struct {
		root   *RootScope
		parent *Scope
		vars   map[string]Value
	}
)

func NewRootScope() (*Scope, CloseFunc) {
	root := &RootScope{
		closed:      false,
		disposables: make([]io.Closer, 0, 10),
	}

	return newScope(root, nil), func() error {
		return root.Close()
	}
}

func (s *RootScope) AddDisposable(disposable io.Closer) {
	if s.closed {
		return
	}

	if disposable != nil {
		s.disposables = append(s.disposables, disposable)
	}
}

func (s *RootScope) Close() error {
	if s.closed {
		return Error(ErrInvalidOperation, "scope is already closed")
	}

	s.closed = true

	// close all values implemented io.Close
	for _, c := range s.disposables {
		c.Close()
	}

	return nil
}

func newScope(root *RootScope, parent *Scope) *Scope {
	return &Scope{
		root:   root,
		parent: parent,
		vars:   make(map[string]Value),
	}
}

func (s *Scope) SetVariable(name string, val Value) error {
	_, exists := s.vars[name]

	// it already has been declared in the current scope
	if exists {
		return Errorf(ErrNotUnique, "variable is already declared '%s'", name)
	}

	disposable, ok := val.(io.Closer)

	if ok {
		s.root.AddDisposable(disposable)
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

		return nil, Errorf(ErrNotFound, "variable '%s'", name)
	}

	return out, nil
}

func (s *Scope) Fork() *Scope {
	child := newScope(s.root, s)

	return child
}
