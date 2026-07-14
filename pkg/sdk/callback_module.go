package sdk

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/module"
)

type callbackModule struct {
	register      func(module.Bootstrap) error
	validationErr error
	name          string
}

// NewModule creates a module backed by a registration callback.
// Invalid names and nil callbacks are reported when the engine registers the module.
func NewModule(name string, register func(module.Bootstrap) error) module.Module {
	m := &callbackModule{
		name:     name,
		register: register,
	}

	if strings.TrimSpace(name) == "" {
		m.validationErr = fmt.Errorf("module name cannot be empty")
	} else if register == nil {
		m.validationErr = fmt.Errorf("module %q: register callback cannot be nil", name)
	}

	return m
}

func (m *callbackModule) Name() string {
	if m == nil {
		return ""
	}

	return m.name
}

func (m *callbackModule) Register(bootstrap module.Bootstrap) error {
	if m == nil {
		return fmt.Errorf("module cannot be nil")
	}

	if m.validationErr != nil {
		return m.validationErr
	}

	if err := m.register(bootstrap); err != nil {
		return fmt.Errorf("module %q: register: %w", m.name, err)
	}

	return nil
}
