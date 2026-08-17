package ferret

import (
	"errors"
	"fmt"

	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
)

func closeEngine(
	hooks *engineHookRegistry,
	filesystem ferretfs.FileSystem,
	network ferretnet.Network,
	ownsNetwork bool,
) error {
	hookErr := hooks.runCloseHooks()
	filesystemErr := closeFileSystem(filesystem)

	if ownsNetwork {
		ferretnet.CloseIdleNetworkConnections(network)
	}

	if hookErr != nil {
		hookErr = errors.Join(hookErr, fmt.Errorf("close hooks: %w", hookErr))
	}

	return errors.Join(hookErr, filesystemErr)
}

func closeEngineOnError(
	err error,
	hooks *engineHookRegistry,
	filesystem ferretfs.FileSystem,
	network ferretnet.Network,
	ownsNetwork bool,
) error {
	if err != nil {
		closeErr := closeEngine(hooks, filesystem, network, ownsNetwork)

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
