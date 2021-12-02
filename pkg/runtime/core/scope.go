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

	return newScope(root, nil), root.Close
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

	var errors []error

	// close all values implemented io.Close
	for _, c := range s.disposables {
		if err := c.Close(); err != nil {
			if errors == nil {
				errors = make([]error, 0, len(s.disposables))
			}

			errors = append(errors, err)
		}
	}

	if errors == nil {
		return nil
	}

	return Errors(errors...)
}

func newScope(root *RootScope, parent *Scope) *Scope {
	return &Scope{
		root:   root,
		parent: parent,
		vars:   make(map[string]Value),
	}
}

func (s *Scope) SetVariable(name string, val Value) error {
	if name != IgnorableVariable {
		_, exists := s.vars[name]

		// it already has been declared in the current scope
		if exists {
			return Errorf(ErrNotUnique, "variable is already declared: '%s'", name)
		}

		s.vars[name] = val
	}

	// we still want to make sure that nothing than needs to be closed leaks out
	disposable, ok := val.(io.Closer)

	if ok {
		s.root.AddDisposable(disposable)
	}

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

		return nil, Errorf(ErrNotFound, "variable: '%s'", name)
	}

	return out, nil
}

func (s *Scope) MustGetVariable(name string) Value {
	out, err := s.GetVariable(name)

	if err != nil {
		panic(err)
	}

	return out
}

func (s *Scope) UpdateVariable(name string, val Value) error {
	_, exists := s.vars[name]

	if !exists {
		return Errorf(ErrNotFound, "variable: '%s'", name)
	}

	delete(s.vars, name)

	return s.SetVariable(name, val)
}

func (s *Scope) Fork() *Scope {
	child := newScope(s.root, s)

	return child
}
