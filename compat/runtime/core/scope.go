package core

import (
	"fmt"
	"io"
)

// RootScope manages the lifecycle of disposable resources for an execution scope.
// Provided for compilation compatibility with v1 code.
type RootScope struct {
	closed      bool
	disposables []io.Closer
}

// AddDisposable registers a Closer to be invoked when the scope is closed.
func (r *RootScope) AddDisposable(d io.Closer) {
	if r.closed || d == nil {
		return
	}

	r.disposables = append(r.disposables, d)
}

// Close closes all registered disposables in reverse order.
func (r *RootScope) Close() error {
	if r.closed {
		return fmt.Errorf("scope is already closed")
	}

	r.closed = true

	var errs []error

	for i := len(r.disposables) - 1; i >= 0; i-- {
		if err := r.disposables[i].Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("scope close errors: %v", errs)
}

// Scope is a variable scope used during FQL expression execution.
// Provided for compilation compatibility with v1 code.
type Scope struct {
	root   *RootScope
	parent *Scope
	vars   map[string]Value
}

func newScope(root *RootScope, parent *Scope) *Scope {
	return &Scope{
		root:   root,
		parent: parent,
		vars:   make(map[string]Value),
	}
}

// SetVariable sets a named variable in this scope.
func (s *Scope) SetVariable(name string, value Value) {
	s.vars[name] = value
}

// GetVariable looks up a variable by name, walking parent scopes.
func (s *Scope) GetVariable(name string) (Value, bool) {
	if v, ok := s.vars[name]; ok {
		return v, true
	}

	if s.parent != nil {
		return s.parent.GetVariable(name)
	}

	return nil, false
}

// Fork creates a child scope.
func (s *Scope) Fork() *Scope {
	return newScope(s.root, s)
}

// Root returns the root scope.
func (s *Scope) Root() *RootScope {
	return s.root
}
