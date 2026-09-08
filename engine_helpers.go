package ferret

import (
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/internal/resource"
	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
)

func closeEngine(
	hooks *engineHookRegistry,
	resources *resource.Manager,
) error {
	var hookErr error

	if hooks != nil {
		hookErr = hooks.runCloseHooks()
	}

	resourceErr := resources.Close()

	if hookErr != nil {
		hookErr = errors.Join(hookErr, fmt.Errorf("close hooks: %w", hookErr))
	}

	return errors.Join(hookErr, resourceErr)
}

func closeEngineOnError(
	err error,
	hooks *engineHookRegistry,
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

func closeFileSystem(filesystem ferretfs.FileSystem) error {
	if err := filesystem.Close(); err != nil {
		return fmt.Errorf("close filesystem: %w", err)
	}

	return nil
}
