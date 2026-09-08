// Package resource manages ownership and cleanup of engine-scoped host resources.
package resource

import (
	"errors"
	"fmt"
	"slices"
)

type (
	// Name identifies a resource slot whose registrations replace one another.
	Name string

	// Manager retains lifecycle metadata, never the services themselves.
	// Calls must be serialized by its owner. Cleanup callbacks must not call back
	// into the manager. The zero value is ready to use.
	Manager struct {
		closeErr error
		entries  []entry
		closed   bool
	}

	entry struct {
		close func() error
		name  Name
	}
)

const (
	// FileSystem identifies the engine's configured filesystem.
	FileSystem Name = "filesystem"
	// Network identifies the engine's configured network service.
	Network Name = "network"
)

// NewManager creates an empty, open ownership scope.
func NewManager() *Manager {
	return &Manager{}
}

// Own takes responsibility for cleanup and retires any previous registration.
// A retirement error does not undo the replacement: the new callback remains
// owned and the previous callback is never retried. A closed manager or nil
// callback is rejected without accepting ownership or changing existing entries.
func (m *Manager) Own(name Name, close func() error) error {
	if m.closed {
		return errors.New("resource manager is closed")
	}

	if close == nil {
		return errors.New("owned resource requires a cleanup callback")
	}

	return m.replace(entry{name: name, close: close})
}

// Borrow records caller ownership and retires any previous registration.
// The borrowed resource is never closed, even if retiring its predecessor fails.
// A closed manager rejects the registration without changing ownership.
func (m *Manager) Borrow(name Name) error {
	if m.closed {
		return errors.New("resource manager is closed")
	}

	return m.replace(entry{name: name})
}

// Close attempts all owned callbacks in reverse acquisition/replacement order.
// Cleanup errors are joined and retained for repeated calls. No callback is
// retried, and no new registrations are accepted after closure.
func (m *Manager) Close() error {
	if m.closed {
		return m.closeErr
	}

	m.closed = true
	entries := m.entries
	m.entries = nil

	var errs []error

	for i := len(entries) - 1; i >= 0; i-- {
		if err := entries[i].cleanup(); err != nil {
			errs = append(errs, err)
		}
	}

	m.closeErr = errors.Join(errs...)

	return m.closeErr
}

func (m *Manager) replace(next entry) error {
	var previous entry

	for i, current := range m.entries {
		if current.name == next.name {
			previous = current
			m.entries = slices.Delete(m.entries, i, i+1)

			break
		}
	}

	m.entries = append(m.entries, next)

	return previous.cleanup()
}

func (e entry) cleanup() error {
	if e.close == nil {
		return nil
	}

	if err := e.close(); err != nil {
		return fmt.Errorf("close %s: %w", e.name, err)
	}

	return nil
}
