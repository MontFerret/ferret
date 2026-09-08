package engine

import (
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
)

func closeEngine(
	hooks *host.EngineHooks,
	resources *resource.Manager,
) error {
	var hookErr error

	if hooks != nil {
		hookErr = hooks.RunClose()
	}

	resourceErr := resources.Close()

	if hookErr != nil {
		hookErr = errors.Join(hookErr, fmt.Errorf("close hooks: %w", hookErr))
	}

	return errors.Join(hookErr, resourceErr)
}

func closeEngineOnError(
	err error,
	hooks *host.EngineHooks,
	resources *resource.Manager,
) error {
	if err != nil {
		closeErr := closeEngine(hooks, resources)

		if closeErr != nil {
			return errors.Join(err, fmt.Errorf("close engine: %w", closeErr))
		}
	}

	return err
}
