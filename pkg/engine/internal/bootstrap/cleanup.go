package bootstrap

import (
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
)

func rollback(err error, hooks *host.EngineHooks, resources *resource.Manager) error {
	if err == nil {
		return nil
	}

	var hookErr error

	if hooks != nil {
		hookErr = hooks.RunClose()
	}

	resourceErr := resources.Close()

	if hookErr != nil {
		hookErr = fmt.Errorf("close hooks: %w", hookErr)
	}

	closeErr := errors.Join(hookErr, resourceErr)
	if closeErr != nil {
		return errors.Join(err, fmt.Errorf("close engine: %w", closeErr))
	}

	return err
}
