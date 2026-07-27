package ferret

import (
	"errors"
	"fmt"

	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
)

func closeEngine(hooks *engineHookRegistry, network ferretnet.Network, ownsNetwork bool) error {
	closeErr := hooks.runCloseHooks()

	if ownsNetwork {
		ferretnet.CloseIdleNetworkConnections(network)
	}

	if closeErr != nil {
		return errors.Join(closeErr, fmt.Errorf("close hooks: %w", closeErr))
	}

	return nil
}

func closeEngineOnError(err error, hooks *engineHookRegistry, network ferretnet.Network, ownsNetwork bool) error {
	if err != nil {
		closeErr := closeEngine(hooks, network, ownsNetwork)

		if closeErr != nil {
			return errors.Join(err, fmt.Errorf("close hooks: %w", closeErr))
		}
	}

	return err
}
