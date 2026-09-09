package engine

import (
	"errors"
	"fmt"
)

func (e *Engine) close() error {
	var hookErr error

	if e.hooks.EngineHooks != nil {
		hookErr = e.hooks.EngineHooks.RunClose()
	}

	resourceErr := e.resources.Close()

	if hookErr != nil {
		hookErr = fmt.Errorf("close hooks: %w", hookErr)
	}

	return errors.Join(hookErr, resourceErr)
}
